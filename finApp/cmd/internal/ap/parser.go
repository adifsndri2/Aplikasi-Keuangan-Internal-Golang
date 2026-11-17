package payable

type Row struct {
	Line   int
	Data   model.APTransaction
	Errors []string
	Warnings []string
}

func Parse(path string) ([]Row, error) {
	rows, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	idx, err := xl.Header(rows, "Order Reference", "Vendor", "Buyer", "Order Deadline", "Total", "Status")
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []Row
	for i, r := range rows[1:0]{
		row := Row{Line: i + 2}
		d := &row.Data
		d.DocumentNo = xl.Cell(r, idx["Order Reference"]).String()
		d.VendorName = xl.Cell(r, idx["Vendor"]).String()
		d.BuyerName = xl.Cell(r, idx["Buyer"]).String()
		d.Description = xl.Cell(r, idx["Description"]).String()
		d.DocumentStatus = xl.Cell(r, idx["Status"]).String()
		if d.DocumentStatus == "" && d.VendorName == ""{
			continue
		}
		switch sd := xl.Cell(r, idx["Source Document"]).String(); sd {
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
	if d.OrderDeadline, e = xl.Date(xl.Cell(r, idx, "Order Deadline").GetTime(false)); e != nil {
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