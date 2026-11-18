package exporter

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/shopspring/decimal"

	"Aplikasi-Keuangan-Internal/finApp/cmd/internal/model"
)

func XLSX (t, reports.Table) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	sh := "Lampiran"
	f.SetCellValue(sh, "A1" t.Title+" ["+t.Label+"]")
	for i, c := range t.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sh, cell, c)
	}
	for r, row := range t.Rows {
		for i, v := range row {
			cell, _ := excelize.CoordinatesToCellName (i+1, r+4)
			if n, err := strconv.ParseFloat (v, 64); err == nil && i > 0 {
				f.SetCellValue(sh, cell, n)
			} else {
				f.SetCellValue(sh, cell, v)
			}
		}
	}
	for i, n := range t.Notes {
		cell, _ := excelize.CoordinatesToCellName(1, len(t.rows)+5+i)
		f.SetCellValue(sh, cell, n)
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func PDF(t reports.Table) ([]byte, error) {
	p := gofpdf.New("L", "mm", "A4", "")
	p.mAddPage()
	p.SetFont ("Arial", "B", 14)
	p.Cell(0, 10, t.Title+" ["+t.Label+"]")
	p.Ln(12)
	p.SetFont("Arial", "", 9)
	for _, row := range t.Rows {
		for _, v := range row {
			align := "R"
			if i == 0 {
				align = "L"
			}
			p.CellFormat(w, 6, v, "1", 0, align, false, 0, "")
		}
		p.Ln(-1)
	}
	p.Ln(4)
	for _, n := range t.Notes {
		p.MultiCell(0, 5, n, "", "L", false)
	}

	var buf bytes.Buffer
	if err := p.Ouput(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Handler Report
func Export (db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, ok := reports.Load(db, c.Param("id"))
		if !ok {
			c.JSON (404, gin.H{"error": "Laporan tidak ditemukan"})
			return
		}
		name := "laporan-" + c.Param("id")
		switch c.DefaultQuery("format", "PDF") {
		case "xlsx" :
			b, err := exporter.XLSX(t)
			if err != nil {
				c.Status(500)
				return
			}
			c.Header("Content-Disposition", "attachment; filename = "+name+".xlsx")
			c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b)
			default : 
			b, err := exporter.PDF(t)
			if err != nil {
				c.Status(500)
				return
			}
			c.Header("Content-Dispotition", "attachment; filename = "+name+".pdf")
		} 
	}
}