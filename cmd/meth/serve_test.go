package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/sqlite"
)

const bodyExport = `{
  "meta": {"timeZone": "Europe/Moscow"},
  "additional": {"body": {
    "units": {"bodyMass": "kg"},
    "daily": [
      {"date": "2026-10-07", "values": {"bodyMass": 80.8}},
      {"date": "2026-10-08", "values": {"bodyMass": 81}}
    ]}}
}`

// wires the package globals the handlers read, the same ones root.go fills in
// before a command runs
func serveTestDB(t *testing.T) {
	t.Helper()
	conn, err := sqlite.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	db = conn
	ms = sqlite.NewMetricService(conn)
	vs = sqlite.NewValueService(conn)
	if err := ms.CreateMetric(&meth.Metric{ID: 1, Name: "weight"}); err != nil {
		t.Fatal(err)
	}
}

func TestHandleImportRejects(t *testing.T) {
	serveTestDB(t)

	cases := []struct {
		name       string
		method     string
		target     string
		body       string
		wantStatus int
		wantErr    string
	}{
		{
			name: "no metric", method: http.MethodPost, target: "/import", body: bodyExport,
			wantStatus: http.StatusBadRequest, wantErr: "metric query parameter is required",
		},
		{
			name: "unknown metric", method: http.MethodPost, target: "/import?metric=hrv", body: bodyExport,
			wantStatus: http.StatusNotFound, wantErr: `metric "hrv"`,
		},
		{
			name: "not json", method: http.MethodPost, target: "/import?metric=weight", body: `{"meta":`,
			wantStatus: http.StatusBadRequest, wantErr: "request body",
		},
		{
			name: "wrong unit", method: http.MethodPost, target: "/import?metric=weight",
			body:       `{"meta":{},"additional":{"body":{"units":{"bodyMass":"lb"},"daily":[{"date":"2026-10-08","values":{"bodyMass":178}}]}}}`,
			wantStatus: http.StatusBadRequest, wantErr: "want kg",
		},
		{
			name: "wrong method", method: http.MethodGet, target: "/import?metric=weight",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			routes().ServeHTTP(rec, httptest.NewRequest(c.method, c.target, strings.NewReader(c.body)))

			if rec.Code != c.wantStatus {
				t.Fatalf("status %d, want %d: %s", rec.Code, c.wantStatus, rec.Body)
			}
			if c.wantErr == "" {
				return
			}
			var body struct{ Error string }
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decoding %s: %v", rec.Body, err)
			}
			if !strings.Contains(body.Error, c.wantErr) {
				t.Errorf("error %q, want it to contain %q", body.Error, c.wantErr)
			}
		})
	}
}

func TestHandleImportIsIdempotent(t *testing.T) {
	serveTestDB(t)

	post := func() importResponse {
		t.Helper()
		rec := httptest.NewRecorder()
		routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/import?metric=weight", strings.NewReader(bodyExport)))
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d: %s", rec.Code, rec.Body)
		}
		var res importResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatal(err)
		}
		return res
	}

	if got, want := post(), (importResponse{Metric: "weight", Inserted: 2}); got != want {
		t.Errorf("first import %+v, want %+v", got, want)
	}
	if got, want := post(), (importResponse{Metric: "weight", Unchanged: 2}); got != want {
		t.Errorf("second import %+v, want %+v", got, want)
	}

	mvs, err := vs.GetMetricValues(1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(mvs) != 2 {
		t.Fatalf("%d values stored, want 2", len(mvs))
	}
	if got := mvs[0].CreatedAt.Local().Format("2006-01-02 15:04 -07:00"); got != "2026-10-08 00:00 +03:00" {
		t.Errorf("newest value at %s", got)
	}
}
