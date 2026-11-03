package parse

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/demolemo/meth"
)

const dateLayout = "2006-01-02"

// healthExport is the weight-carrying part of a Health Export Kit file
// (meta.schemaVersion 2). activity, sleep sessions, workouts and the sample
// streams are other shapes and are ignored on purpose.
type healthExport struct {
	Meta struct {
		TimeZone string `json:"timeZone"`
	} `json:"meta"`
	Additional struct {
		Body struct {
			Units map[string]string `json:"units"`
			Daily []struct {
				Date   string             `json:"date"`
				Values map[string]float64 `json:"values"`
			} `json:"daily"`
		} `json:"body"`
	} `json:"additional"`
}

// HealthExportWeight reads bodyMass out of a Health Export Kit file as values
// of metricID. The file gives a date and no clock time, so a reading lands at
// midnight in the export's own zone, not the machine's.
func HealthExportWeight(path string, metricID int) ([]*meth.MetricValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return HealthExportWeightFrom(f, path, metricID)
}

// HealthExportWeightFrom is HealthExportWeight over an open export. source
// only names it in error messages, so a request body can say so.
func HealthExportWeightFrom(r io.Reader, source string, metricID int) ([]*meth.MetricValue, error) {
	// read then unmarshal, not a streaming decode: Unmarshal rejects trailing
	// garbage after the object and the decoder accepts it
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var export healthExport
	if err := json.Unmarshal(raw, &export); err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}

	body := export.Additional.Body
	if len(body.Daily) == 0 {
		return nil, fmt.Errorf("%s: no additional.body.daily rows, was the export made with the body category?", source)
	}
	// value is a bare float in the db, so pounds would be stored as kilograms
	if unit := body.Units["bodyMass"]; unit != "kg" {
		return nil, fmt.Errorf("%s: bodyMass unit is %q, want kg", source, unit)
	}

	loc := time.Local
	if export.Meta.TimeZone != "" {
		loc, err = time.LoadLocation(export.Meta.TimeZone)
		if err != nil {
			return nil, fmt.Errorf("%s: meta.timeZone: %w", source, err)
		}
	}

	var mvs []*meth.MetricValue
	for _, day := range body.Daily {
		value, ok := day.Values["bodyMass"]
		if !ok {
			continue
		}
		createdAt, err := time.ParseInLocation(dateLayout, day.Date, loc)
		if err != nil {
			return nil, fmt.Errorf("%s: date %q: %w", source, day.Date, err)
		}
		mvs = append(mvs, &meth.MetricValue{
			MetricID:  metricID,
			Value:     value,
			CreatedAt: createdAt,
		})
	}
	return mvs, nil
}
