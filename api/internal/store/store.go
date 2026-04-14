package store

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nenrinyear/small-library-system/api/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrConflict         = errors.New("state conflict")
	ErrValidationFailed = errors.New("validation failed")
	ErrForbiddenAction  = errors.New("forbidden action")
)

type Store struct {
	db           *gorm.DB
	historyLimit int
}

type RegisterItemInput struct {
	Title         string
	Publisher     string
	PublishedDate *time.Time
	Description   *string
	TagIDs        []int64
}

type RentItemInput struct {
	RentedBy string
	DueDate  time.Time
}

type ReturnItemInput struct {
	ReturnedBy string
}

type SearchItemsInput struct {
	Query     string
	Publisher string
	TagID     *int64
	Page      int
	PageSize  int
}

func New(db *gorm.DB, historyLimit int) *Store {
	if historyLimit <= 0 {
		historyLimit = 10
	}

	return &Store{
		db:           db,
		historyLimit: historyLimit,
	}
}

func (s *Store) ListTags(ctx context.Context) ([]domain.Tag, error) {
	var models []tagModel
	if err := s.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Order("name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(models))
	for _, model := range models {
		tags = append(tags, mapTag(model))
	}

	return tags, nil
}

func (s *Store) GetDashboard(ctx context.Context) (domain.Dashboard, error) {
	availableItems, err := s.listAvailableItems(ctx, 50)
	if err != nil {
		return domain.Dashboard{}, err
	}

	rentedItems, err := s.listRentedItems(ctx, 50)
	if err != nil {
		return domain.Dashboard{}, err
	}

	recentRentals, err := s.listRecentRentals(ctx, s.historyLimit)
	if err != nil {
		return domain.Dashboard{}, err
	}

	return domain.Dashboard{
		AvailableItems: availableItems,
		RentedItems:    rentedItems,
		RecentRentals:  recentRentals,
	}, nil
}

func (s *Store) GetItemDetailByQRID(ctx context.Context, qrID string) (domain.ItemDetail, error) {
	if strings.TrimSpace(qrID) == "" {
		return domain.ItemDetail{}, ErrValidationFailed
	}

	var model itemModel
	err := s.db.WithContext(ctx).
		Where("qr_id = ? AND is_deleted = ?", qrID, false).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ItemDetail{}, ErrNotFound
		}
		return domain.ItemDetail{}, err
	}

	item := mapItemSummary(model)
	tags, err := s.listTagsByItemID(ctx, item.ID)
	if err != nil {
		return domain.ItemDetail{}, err
	}
	item.Tags = tags

	activeRental, err := s.findActiveRental(ctx, item.ID)
	if err != nil {
		return domain.ItemDetail{}, err
	}

	history, err := s.ListRentalsByItemID(ctx, item.ID, s.historyLimit)
	if err != nil {
		return domain.ItemDetail{}, err
	}

	return domain.ItemDetail{
		Item:          item,
		ActiveRental:  activeRental,
		RentalHistory: history,
	}, nil
}

func (s *Store) ListRentalsByItemID(ctx context.Context, itemID int64, limit int) ([]domain.Rental, error) {
	if itemID <= 0 {
		return nil, ErrValidationFailed
	}
	if limit <= 0 {
		limit = s.historyLimit
	}

	var models []rentalModel
	if err := s.db.WithContext(ctx).
		Where("item_id = ?", itemID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	rentals := make([]domain.Rental, 0, len(models))
	for _, model := range models {
		rentals = append(rentals, mapRental(model))
	}

	return rentals, nil
}

func (s *Store) GenerateUnregisteredItems(ctx context.Context, count int) ([]domain.ItemSummary, error) {
	if count <= 0 || count > 100 {
		return nil, ErrValidationFailed
	}

	created := make([]domain.ItemSummary, 0, count)
	for i := 0; i < count; i++ {
		item, err := s.generateOneItem(ctx)
		if err != nil {
			return nil, err
		}
		created = append(created, item)
	}

	return created, nil
}

func (s *Store) RegisterItemByQRID(ctx context.Context, qrID string, input RegisterItemInput) (domain.ItemDetail, error) {
	if strings.TrimSpace(qrID) == "" || strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Publisher) == "" || len(input.TagIDs) == 0 {
		return domain.ItemDetail{}, ErrValidationFailed
	}

	tagIDs, err := normalizeTagIDs(input.TagIDs)
	if err != nil {
		return domain.ItemDetail{}, err
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domain.ItemDetail{}, tx.Error
	}
	defer rollback(tx)

	var item itemModel
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("qr_id = ?", qrID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ItemDetail{}, ErrNotFound
		}
		return domain.ItemDetail{}, err
	}

	if item.IsDeleted || item.IsRegistered {
		return domain.ItemDetail{}, ErrConflict
	}

	if err := ensureTagsActive(tx, tagIDs); err != nil {
		return domain.ItemDetail{}, err
	}

	err = tx.Model(&itemModel{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"title":          strings.TrimSpace(input.Title),
			"publisher":      strings.TrimSpace(input.Publisher),
			"published_date": input.PublishedDate,
			"description":    nullableText(input.Description),
			"is_registered":  true,
			"updated_at":     nowInTokyo(),
		}).Error
	if err != nil {
		return domain.ItemDetail{}, err
	}

	if err := tx.Where("item_id = ?", item.ID).Delete(&itemTagModel{}).Error; err != nil {
		return domain.ItemDetail{}, err
	}

	links := make([]itemTagModel, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		links = append(links, itemTagModel{ItemID: item.ID, TagID: tagID})
	}
	if err := tx.Create(&links).Error; err != nil {
		return domain.ItemDetail{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return domain.ItemDetail{}, err
	}

	return s.GetItemDetailByQRID(ctx, qrID)
}

func (s *Store) RentItem(ctx context.Context, itemID int64, input RentItemInput) error {
	if itemID <= 0 || strings.TrimSpace(input.RentedBy) == "" || input.DueDate.IsZero() {
		return ErrValidationFailed
	}

	today := nowInTokyo().Truncate(24 * time.Hour)
	due := input.DueDate.Truncate(24 * time.Hour)
	if due.Before(today) {
		return ErrValidationFailed
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer rollback(tx)

	var item itemModel
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND is_deleted = ? AND is_registered = ?", itemID, false, true).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if item.RentalStatus != domain.RentalStatusAvailable {
		return ErrConflict
	}

	var activeCount int64
	err = tx.Model(&rentalModel{}).
		Where("item_id = ? AND status = ?", itemID, domain.RentalRecordStatusActive).
		Count(&activeCount).Error
	if err != nil {
		return err
	}
	if activeCount > 0 {
		return ErrConflict
	}

	now := nowInTokyo()
	if err := tx.Create(&rentalModel{
		ItemID:   itemID,
		RentedBy: strings.TrimSpace(input.RentedBy),
		RentedAt: now,
		DueDate:  due,
		Status:   domain.RentalRecordStatusActive,
	}).Error; err != nil {
		return err
	}

	if err := tx.Model(&itemModel{}).
		Where("id = ?", itemID).
		Updates(map[string]any{
			"rental_status": domain.RentalStatusRented,
			"updated_at":    now,
		}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *Store) ReturnItem(ctx context.Context, itemID int64, input ReturnItemInput) error {
	if itemID <= 0 || strings.TrimSpace(input.ReturnedBy) == "" {
		return ErrValidationFailed
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer rollback(tx)

	var item itemModel
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND is_deleted = ? AND is_registered = ?", itemID, false, true).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if item.RentalStatus != domain.RentalStatusRented {
		return ErrConflict
	}

	var rental rentalModel
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("item_id = ? AND status = ?", itemID, domain.RentalRecordStatusActive).
		Order("id DESC").
		First(&rental).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConflict
		}
		return err
	}

	if rental.RentedBy != strings.TrimSpace(input.ReturnedBy) {
		return ErrForbiddenAction
	}

	now := nowInTokyo()
	if err := tx.Model(&rentalModel{}).
		Where("id = ?", rental.ID).
		Updates(map[string]any{
			"status":      domain.RentalRecordStatusReturned,
			"returned_at": now,
			"updated_at":  now,
		}).Error; err != nil {
		return err
	}

	if err := tx.Model(&itemModel{}).
		Where("id = ?", itemID).
		Updates(map[string]any{
			"rental_status": domain.RentalStatusAvailable,
			"updated_at":    now,
		}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *Store) SearchItems(ctx context.Context, input SearchItemsInput) (domain.SearchResult, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := (page - 1) * pageSize

	base := s.db.WithContext(ctx).
		Model(&itemModel{}).
		Where("items.is_deleted = ? AND items.is_registered = ?", false, true)

	if q := strings.TrimSpace(input.Query); q != "" {
		like := "%" + q + "%"
		base = base.Where(s.db.Where("items.title LIKE ?", like).Or("items.description LIKE ?", like))
	}
	if publisher := strings.TrimSpace(input.Publisher); publisher != "" {
		base = base.Where("items.publisher = ?", publisher)
	}
	if input.TagID != nil {
		base = base.Joins("JOIN item_tags ON item_tags.item_id = items.id").
			Where("item_tags.tag_id = ?", *input.TagID)
	}

	var total int64
	if err := base.Distinct("items.id").Count(&total).Error; err != nil {
		return domain.SearchResult{}, err
	}

	var models []itemModel
	if err := base.
		Distinct("items.id").
		Select("items.*").
		Order("items.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&models).Error; err != nil {
		return domain.SearchResult{}, err
	}

	items := make([]domain.ItemSummary, 0, len(models))
	for _, model := range models {
		items = append(items, mapItemSummary(model))
	}

	if err := s.attachTags(ctx, items); err != nil {
		return domain.SearchResult{}, err
	}

	return domain.SearchResult{
		Items:      items,
		TotalCount: int(total),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *Store) ExportSearch(ctx context.Context) ([]map[string]any, error) {
	var models []itemModel
	if err := s.db.WithContext(ctx).
		Model(&itemModel{}).
		Select("qr_id", "title", "publisher", "description").
		Where("is_deleted = ? AND is_registered = ?", false, true).
		Find(&models).Error; err != nil {
		return nil, err
	}

	payload := make([]map[string]any, 0, len(models))
	for _, model := range models {
		payload = append(payload, map[string]any{
			"objectID":    model.QrID,
			"title":       nullableString(model.Title),
			"publisher":   nullableString(model.Publisher),
			"description": nullableString(model.Description),
		})
	}

	return payload, nil
}

func (s *Store) listAvailableItems(ctx context.Context, limit int) ([]domain.ItemSummary, error) {
	var models []itemModel
	if err := s.db.WithContext(ctx).
		Where("is_deleted = ? AND is_registered = ? AND rental_status = ?", false, true, domain.RentalStatusAvailable).
		Order("created_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.ItemSummary, 0, len(models))
	for _, model := range models {
		items = append(items, mapItemSummary(model))
	}

	if err := s.attachTags(ctx, items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Store) listRentedItems(ctx context.Context, limit int) ([]domain.RentedSummary, error) {
	var rentals []rentalModel
	if err := s.db.WithContext(ctx).
		Where("status = ?", domain.RentalRecordStatusActive).
		Order("due_date ASC").
		Limit(limit).
		Find(&rentals).Error; err != nil {
		return nil, err
	}
	if len(rentals) == 0 {
		return []domain.RentedSummary{}, nil
	}

	itemIDs := make([]int64, 0, len(rentals))
	for _, rental := range rentals {
		itemIDs = append(itemIDs, rental.ItemID)
	}

	var items []itemModel
	if err := s.db.WithContext(ctx).
		Where("id IN ?", itemIDs).
		Where("is_deleted = ? AND is_registered = ? AND rental_status = ?", false, true, domain.RentalStatusRented).
		Find(&items).Error; err != nil {
		return nil, err
	}

	itemsByID := make(map[int64]domain.ItemSummary, len(items))
	for _, model := range items {
		itemsByID[model.ID] = mapItemSummary(model)
	}

	rented := make([]domain.RentedSummary, 0, len(rentals))
	for _, rental := range rentals {
		item, ok := itemsByID[rental.ItemID]
		if !ok {
			continue
		}
		rented = append(rented, domain.RentedSummary{
			ItemSummary: item,
			RentedBy:    rental.RentedBy,
			RentedAt:    rental.RentedAt,
			DueDate:     rental.DueDate,
		})
	}

	itemSummaries := make([]domain.ItemSummary, 0, len(rented))
	for _, row := range rented {
		itemSummaries = append(itemSummaries, row.ItemSummary)
	}
	if err := s.attachTags(ctx, itemSummaries); err != nil {
		return nil, err
	}
	for i := range rented {
		rented[i].Tags = itemSummaries[i].Tags
	}

	return rented, nil
}

func (s *Store) listRecentRentals(ctx context.Context, limit int) ([]domain.RentalEvent, error) {
	var rentals []rentalModel
	if err := s.db.WithContext(ctx).
		Order("updated_at DESC").
		Limit(limit).
		Find(&rentals).Error; err != nil {
		return nil, err
	}
	if len(rentals) == 0 {
		return []domain.RentalEvent{}, nil
	}

	itemIDs := make([]int64, 0, len(rentals))
	for _, rental := range rentals {
		itemIDs = append(itemIDs, rental.ItemID)
	}

	var items []itemModel
	if err := s.db.WithContext(ctx).
		Where("id IN ? AND is_deleted = ?", itemIDs, false).
		Find(&items).Error; err != nil {
		return nil, err
	}

	itemsByID := make(map[int64]itemModel, len(items))
	for _, item := range items {
		itemsByID[item.ID] = item
	}

	events := make([]domain.RentalEvent, 0, len(rentals))
	for _, rental := range rentals {
		item, ok := itemsByID[rental.ItemID]
		if !ok {
			continue
		}
		events = append(events, domain.RentalEvent{
			RentalID:   rental.ID,
			ItemID:     rental.ItemID,
			QrID:       item.QrID,
			Title:      item.Title,
			RentedBy:   rental.RentedBy,
			RentedAt:   rental.RentedAt,
			DueDate:    rental.DueDate,
			ReturnedAt: rental.ReturnedAt,
			UpdatedAt:  rental.UpdatedAt,
		})
	}

	return events, nil
}

func (s *Store) listTagsByItemID(ctx context.Context, itemID int64) ([]domain.Tag, error) {
	if itemID <= 0 {
		return nil, ErrValidationFailed
	}

	type row struct {
		ID   int64
		Name string
	}

	var rows []row
	err := s.db.WithContext(ctx).
		Table("tags").
		Select("tags.id, tags.name").
		Joins("JOIN item_tags ON item_tags.tag_id = tags.id").
		Where("item_tags.item_id = ? AND tags.is_deleted = ?", itemID, false).
		Order("tags.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, domain.Tag{ID: row.ID, Name: row.Name})
	}

	return tags, nil
}

func (s *Store) attachTags(ctx context.Context, items []domain.ItemSummary) error {
	if len(items) == 0 {
		return nil
	}

	itemIDs := make([]int64, 0, len(items))
	for _, item := range items {
		itemIDs = append(itemIDs, item.ID)
	}

	type row struct {
		ItemID int64
		TagID  int64
		Name   string
	}

	var rows []row
	err := s.db.WithContext(ctx).
		Table("item_tags").
		Select("item_tags.item_id, tags.id AS tag_id, tags.name").
		Joins("JOIN tags ON tags.id = item_tags.tag_id").
		Where("tags.is_deleted = ?", false).
		Where("item_tags.item_id IN ?", itemIDs).
		Order("tags.name ASC").
		Scan(&rows).Error
	if err != nil {
		return err
	}

	tagMap := make(map[int64][]domain.Tag, len(items))
	for _, row := range rows {
		tagMap[row.ItemID] = append(tagMap[row.ItemID], domain.Tag{ID: row.TagID, Name: row.Name})
	}

	for i := range items {
		items[i].Tags = tagMap[items[i].ID]
	}

	return nil
}

func (s *Store) findActiveRental(ctx context.Context, itemID int64) (*domain.Rental, error) {
	var model rentalModel
	err := s.db.WithContext(ctx).
		Where("item_id = ? AND status = ?", itemID, domain.RentalRecordStatusActive).
		Order("id DESC").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	rental := mapRental(model)
	return &rental, nil
}

func (s *Store) generateOneItem(ctx context.Context) (domain.ItemSummary, error) {
	for attempts := 0; attempts < 10; attempts++ {
		qrID, err := generateQRID(16)
		if err != nil {
			return domain.ItemSummary{}, err
		}

		model := itemModel{
			QrID:         qrID,
			RentalStatus: domain.RentalStatusAvailable,
			IsRegistered: false,
			IsDeleted:    false,
		}

		err = s.db.WithContext(ctx).Create(&model).Error
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				continue
			}
			return domain.ItemSummary{}, err
		}

		return mapItemSummary(model), nil
	}

	return domain.ItemSummary{}, ErrConflict
}

func generateQRID(length int) (string, error) {
	if length <= 0 {
		return "", ErrValidationFailed
	}

	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	max := big.NewInt(int64(len(alphabet)))
	b := strings.Builder{}
	b.Grow(length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteByte(alphabet[n.Int64()])
	}

	return b.String(), nil
}

func normalizeTagIDs(tagIDs []int64) ([]int64, error) {
	if len(tagIDs) == 0 {
		return nil, ErrValidationFailed
	}

	normalized := make([]int64, 0, len(tagIDs))
	seen := make(map[int64]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID <= 0 {
			return nil, ErrValidationFailed
		}
		if _, ok := seen[tagID]; ok {
			continue
		}
		seen[tagID] = struct{}{}
		normalized = append(normalized, tagID)
	}
	slices.Sort(normalized)

	return normalized, nil
}

func ensureTagsActive(db *gorm.DB, tagIDs []int64) error {
	var count int64
	if err := db.Model(&tagModel{}).
		Where("id IN ?", tagIDs).
		Where("is_deleted = ?", false).
		Count(&count).Error; err != nil {
		return err
	}

	if count != int64(len(tagIDs)) {
		return ErrValidationFailed
	}

	return nil
}

func nullableText(value *string) *string {
	if value == nil {
		return nil
	}
	v := strings.TrimSpace(*value)
	if v == "" {
		return nil
	}
	return &v
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func nowInTokyo() time.Time {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

func rollback(tx *gorm.DB) {
	if tx != nil {
		_ = tx.Rollback().Error
	}
}
