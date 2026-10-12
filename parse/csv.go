package parse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/demolemo/meth"
)

var requiredCols = []string{"value", "created_at"}

func missingCols(cols map[string]int) []string {
	var missing []string
	for _, col := range requiredCols {
		if _, ok := cols[col]; !ok {
			missing = append(missing, col)
		}
	}
	return missing
}

// CSV parses path into values belonging to metricID. Columns are found by
// header name in any order, extra ones ignored, created_at in TimeLayout.
// One bad row fails the whole file.
func CSV(path string, metricID int) ([]*meth.MetricValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// not for trimming, TrimSpace below does that and both ends. this is what
	// makes `, "quoted, cell"` parse instead of failing on the bare quote
	r.TrimLeadingSpace = true

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	cols := make(map[string]int, len(header))
	for i, name := range header {
		cols[strings.TrimSpace(name)] = i
	}
	if missing := missingCols(cols); len(missing) > 0 {
		return nil, fmt.Errorf("%s: missing columns: %s", path, strings.Join(missing, ", "))
	}

	var mvs []*meth.MetricValue
	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		// a short row is ErrFieldCount, so the indexes below are in range
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		line, _ := r.FieldPos(cols["value"])

		value, err := strconv.ParseFloat(strings.TrimSpace(rec[cols["value"]]), 64)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: value: %w", path, line, err)
		}
		createdAt, err := ParseTime(strings.TrimSpace(rec[cols["created_at"]]))
		if err != nil {
			return nil, fmt.Errorf("%s:%d: created_at: %w", path, line, err)
		}

		mvs = append(mvs, &meth.MetricValue{
			MetricID:  metricID,
			Value:     value,
			CreatedAt: createdAt,
		})
	}
	return mvs, nil
}
