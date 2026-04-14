package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nenrinyear/small-library-system/api/internal/auth"
	"github.com/nenrinyear/small-library-system/api/internal/config"
	"github.com/nenrinyear/small-library-system/api/internal/store"
	"gorm.io/gorm"
)

type Server struct {
	e      *echo.Echo
	store  *store.Store
	authMW auth.Middleware
}

func New(cfg config.Config, db *gorm.DB) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	s := &Server{
		e:      e,
		store:  store.New(db, cfg.HistoryLimit),
		authMW: auth.New(cfg.AdminToken),
	}
	s.registerRoutes()

	return s
}

func (s *Server) Engine() *echo.Echo {
	return s.e
}

func (s *Server) registerRoutes() {
	s.e.GET("/healthz", s.handleHealth)

	api := s.e.Group("/api")
	api.GET("/health", s.handleHealth)
	api.GET("/dashboard", s.handleDashboard)
	api.GET("/tags", s.handleListTags)
	api.GET("/items/:qrId", s.handleGetItemByQR)
	api.GET("/items/:itemId/rentals", s.handleListRentals)
	api.GET("/search", s.handleSearch)
	api.GET("/search/export", s.handleExport)

	admin := api.Group("", s.authMW.RequireAdmin)
	admin.POST("/items/generate", s.handleGenerateItems)
	admin.POST("/items/:qrId/register", s.handleRegisterItem)
	admin.POST("/items/:itemId/rent", s.handleRentItem)
	admin.POST("/items/:itemId/return", s.handleReturnItem)
}

func (s *Server) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func (s *Server) handleDashboard(c echo.Context) error {
	dashboard, err := s.store.GetDashboard(c.Request().Context())
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, dashboard)
}

func (s *Server) handleListTags(c echo.Context) error {
	tags, err := s.store.ListTags(c.Request().Context())
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"tags": tags})
}

func (s *Server) handleGetItemByQR(c echo.Context) error {
	item, err := s.store.GetItemDetailByQRID(c.Request().Context(), strings.TrimSpace(c.Param("qrId")))
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

func (s *Server) handleListRentals(c echo.Context) error {
	itemID, err := parseInt64Param(c, "itemId")
	if err != nil {
		return writeError(c, err)
	}
	limit := parseIntWithDefault(c.QueryParam("limit"), 10)
	rentals, err := s.store.ListRentalsByItemID(c.Request().Context(), itemID, limit)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"rentals": rentals})
}

func (s *Server) handleSearch(c echo.Context) error {
	var tagID *int64
	tagRaw := strings.TrimSpace(c.QueryParam("tagId"))
	if tagRaw != "" {
		parsed, err := strconv.ParseInt(tagRaw, 10, 64)
		if err != nil {
			return writeError(c, store.ErrValidationFailed)
		}
		tagID = &parsed
	}

	result, err := s.store.SearchItems(c.Request().Context(), store.SearchItemsInput{
		Query:     strings.TrimSpace(c.QueryParam("q")),
		Publisher: strings.TrimSpace(c.QueryParam("publisher")),
		TagID:     tagID,
		Page:      parseIntWithDefault(c.QueryParam("page"), 1),
		PageSize:  parseIntWithDefault(c.QueryParam("pageSize"), 20),
	})
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (s *Server) handleExport(c echo.Context) error {
	rows, err := s.store.ExportSearch(c.Request().Context())
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, rows)
}

func (s *Server) handleGenerateItems(c echo.Context) error {
	var body struct {
		Count int `json:"count"`
	}
	if err := c.Bind(&body); err != nil {
		return writeError(c, store.ErrValidationFailed)
	}
	if body.Count == 0 {
		body.Count = 1
	}

	items, err := s.store.GenerateUnregisteredItems(c.Request().Context(), body.Count)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"items": items,
	})
}

func (s *Server) handleRegisterItem(c echo.Context) error {
	var body struct {
		Title         string  `json:"title"`
		Publisher     string  `json:"publisher"`
		PublishedDate *string `json:"publishedDate"`
		Description   *string `json:"description"`
		TagIDs        []int64 `json:"tagIds"`
	}
	if err := c.Bind(&body); err != nil {
		return writeError(c, store.ErrValidationFailed)
	}

	var publishedDate *time.Time
	if body.PublishedDate != nil && strings.TrimSpace(*body.PublishedDate) != "" {
		parsed, err := parseDate(*body.PublishedDate)
		if err != nil {
			return writeError(c, store.ErrValidationFailed)
		}
		publishedDate = &parsed
	}

	item, err := s.store.RegisterItemByQRID(c.Request().Context(), strings.TrimSpace(c.Param("qrId")), store.RegisterItemInput{
		Title:         body.Title,
		Publisher:     body.Publisher,
		PublishedDate: publishedDate,
		Description:   body.Description,
		TagIDs:        body.TagIDs,
	})
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, item)
}

func (s *Server) handleRentItem(c echo.Context) error {
	itemID, err := parseInt64Param(c, "itemId")
	if err != nil {
		return writeError(c, err)
	}

	var body struct {
		RentedBy string `json:"rentedBy"`
		DueDate  string `json:"dueDate"`
	}
	if err := c.Bind(&body); err != nil {
		return writeError(c, store.ErrValidationFailed)
	}

	dueDate, err := parseDate(body.DueDate)
	if err != nil {
		return writeError(c, store.ErrValidationFailed)
	}

	if err := s.store.RentItem(c.Request().Context(), itemID, store.RentItemInput{
		RentedBy: body.RentedBy,
		DueDate:  dueDate,
	}); err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) handleReturnItem(c echo.Context) error {
	itemID, err := parseInt64Param(c, "itemId")
	if err != nil {
		return writeError(c, err)
	}

	var body struct {
		ReturnedBy string `json:"returnedBy"`
	}
	if err := c.Bind(&body); err != nil {
		return writeError(c, store.ErrValidationFailed)
	}

	if err := s.store.ReturnItem(c.Request().Context(), itemID, store.ReturnItemInput{
		ReturnedBy: body.ReturnedBy,
	}); err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func writeError(c echo.Context, err error) error {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, store.ErrValidationFailed):
		status = http.StatusBadRequest
		message = "validation failed"
	case errors.Is(err, store.ErrNotFound):
		status = http.StatusNotFound
		message = "not found"
	case errors.Is(err, store.ErrConflict):
		status = http.StatusConflict
		message = "state conflict"
	case errors.Is(err, store.ErrForbiddenAction):
		status = http.StatusForbidden
		message = "forbidden"
	}

	return c.JSON(status, map[string]any{
		"error": message,
	})
}

func parseInt64Param(c echo.Context, key string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(c.Param(key)), 10, 64)
	if err != nil || parsed <= 0 {
		return 0, store.ErrValidationFailed
	}
	return parsed, nil
}

func parseIntWithDefault(raw string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseDate(raw string) (time.Time, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return time.Time{}, errors.New("date is required")
	}
	for _, layout := range []string{"2006-01-02", "2006-01"} {
		if t, err := time.ParseInLocation(layout, v, time.FixedZone("JST", 9*60*60)); err == nil {
			if layout == "2006-01" {
				return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
			}
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid date")
}
