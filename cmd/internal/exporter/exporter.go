package exporter

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/finApp/cmd/internal/reports"
)

func XLSX(t reports.Table) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	sh := "Lampiran"
	f.SetCellValue(sh, "A1", t.Title+" ["+t.Label+"]")
	for i, c := range t.Columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sh, cell, c)
	}
	for r, row := range t.Rows {
		keys := reportRowKeys(row)
		for i, key := range keys {
			v := fmt.Sprint(row[key])
			cell, _ := excelize.CoordinatesToCellName(i+1, r+4)
			if n, err := strconv.ParseFloat(v, 64); err == nil && i > 0 {
				f.SetCellValue(sh, cell, n)
			} else {
				f.SetCellValue(sh, cell, v)
			}
		}
	}
	for i, n := range t.Notes {
		cell, _ := excelize.CoordinatesToCellName(1, len(t.Rows)+5+i)
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
	p.AddPage()
	p.SetFont("Arial", "B", 14)
	p.Cell(0, 10, t.Title+" ["+t.Label+"]")
	p.Ln(12)
	p.SetFont("Arial", "", 9)
	for _, row := range t.Rows {
		keys := reportRowKeys(row)
		for i, key := range keys {
			align := "R"
			if i == 0 {
				align = "L"
			}
			p.CellFormat(65, 6, fmt.Sprint(row[key]), "1", 0, align, false, 0, "")
		}
		p.Ln(-1)
	}
	p.Ln(4)
	for _, n := range t.Notes {
		p.MultiCell(0, 5, n, "", "L", false)
	}

	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func reportRowKeys(row map[string]any) []string {
	preferred := []string{"period", "in", "out", "value", "total_in", "total_out"}
	keys := make([]string, 0, len(row))
	for _, key := range preferred {
		if _, ok := row[key]; ok {
			keys = append(keys, key)
		}
	}
	extras := make([]string, 0)
	for key := range row {
		found := false
		for _, existing := range keys {
			if key == existing {
				found = true
				break
			}
		}
		if !found {
			extras = append(extras, key)
		}
	}
	sort.Strings(extras)
	keys = append(keys, extras...)
	return keys
}

// Handler Report
func Export(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, ok := reports.Load(db, c.Param("id"))
		if !ok {
			c.JSON(404, gin.H{"error": "Laporan tidak ditemukan"})
			return
		}
		name := "laporan-" + c.Param("id")
		switch c.DefaultQuery("format", "PDF") {
		case "xlsx":
			b, err := XLSX(t)
			if err != nil {
				c.Status(500)
				return
			}
			c.Header("Content-Disposition", "attachment; filename = "+name+".xlsx")
			c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b)
		default:
			b, err := PDF(t)
			if err != nil {
				c.Status(500)
				return
			}
			c.Header("Content-Disposition", "attachment; filename = "+name+".pdf")
			c.Data(200, "application/pdf", b)
		}
	}
}
