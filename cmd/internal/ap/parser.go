package ap

import (
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/model"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/xl"
)

type Row struct {
	Line   int
	Data   model.APTransaction
	Errors []string
	Warnings []string
}

func Parse(path string) ([]Row, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}
	idx, err := xl.Header(rows, "Order Reference", "Vendor", "Buyer", "Order Deadline", "Total", "Status")
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []Row
	for i, r := range rows[1:] {
		row := Row{Line: i + 2}
		d := &row.Data
		d.DocumentNo = xl.Cell(r, idx, "Order Reference")
		d.VendorName = xl.Cell(r, idx, "Vendor")
		d.BuyerName = xl.Cell(r, idx, "Buyer")
		d.Description = xl.Cell(r, idx, "Description")
		d.DocumentStatus = xl.Cell(r, idx, "Status")
		if d.DocumentStatus == "" && d.VendorName == ""{
			continue
		}
		switch sd := xl.Cell(r, idx, "Source Document"); sd {
			case "Paid", "Unpaid":
				d.PaymentStatus = sd
			default:
				d.SourceDocument = sd
				row.Warnings = append(row.Warnings, "Payment Status Kosong")
		}

		if d.DocumentNo == "" {
			row.Errors = append(row.Errors, "Order Reference Kosong")
		} else if seen[d.DocumentNo] {
			row.Errors = append(row.Errors, "Order Reference Duplikat di File")
		}
		seen[d.DocumentNo] = true
		var e error
		if d.OrderDeadline, e = xl.Date(xl.Cell(r, idx, "Order Deadline")); e != nil {
			row.Errors = append(row.Errors, e.Error())
		}
		if d.DocumentTotal, e = xl.Money(xl.Cell(r, idx, "Total")); e != nil {
			row.Errors = append(row.Errors, e.Error())
		} else if d.DocumentTotal.LessThan(decimal.NewFromInt(1000)) {
			row.Warnings = append(row.Warnings, "Total Sangat Kecil")
		}
		if d.DocumentStatus != "Purchase Order" && d.DocumentStatus != "RFQ" {
			row.Errors = append(row.Errors, "Status Tidak Dikenal:"+d.DocumentStatus)
		}
		if d.VendorName == d.BuyerName {
			row.Warnings = append(row.Warnings, "Vendor dan Buyer Sama")
		}
		out = append(out, row)
	}
	return out, nil
}
