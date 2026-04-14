package store

import (
	"time"

	"github.com/nenrinyear/small-library-system/go-api/internal/domain"
)

type itemModel struct {
	ID            int64               `gorm:"column:id;primaryKey;autoIncrement"`
	QrID          string              `gorm:"column:qr_id"`
	Title         *string             `gorm:"column:title"`
	Publisher     *string             `gorm:"column:publisher"`
	PublishedDate *time.Time          `gorm:"column:published_date"`
	Description   *string             `gorm:"column:description"`
	RentalStatus  domain.RentalStatus `gorm:"column:rental_status"`
	IsRegistered  bool                `gorm:"column:is_registered"`
	IsDeleted     bool                `gorm:"column:is_deleted"`
	CreatedAt     time.Time           `gorm:"column:created_at"`
	UpdatedAt     time.Time           `gorm:"column:updated_at"`
}

func (itemModel) TableName() string {
	return "items"
}

type tagModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	IsDeleted bool      `gorm:"column:is_deleted"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (tagModel) TableName() string {
	return "tags"
}

type itemTagModel struct {
	ItemID    int64     `gorm:"column:item_id;primaryKey"`
	TagID     int64     `gorm:"column:tag_id;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (itemTagModel) TableName() string {
	return "item_tags"
}

type rentalModel struct {
	ID         int64                     `gorm:"column:id;primaryKey;autoIncrement"`
	ItemID     int64                     `gorm:"column:item_id"`
	RentedBy   string                    `gorm:"column:rented_by"`
	RentedAt   time.Time                 `gorm:"column:rented_at"`
	DueDate    time.Time                 `gorm:"column:due_date"`
	ReturnedAt *time.Time                `gorm:"column:returned_at"`
	Status     domain.RentalRecordStatus `gorm:"column:status"`
	CreatedAt  time.Time                 `gorm:"column:created_at"`
	UpdatedAt  time.Time                 `gorm:"column:updated_at"`
}

func (rentalModel) TableName() string {
	return "rentals"
}

func mapItemSummary(model itemModel) domain.ItemSummary {
	return domain.ItemSummary{
		ID:            model.ID,
		QrID:          model.QrID,
		Title:         model.Title,
		Publisher:     model.Publisher,
		PublishedDate: model.PublishedDate,
		Description:   model.Description,
		RentalStatus:  model.RentalStatus,
		CreatedAt:     model.CreatedAt,
		Tags:          []domain.Tag{},
	}
}

func mapRental(model rentalModel) domain.Rental {
	return domain.Rental{
		ID:         model.ID,
		ItemID:     model.ItemID,
		RentedBy:   model.RentedBy,
		RentedAt:   model.RentedAt,
		DueDate:    model.DueDate,
		ReturnedAt: model.ReturnedAt,
		Status:     model.Status,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func mapTag(model tagModel) domain.Tag {
	return domain.Tag{
		ID:   model.ID,
		Name: model.Name,
	}
}
