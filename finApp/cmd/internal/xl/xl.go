package xl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

func Open(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.GetRows(f.GetSheetName(0), excelize.Options{RawCellValue: true})
}

func Header(rows [][]string, required ...string) (map[string]int, error) {
	if len(rows) == 0 {
		return nil, errors.New("sheet kosong")
	}
	idx := map[string]int{}
	for i, h := range rows[0] {
		idx[strings.TrimSpace(h)] = i
	}

	var missing []string
	for _, r := range required {
		if _, ok := idx[r]; !ok {
			missing = append(missing, r)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("header tidak lengkap, kolom yang hilang: %s", strings.Join(missing, ", "))
	}

	return idx, nil
}

func Cell(row []string, idx map[string]int, name string) string {
	i, ok := idx[name]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

func Money(s string) (decimal.Decimal, error) {
	s = strings.NewReplacer("Rp", "", " ", "").Replace(s)
	return decimal.NewFromString(s)
}

func Date(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return excelize.ExcelDateToTime(f, false)
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02", "02/01/2006 15:04:05", "02/01/2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("tanggal tidak valid: %q", s)
}
