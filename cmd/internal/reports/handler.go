package reports

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/model"
)

func Generate(db *gorm.DB, kind string) gin.HandlerFunc {
	svc := New(db)
	return func(c *gin.Context) {
		var in struct {
			Start string `json:"start_date"`
			End   string `json:"end_date"`
			OpeningBalance decimal.Decimal `json:"opening_balance"`
		}
		if c.ShouldBindJSON(&in) != nil {
			c.JSON(400, gin.H{"error": "start_date dan end_date harus diisi (format YYYY-MM-DD), end >= start"})
			return
		}
		start, err := time.Parse("2006-01-02", in.Start)
		if err != nil {
			c.JSON(400, gin.H{"error": "start_date harus menggunakan format YYYY-MM-DD"})
			return
		}
		end, err := time.Parse("2006-01-02", in.End)
		if err != nil || end.Before(start) {
			c.JSON(400, gin.H{"error": "end_date harus menggunakan format YYYY-MM-DD dan tidak boleh sebelum start_date"})
			return
		}
		p := Params{Start: start, End: end, Opening: in.OpeningBalance}
		var t Table
		switch kind {
		case "forecasting":
			t, err = svc.Forecast(p)
		case "profit-loss":
			t, err = svc.ProfitLoss(p)
		default:
			t, err = svc.Cashflow(p)
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "Gagal menghitung laporan"})
			return
		}
		js, _ := json.Marshal(t)
		r := model.Report{Type: kind, Result: string(js), CreatedBy: c.GetUint("uid")}
		db.Create(&r)
		c.JSON(200, gin.H{"report_id": r.ID, "table": t})
	}
}

func Load(db *gorm.DB, id string) (Table, bool) {
	var r model.Report
	var t Table
	if db.First(&r, id).Error != nil || json.Unmarshal([]byte(r.Result), &t) != nil {
		return t, false
	}
	return t, true
}
