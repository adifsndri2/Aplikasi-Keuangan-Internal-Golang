package ap

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/finApp/cmd/internal/model"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Preview(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || !strings.HasSuffix(strings.ToLower(fh.Filename), ".xlsx") || fh.Size > 5<<20 {
		c.JSON(400, gin.H{"error": "File harus berupa .xlsx dan ukuran maksimal 5MB"})
		return
	}
	b := model.ImportBatch{Module: "ap", FileName: fh.Filename, UploadedBy: uint(c.GetUint("uid")), Status: "PREVIEW"}
	h.DB.Create(&b)
	path := fmt.Sprintf("storage/uploads/ap_%d.xlsx", b.ID)
	if err := c.SaveUploadedFile(fh, path); err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan file"})
		return
	}
	rows, err := Parse(path)
	if err != nil {
		c.JSON(422, gin.H{"error": "Gagal membaca file: " + err.Error()})
		return
	}
	valid, invalid := 0, 0
	var nos, errs, warns []string
	for _, r := range rows {
		if len(r.Errors) > 0 {
			invalid++
			errs = append(errs, fmt.Sprintf("Baris %d: %s", r.Line, strings.Join(r.Errors, "; ")))
			continue
		}
		valid++
		nos = append(nos, r.Data.DocumentNo)
		for _, w := range r.Warnings {
			warns = append(warns, fmt.Sprintf("Baris %d (%s): %s", r.Line, r.Data.DocumentNo, w))
		}
	}
	var dup int64
	h.DB.Model(&model.APTransaction{}).Where("document_no IN ?", nos).Count(&dup)
	h.DB.Model(&b).Update("row_count", len(rows))
	c.JSON(200, gin.H{
		"batch_id": b.ID,
		"valid":    valid,
		"invalid":  invalid,
		"already_in_db": dup,
		"errors":   head(errs, 20),
		"warnings": head(warns, 20),
	})
}

func head(s []string, n int) []string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func (h *Handler) Import(c *gin.Context) {
	var in struct{ BatchID uint `json:"batch_id"` }
	var b model.ImportBatch
	if c.ShouldBindJSON(&in) != nil || h.DB.First(&b, in.BatchID).Error != nil ||
		b.Module != "ap" || b.Status != "PREVIEW" {
		c.JSON(400, gin.H{"error": "Batch ID Tidak Valid"})
		return
	}
	rows, err := Parse(fmt.Sprintf("storage/uploads/ap_%d.xlsx", b.ID))
	if err != nil {
		c.JSON(422, gin.H{"error": "Gagal membaca file: " + err.Error()})
		return
	}
	var data []model.APTransaction
	for _, r := range rows {
		if len(r.Errors) == 0 {
			r.Data.ImportBatchID = b.ID
			data = append(data, r.Data)
		}
	}
	var inserted int64
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if len(data) > 0 {
			res := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&data, 100)
			if res.Error != nil {
				return res.Error
			}
			inserted = res.RowsAffected
		}
		return tx.Model(&b).Updates(map[string]any{"status": "Comitted", "success_count": inserted, "failed_count": len(rows) - len(data)}).Error
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan data"})
		return
	}
	c.JSON(200, gin.H{"inserted": inserted, "skipped_duplicates": int64(len(data)) - inserted})
}

func (h *Handler) List(c *gin.Context) {
	q := h.DB.Model(&model.APTransaction{})
	if s, err := time.Parse("2006-01-02", c.Query("start_date")); err == nil {
		q = q.Where("order_deadline >= ?", s)
	}
	if e, err := time.Parse("2006-01-02", c.Query("end_date")); err == nil {
		q = q.Where("order_deadline <= ?", e.AddDate(0, 0, 1))
	}
	if v := c.Query("vendor"); v != "" {
		q = q.Where("vendor_name ILIKE ?", "%"+v+"%")
	}
	if v := c.Query("buyer"); v != "" {
		q = q.Where("buyer_name ILIKE ?", "%"+v+"%")
	}
	if v := c.Query("status"); v != "" {
		q = q.Where("document_status = ?", v)
	}
	if v := c.Query("payment_status"); v != "" {
		q = q.Where("payment_status = ?", v)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var total int64
	var items []model.APTransaction
	q.Count(&total)
	q.Order("order_deadline DESC").Offset((page - 1) * limit).Limit(limit).Find(&items)
	c.JSON(200, gin.H{"total": total, "page": page, "items": items})
}

func (h *Handler) Get(c *gin.Context) {
	var t model.APTransaction
	if h.DB.First(&t, c.Param("id")).Error != nil {
		c.JSON(404, gin.H{"error": "Data Tidak Ditemukan"})
		return
	}
	c.JSON(200, t)
}
