package parse

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// one night, one nap on the same day, and a recovery for each
const whoopFixture = `{
  "source": "whoop",
  "sleeps": [
    {"id": "night", "day": "2026-10-12", "start_utc": "2026-10-11T22:56:48.180Z",
     "end_utc": "2026-10-12T07:00:54.840Z", "tz": "+03:00", "nap": 0,
     "in_bed_h": 8.869, "light_h": 3.159, "sws_h": 2.625, "rem_h": 2.167, "awake_h": 0.918,
     "efficiency": 89.90967, "perf_pct": 82.0, "consistency": 70.0, "disturbances": 11,
     "resp_rate": 14.707031, "score_state": "SCORED"},
    {"id": "doze", "day": "2026-10-12", "start_utc": "2026-10-12T12:00:00.000Z",
     "end_utc": "2026-10-12T13:30:00.000Z", "tz": "+03:00", "nap": 1,
     "in_bed_h": 1.5, "light_h": 1.0, "sws_h": 0.4, "rem_h": 0.0, "awake_h": 0.1,
     "efficiency": 93.0, "perf_pct": 10.0, "consistency": 70.0, "disturbances": 1,
     "resp_rate": 14.1, "score_state": "SCORED"}
  ],
  "recovery": [
    {"day": "2026-10-12", "cycle_id": 1, "sleep_id": "night", "recovery_pct": 87.0,
     "rhr": 67.0, "hrv_milli": 55.20474, "spo2": 95.611115, "skin_temp_c": 34.232666,
     "calibrating": 0},
    {"day": "2026-10-12", "cycle_id": 2, "sleep_id": "doze", "recovery_pct": 40.0,
     "rhr": 70.0, "hrv_milli": 30.5, "spo2": null, "skin_temp_c": null, "calibrating": 0}
  ]
}`

func writeFixture(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "whoop-dump.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestWhoopMapping(t *testing.T) {
	readings, err := Whoop(writeFixture(t, whoopFixture))
	if err != nil {
		t.Fatalf("Whoop: %v", err)
	}

	night := time.Date(2026, 10, 12, 7, 0, 54, 840000000, time.UTC)
	doze := time.Date(2026, 10, 12, 13, 30, 0, 0, time.UTC)

	cases := []struct {
		metric string
		at     time.Time
		unit   string
		value  float64
	}{
		{"sleep.in_bed", night, "h", 8.869},
		{"sleep.light", night, "h", 3.159},
		{"sleep.sws", night, "h", 2.625},
		{"sleep.rem", night, "h", 2.167},
		{"sleep.awake", night, "h", 0.918},
		{"sleep.efficiency", night, "%", 89.90967},
		{"sleep.performance", night, "%", 82.0},
		{"sleep.consistency", night, "%", 70.0},
		{"sleep.disturbances", night, "count", 11},
		{"sleep.resp_rate", night, "brpm", 14.707031},

		// nap=1 goes to its own family, anchored at its own wake instant, so
		// it never lands on the night's slot
		{"nap.in_bed", doze, "h", 1.5},
		{"nap.rem", doze, "h", 0.0},
		{"nap.disturbances", doze, "count", 1},

		// a recovery borrows the anchor of the sleep it scored
		{"recovery.score", night, "%", 87.0},
		{"recovery.rhr", night, "bpm", 67.0},
		{"recovery.hrv", night, "ms", 55.20474},
		{"recovery.spo2", night, "%", 95.611115},
		{"recovery.skin_temp", night, "degC", 34.232666},
		{"recovery.score", doze, "%", 40.0},
	}

	index := make(map[string]Reading, len(readings))
	for _, r := range readings {
		key := r.Metric + "@" + r.At.Format(time.RFC3339Nano)
		if prev, dup := index[key]; dup {
			t.Fatalf("%s written twice, %g then %g", key, prev.Value, r.Value)
		}
		index[key] = r
	}

	for _, c := range cases {
		key := c.metric + "@" + c.at.Format(time.RFC3339Nano)
		got, ok := index[key]
		if !ok {
			t.Errorf("%s: missing", key)
			continue
		}
		if got.Value != c.value {
			t.Errorf("%s: value %g, want %g", key, got.Value, c.value)
		}
		if got.Unit != c.unit {
			t.Errorf("%s: unit %q, want %q", key, got.Unit, c.unit)
		}
	}

	// 2 sessions x 10 fields + 5 recovery fields + 3 of 5 for the row whose
	// spo2 and skin_temp_c are null
	if want := 28; len(readings) != want {
		t.Errorf("%d readings, want %d", len(readings), want)
	}
	for _, r := range readings {
		if r.Metric == "recovery.spo2" && r.At.Equal(doze) {
			t.Error("null spo2 produced a reading")
		}
		if r.Metric == "recovery.skin_temp" && r.At.Equal(doze) {
			t.Error("null skin_temp_c produced a reading")
		}
	}
}

func TestWhoopErrors(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"not a whoop dump", `{"source": "garmin", "sleeps": [], "recovery": []}`},
		{"empty", `{"source": "whoop", "sleeps": [], "recovery": []}`},
		{"recovery without its sleep", `{"source": "whoop", "sleeps": [], "recovery": [
			{"day": "2026-10-12", "sleep_id": "gone", "recovery_pct": 87.0}]}`},
		{"unparseable end_utc", `{"source": "whoop", "sleeps": [
			{"id": "x", "day": "2026-10-12", "end_utc": "yesterday", "nap": 0, "rem_h": 1.0}]}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := Whoop(writeFixture(t, c.body)); err == nil {
				t.Error("want an error, got none")
			}
		})
	}
}

func TestValuesResolvesEachNameOnce(t *testing.T) {
	readings, err := Whoop(writeFixture(t, whoopFixture))
	if err != nil {
		t.Fatalf("Whoop: %v", err)
	}

	calls := make(map[string]int)
	ids := make(map[string]int)
	mvs, err := Values(readings, func(name, unit string) (int, error) {
		calls[name]++
		if _, ok := ids[name]; !ok {
			ids[name] = len(ids) + 1
		}
		return ids[name], nil
	})
	if err != nil {
		t.Fatalf("Values: %v", err)
	}
	if len(mvs) != len(readings) {
		t.Errorf("%d values from %d readings", len(mvs), len(readings))
	}
	for name, n := range calls {
		if n != 1 {
			t.Errorf("%s resolved %d times", name, n)
		}
	}
	// 10 leaves under sleep., 10 under nap., 5 under recovery.
	if want := 25; len(calls) != want {
		t.Errorf("%d distinct metrics, want %d", len(calls), want)
	}
}
