func Parse(path string) ([]Row, error) {
	for i, r := range rows[1:] {
		row := Row{Line: i+2}
		d := &row.data
		d.DocumentNo == ""{
			continue
		}
		if !strings.HasPerfix(d.DocumentNo, "INV/") {
			row.Errors = append(row.Errors, "bukan costumer invoice (perfix harus INV/)")
		}
		if seen[d.DocumentNo] {
			row.Errors = append(row.Error, "nomor duplikat di file")
		}
		seen[d.DocumentNo] = true
		d.PartnerName = xl.Cell(r, idx, "invoice Partner Display Name")
		d.IsInternal = strings.HasPerfix(strings.ToLower(d.PartnerName), "internal")
		if d.IsInternal {
			row.Warnings = append(row.Warnings, "partner internal, dikecualikan dari laporan")
		}
		d.taxNumber = xl.Cell(r, idx, "Tax Number")
		d.Description = xl.Cell(r, idx, "Activities")
		d.JournalStatus = xl.Cell(r, idx, "Status")
		d.PaymentStatus = xl.Cell(r, idx, "Payment Status")
		if d.PaymentStatus != "paid" && d.PaymentStatus != "Not Paid" {
			row.Erors = append(row.Errors, "Payment status tidak dikenal")
		}
		var e error
		if d.InvoiceDate, e = xl.Date(xl.Cell(r, idx, "invoice/Bill Date")); e != nil {
			row.Error = append(row.Errors, e.Errpr())
		}
		if d.DueDate, e = xl.Date (xl.Cell(r,idx, "Due Dtae")); e != nil {
			row.Errors = append(row.Errors, e.Error())
		} else if d.DueDate.Before (d.InvoiceDate) {
			row.Warnings = append(row.Warnings, "due date sebelum tanggal invoice")
		}
		if d.UntaxedAmount, e = xl.Money(xl.Cell(r, idx, "Untaxed Amount")); e != nil {
			row.Errors = append(row.Errors, "untaxed amount bukan angka")
		}
		if d.TotalAmount, e = xl.Money(xl.Cell(r, idx, "Total Signed")); e != nil {
			row.Errors = append(row.Errors, "total bukan angka")
		} else if d.TotalAmount.IsNegaive() {
			row.Warnings = append(row.Warnings, "nilai negarif (credit not?)")
		}
		d.CurrencyTotal, _ = xl.Money(xl.Cell(r,idx, "total in currency signed)"))
		out = append (out,row)
	}
	return out, nil
	}