package model

import (
	"time"
)

type User struct {
	ID string `gorm:"primaryKey"`
	Name string
	Email string `gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-"` 
	Role string `gorm:"not null"` //admin dan user
	IsActive bool `gorm:"default:true"`
	CreatedAt time.Time
}

type ImportBatch struct {
	ID string `gorm:"primaryKey"`
	Module string
	FileName string
	UploadedBy unit
	RowCount int
	SuccessCount int
	FailedCount int
	Status string //preview atau commited
	CreatedAt time.Time
}

type APTransaction struct {
	ID string `gorm:"primaryKey"`
	DocumentNo string `gorm:"uniqueIndex;not null"`
	VendorName string `gorm:"index"`
	BuyerName string `gorm:"index"`
	OrderDeadline time.Time `gorm:"index"`
	Description string
	SourceDocument string
	PayentStatus string	`gorm:"index"`
	DocumentTotal decimal.Decimal `gorm:"type:decimal(18,2)"`
	DocumentStatus string `gorm:"index"`
	ImportBatchID unit
}

type ARTransaction struct {
	ID string `gorm:"primaryKey"`
	DocumentNo string `gorm:"uniqueIndex;not null"`
	PartenerName string `gorm:"index"`
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
	ImportBatchID unit
}

type Report struct {
	ID unit `gorm:"primaryKey"`
	Type string
	Result string //JSON dari reports.Table
	CreatedBy unit
	CreatedAt time.Time
}

