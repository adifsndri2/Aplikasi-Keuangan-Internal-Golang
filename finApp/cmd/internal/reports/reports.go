package reports

type Table struct {
	Ttitle string `json:"title"`
	Label string `json:"label"`
	Columns []string `json:"columns"`
	Rows []map[string]any `json:"rows"`
	Notes []string `json:"notes"`
}

type Params struct {
	Start, End time.Time
	Opening decimal.Decimal
}

type Service struct {
	DB *gorm.DB
}
func New(db *gorm.DB) *Service {
	return &Service{DB: db}
}
type PeriodAmt struct {
	Period string
	Amt decimal.Decimal
}

func (s *Service) monthly(table, dateCol, ValCol, where string, p Params) (map[string]decimal.Decimal, error) {
	q := fmt.Sprintf("SELECT to_char(date_trunc('month', %s), 'YYYY-MM') AS period, "+ 
		"COALESCE(SUM(%s), 0) AS amt FROM %s WHERE %s AND %s >= ? AND %s < ? GROUP BY 1", 
		dateCol, ValCol, table, where, dateCol, dateCol)
		var res []PeriodAmt
		if err := s.DB.Raw(q, p.Start, p.End.AddDate(0, 0, 1)).Scan(&res).Error; err != nil {
			return nil, err
		}
			m := map[string]decimal.Decimal{}
		for _, r := range res {
			m[r.Period] = r.Amt
		}
		return m, nil
}	

func build(title, label string, cols []string, in, out []map[string]decimal.Decimal,
	opening decimal.Decimal, running bool, notes []string) Table {
		set := map[string]bool{}
		for k := range in {
			set[k] = true
		}
		for k := range out {
			set[k] = true
		}
		var periods []string
		for k := range set {
			periods = append(periods, k)
		}
		sort.Strings(periods)
		t := Table{Ttitle: title, Label: label, Columns: cols, Notes: notes}
		bal, totIn, totOut := opening, decimal.Zero, decimal.Zero
		for _, p := range periods {
			a, b := in[p], out[p]
			totIn, totOut = totIn.Add(a), totOut.Add(b)
			last := a.Sub(b)
			if running {
				bal = bal.Add(last)
				last = bal
			}
			t.Rows = append(t.Rows, []string{"TOTAL", totIn.StringFixed(2), totOut.StringFixed(2), totIn.Sub(totOut).StringFixed(2)})
		}
		return t
}

func (s *Service) Forecast(p Params) (Table, error) {
	in, err := s.monthly("ar_transactions", "due_date", "total_amount", "payment_status = 'Not Paid' AND is_internal = false", p)
	if err != nil {
		return Table{}, err
	}
	out, err := s.monthly("ap_transactions", "order_deadline", "document_total", "payment_status = 'Unpaid' AND document_status = 'Purchase Order'", p)
	if err != nil {
		return Table{}, err
	}
	return build("Laporan Forecasting", "ESTIMASI", []string{"Periode", "Estimasi Penerimaan", "Estimasi Pengeluaran", "Saldo Akhir"}, 
	in, out, p.Opening, true, []string{
		"Penerimaan dari AR Not Paid Berdasarkan due date: pengeluara dari AP Not Padi berdasarkan order deadline", 
		"Asumsi: Pembayaran terjadi tepat pada due date; RFQ dan partner internal dikecualikan"}), nil
}

func (s *Service) ProfitLoss(p Params) (Table, error) {
	in, err := s.monthly("ar_transactions", "invoice_due", "untaxed_amount", "is_internal = false", p)
	if err != nil {
		return Table{}, err
	}
	out, err := s.monthly("ap_transactions", "order_deadline", "document_total", "document_status = 'Purchase Order'", p)
	if err != nil {
		return Table{}, err
	}
	return build("Laporan Laba Rugi", "ESTIMASI", []string{"Periode", "Pendapatan", "Beban", "Laba/Rugi"}, 
	in, out, decimal.Zero, false, []string{
		"Estimasi operasional, bukan laba rugi akuntansi: belum ad pemetaan akun dan pengakuan periode."}), nil
}

func (s *Service) Cashflow(p Params) (Table, error) {
	in, err := s.monthly("ar_transactions", "due_date", "total_amount", "payment_status = 'Not Paid' AND is_internal = false", p)
	if err != nil {
		return Table{}, err
	}
	out, err := s.monthly("ap_transactions", "order_deadline", "document_total", "payment_status = 'Unpaid'", p)
	if err != nil {
		return Table{}, err
	}
	return build("Laporan Arus Kas", "PROVISIONAL", []string{"Periode", "Penerimaan Kas", "Pengeluaran Kas", "Saldo Kas"},
	in, out, p.Opening, true, []string{
		"Tanggal Pembayaran aktual belum tersedia; due date/order deadline dipakai sebagai pendekatan.",
		"Bukan arus kas aktual dan belum dikelompokkan ke aktivitas operasi/investasi/pendanaan."}), nil
}
