package domain

import "time"

type RentalStatus string

const (
	RentalStatusAvailable RentalStatus = "AVAILABLE"
	RentalStatusRented    RentalStatus = "RENTED"
)

type RentalRecordStatus string

const (
	RentalRecordStatusActive   RentalRecordStatus = "ACTIVE"
	RentalRecordStatusReturned RentalRecordStatus = "RETURNED"
)

type ItemSummary struct {
	ID            int64        `json:"id"`
	QrID          string       `json:"qrId"`
	Title         *string      `json:"title"`
	Publisher     *string      `json:"publisher"`
	PublishedDate *time.Time   `json:"publishedDate"`
	Description   *string      `json:"description"`
	RentalStatus  RentalStatus `json:"rentalStatus"`
	CreatedAt     time.Time    `json:"createdAt"`
	Tags          []Tag        `json:"tags"`
}

type ItemDetail struct {
	Item          ItemSummary `json:"item"`
	ActiveRental  *Rental     `json:"activeRental"`
	RentalHistory []Rental    `json:"rentalHistory"`
}

type Rental struct {
	ID         int64              `json:"id"`
	ItemID     int64              `json:"itemId"`
	RentedBy   string             `json:"rentedBy"`
	RentedAt   time.Time          `json:"rentedAt"`
	DueDate    time.Time          `json:"dueDate"`
	ReturnedAt *time.Time         `json:"returnedAt"`
	Status     RentalRecordStatus `json:"status"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Dashboard struct {
	AvailableItems []ItemSummary   `json:"availableItems"`
	RentedItems    []RentedSummary `json:"rentedItems"`
	RecentRentals  []RentalEvent   `json:"recentRentals"`
}

type RentedSummary struct {
	ItemSummary
	RentedBy string    `json:"rentedBy"`
	RentedAt time.Time `json:"rentedAt"`
	DueDate  time.Time `json:"dueDate"`
}

type RentalEvent struct {
	RentalID   int64      `json:"rentalId"`
	ItemID     int64      `json:"itemId"`
	QrID       string     `json:"qrId"`
	Title      *string    `json:"title"`
	RentedBy   string     `json:"rentedBy"`
	RentedAt   time.Time  `json:"rentedAt"`
	DueDate    time.Time  `json:"dueDate"`
	ReturnedAt *time.Time `json:"returnedAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type SearchResult struct {
	Items      []ItemSummary `json:"items"`
	TotalCount int           `json:"totalCount"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
}
