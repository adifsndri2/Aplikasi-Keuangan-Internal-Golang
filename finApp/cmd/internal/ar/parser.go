package ar

import (
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/finApp/cmd/internal/model"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/finApp/cmd/internal/xl"
)

type Row struct {
	Line     int
	Data     model.ARTransaction
	Errors   []string
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
	idx, err := xl.Header(rows,
		"Number",
		"Invoice Partner Display Name",
		"Tax Number",
		"Activities",
		"Status",
		"Payment Status",
		"Invoice/Bill Date",
		"Due Date",
		"Untaxed Amount Signed",
		"Total Signed",
		"Total in Currency Signed",
	)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var out []Row
	for i, r := range rows[1:] {
		row := Row{Line: i + 2}
		d := &row.Data
		d.DocumentNo = xl.Cell(r, idx, "Number")
		if d.DocumentNo == "" {
			continue
		}
		if !strings.HasPrefix(d.DocumentNo, "INV/") {
			row.Errors = append(row.Errors, "bukan customer invoice (prefix harus INV/)")
		}
		if seen[d.DocumentNo] {
			row.Errors = append(row.Errors, "nomor duplikat di file")
		}
		seen[d.DocumentNo] = true
		d.PartnerName = xl.Cell(r, idx, "Invoice Partner Display Name")
		d.IsInternal = strings.HasPrefix(strings.ToLower(d.PartnerName), "internal")
		if d.IsInternal {
			row.Warnings = append(row.Warnings, "partner internal, dikecualikan dari laporan")
		}
		d.TaxNumber = xl.Cell(r, idx, "Tax Number")
		d.Description = xl.Cell(r, idx, "Activities")
		d.JournalStatus = xl.Cell(r, idx, "Status")
		d.PaymentStatus = xl.Cell(r, idx, "Payment Status")
		if d.PaymentStatus != "Paid" && d.PaymentStatus != "Not Paid" {
			row.Errors = append(row.Errors, "payment status tidak dikenal")
		}
		var e error
		if d.InvoiceDate, e = xl.Date(xl.Cell(r, idx, "Invoice/Bill Date")); e != nil {
			row.Errors = append(row.Errors, e.Error())
		}
		if d.DueDate, e = xl.Date(xl.Cell(r, idx, "Due Date")); e != nil {
			row.Errors = append(row.Errors, e.Error())
		} else if d.DueDate.Before(d.InvoiceDate) {
			row.Warnings = append(row.Warnings, "due date sebelum tanggal invoice")
		}
		if d.UntaxedAmount, e = xl.Money(xl.Cell(r, idx, "Untaxed Amount Signed")); e != nil {
			row.Errors = append(row.Errors, "untaxed amount bukan angka")
		}
		if d.TotalAmount, e = xl.Money(xl.Cell(r, idx, "Total Signed")); e != nil {
			row.Errors = append(row.Errors, "total bukan angka")
		} else if d.TotalAmount.IsNegative() {
			row.Warnings = append(row.Warnings, "nilai negatif (credit note?)")
		}
		d.CurrencyTotal, _ = xl.Money(xl.Cell(r, idx, "Total in Currency Signed"))
		out = append(out, row)
	}
	return out, nil
}
