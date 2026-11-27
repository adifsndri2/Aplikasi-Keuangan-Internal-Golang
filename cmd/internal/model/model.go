package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID uint `gorm:"primaryKey"`
	Name string
	Email string `gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"` 
	Role string `gorm:"not null"` //admin dan user
	IsActive bool `gorm:"default:true"`
	CreatedAt time.Time
}

type ImportBatch struct {
	ID uint `gorm:"primaryKey"`
	Module string
	FileName string
	UploadedBy uint
	RowCount int
	SuccessCount int
	FailedCount int
	Status string //preview atau commited
	CreatedAt time.Time
}

type APTransaction struct {
	ID uint `gorm:"primaryKey"`
	DocumentNo string `gorm:"uniqueIndex;not null"`
	VendorName string `gorm:"index"`
	BuyerName string `gorm:"index"`
	OrderDeadline time.Time `gorm:"index"`
	Description string
	SourceDocument string
	PaymentStatus string `gorm:"index"`
	DocumentTotal decimal.Decimal `gorm:"type:decimal(18,2)"`
	DocumentStatus string `gorm:"index"`
	ImportBatchID uint
}

type ARTransaction struct {
	ID uint `gorm:"primaryKey"`
	DocumentNo string `gorm:"uniqueIndex;not null"`
	PartnerName string `gorm:"index"`
	IsInternal bool
	InvoiceDate	time.Time `gorm:"index"`
	DueDate time.Time
	TaxNumber string
	Description string
	UntaxedAmount decimal.Decimal `gorm:"type:decimal(18,2)"`
	TotalAmount decimal.Decimal `gorm:"type:decimal(18,2)"`
	CurrencyTotal decimal.Decimal `gorm:"type:decimal(18,2)"`
	PaymentStatus string `gorm:"index"`
	JournalStatus string
	ImportBatchID uint
}

type Report struct {
	ID uint `gorm:"primaryKey"`
	Type string
	Result string //JSON dari reports.Table
	CreatedBy uint
	CreatedAt time.Time
}

