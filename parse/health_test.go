package parse

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

const bodyExport = `{
  "meta": {"timeZone": "Europe/Moscow", "schemaVersion": 2},
  "additional": {"body": {
    "units": {"bodyMass": "kg"},
    "daily": [
      {"date": "2026-10-07", "values": {"bodyMass": 80.8}},
      {"date": "2026-10-08", "values": {"bodyMass": 81}}
    ]}}
}`

func TestHealthExportWeightFrom(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		want    []string // "RFC3339=value" per reading, in order
		wantErr string
	}{
		{
			name: "two days land at midnight in the export's zone",
			body: bodyExport,
			want: []string{"2026-10-07T00:00:00+03:00=80.8", "2026-10-08T00:00:00+03:00=81"},
		},
		{
			name: "a day without bodyMass is skipped, not zero",
			body: `{"meta":{"timeZone":"UTC"},"additional":{"body":{"units":{"bodyMass":"kg"},
			        "daily":[{"date":"2026-10-07","values":{}},{"date":"2026-10-08","values":{"bodyMass":81}}]}}}`,
			want: []string{"2026-10-08T00:00:00Z=81"},
		},
		{
			name:    "no body category",
			body:    `{"meta":{"timeZone":"UTC"},"additional":{}}`,
			wantErr: "no additional.body.daily rows",
		},
		{
			name: "pounds would be stored as kilograms",
			body: `{"meta":{"timeZone":"UTC"},"additional":{"body":{"units":{"bodyMass":"lb"},
			        "daily":[{"date":"2026-10-08","values":{"bodyMass":178}}]}}}`,
			wantErr: `bodyMass unit is "lb", want kg`,
		},
		{
			name: "unknown zone",
			body: `{"meta":{"timeZone":"Mars/Olympus"},"additional":{"body":{"units":{"bodyMass":"kg"},
			        "daily":[{"date":"2026-10-08","values":{"bodyMass":81}}]}}}`,
			wantErr: "meta.timeZone",
		},
		{
			name: "unparseable date",
			body: `{"meta":{"timeZone":"UTC"},"additional":{"body":{"units":{"bodyMass":"kg"},
			        "daily":[{"date":"08.10.2026","values":{"bodyMass":81}}]}}}`,
			wantErr: `date "08.10.2026"`,
		},
		{
			name:    "not json",
			body:    `{"meta":`,
			wantErr: "unexpected end of JSON input",
		},
		{
			name:    "trailing garbage after the object",
			body:    bodyExport + `}`,
			wantErr: "invalid character",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mvs, err := HealthExportWeightFrom(strings.NewReader(c.body), "test", 7)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("want error containing %q, got %d values", c.wantErr, len(mvs))
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("want error containing %q, got %v", c.wantErr, err)
				}
				if !strings.HasPrefix(err.Error(), "test: ") {
					t.Errorf("error should name its source: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(mvs) != len(c.want) {
				t.Fatalf("got %d values, want %d", len(mvs), len(c.want))
			}
			for i, want := range c.want {
				got := mvs[i].CreatedAt.Format(time.RFC3339) + "=" + strconv.FormatFloat(mvs[i].Value, 'g', -1, 64)
				if got != want {
					t.Errorf("value %d: got %s, want %s", i, got, want)
				}
				if mvs[i].MetricID != 7 {
					t.Errorf("value %d: metric id %d, want 7", i, mvs[i].MetricID)
				}
			}
		})
	}
}
