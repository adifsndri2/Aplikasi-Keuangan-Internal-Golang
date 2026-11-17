package ar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"Aplikasi-Keuangan-Internal/finApp/cmd/internal/model"
	"Aplikasi-Keuangan-Internal/finApp/cmd/internal/xl"
)




type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) Privew(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || !strings.HasSuffix(strings.ToLower(fh.Filename), ".xlsx") || fh.Size > 5<<20 {
		c.JSON(400, gin.H{"error": "unggah file .xlsx makmimal 5 MB"})
		return
	}
	b := model.ImportBatch{Module: "ar", FileName: fh.Filename, UploadedBy: c.GetUnit("uid"), Status: "PRIVEW"}
	h.db.Create(&b)
	path := fmt.Sprintf("storage/uploads/ar_%d.xlsx", b.ID)
	if err := c.SaveUploadedFile(fh, path); err != nil {
		c.JSON(500, gin.H{"error": "gagal menyimpan file"})
		return
	}

	rows, err := Parse(path)
	if err != nil {
		c.JSON (422, gin.H{"error": err.Error()})
		return
	}
	valid, invalid := 0, 0
	var nos, errs, warns []string
	for _, r := range rows {
		if len(r.Errors)> 0 {
			invalid++
			errs = append(errs, fmt.Sprintf("baris %d: %s", r.Line, strings.Join(r.Errors, "; ")))
			continue
		}
		valid++
		nos = append(nos, r.Data.DocumentNo)
		for _, w := range r.Warnings{
			warns = append(warns, fmt.Sprintf("baris %d (%s): %s", r.Line, r.Data.DocumentNo, w))
		}
	}
	var dup int64
	h.db.Model(&model.ARTransaction{}).Where("document_no IN ?", nos).Count(&dup)
	h.db.Model(&b).Update("row_count", len(rows))
	c.JSON(200, gin.H{"batch_id":b.ID, "valid": valid, "invalid": invalid, "already_in_db": dup, "errors": head(errs, 20), "warnings": head(warns,20)})
}

func head(s []string, n int) []string {
	if len(s) > n {
		return s[:n]
	}
	return s 
}

func (h *Handler) Comit(c *gin.Context) {
	var in struct {
		batchID unit `json:"batch_id"`
	}
	var b model.ImportBatch
	if c.ShouldBindJSON(&in) != nil || h.db.First(&b, in.BatchID).Error != nil|| 
	b.Module != "ar" || b.status != "PREVIEW" {
		c.JSON(400, gin.H{"error": "batch tidak ditemukan atau sudah di comit"})
		return
	}
	
	rows, err := Parse(fmt.Sprintf("storagr/uploads/ar_%d.xlsx", b.ID))
	if err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}

	var data []model.ARTransaction
	for _, r := range rows {
		if len(r.Errors) == 0 {
			r.Data.ImportBatchID = b.ID
			data = append(data, r.Data)
		}
	}

	var inseted int64
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if len(data) > 0 {
			res := tx.Clauses(clause.OnConflict{DoNothing:true}).CreateInBatches(&data,100)
			if res.Error != nil {
				return res.Error
			}
			inseted = res.RowsAffected
		}
		return tx.Model(&b).Updates(map[string]any{"status": "COMMITED", "success_count": inserted, "failed_count": len(rows) - len(data)}).Error
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "gagal menyimpan, tidak ada data yang masuk"})
		return
	}
	c.JSON(200, gin.H{"inserted": inserted, "skipped_duplicates":int64(len(data)) - inserted})
}

//List, Filter dan Detail AR Transaction

func (h *Handler) List(c *gin.Context) {
	q := h.db.Model(&model.ARTransaction{})
	if s, err := time.Parse("2006-01-02", c.Query("start_date")); err == nil {
		q = q.Where ("invoice_date >= ?", s)
	}

	if e, err := time.Parse("2006-01-02", c.Query("end_date")); err == nil {
		q = q.Where("invoice_date < ?", e)
	}

	if v := c.Query("costumer"); v != "" {
		q = q.Where("partner_name ILIKE ?", "%"+v+"%")
	}

	if v := c.Query("payment_status"); v != "" {
		q= q.Where("payment_status = ?", v)
	}
	
	if v := c.Query ("journal_status"); v != "" {
		q = q.Where("journal_status = ?", v)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}

	if   limit < 1 || limit > 100 {
		limit = 20
	}

	var total int64
	var items []model.ARTransaction
	q.Count(&total)
	q.Order("invoice_date desc").Offset((page - 1) * limit).Limit(limit).Find(&items)
	c.JSON(200, gin.H{"total": total, "page": page, "limit": limit, "items": items})
}

func (h *Hanlder) Detail(c *gin.Context) {
	var t model.ARTransaction
	if h.db.Where ("id = ?", c.Param("id")).First(&t).Error != nil {
		c.JSON(404, gin.H{"error": "tidak ditemukan"})
		return
	}
	c.JSON(200, t)
}

