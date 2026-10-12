package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/demolemo/meth"
)

// Reading is a number that knows the name and unit of the metric it belongs
// to but not its id. A file holding one metric parses straight to
// []*meth.MetricValue; a file that invents two dozen metrics cannot, because
// the ids do not exist until something writes them.
type Reading struct {
	Metric string
	Unit   string
	Value  float64
	At     time.Time
}

type whoopDump struct {
	Source   string          `json:"source"`
	Sleeps   []whoopSleep    `json:"sleeps"`
	Recovery []whoopRecovery `json:"recovery"`
}

// Every number is a pointer so a missing one stays missing. spo2 and
// skin_temp_c are null in real dumps and a null that becomes 0 is a 0% blood
// oxygen reading.
type whoopSleep struct {
	ID           string   `json:"id"`
	Day          string   `json:"day"`
	EndUTC       string   `json:"end_utc"`
	Nap          int      `json:"nap"`
	InBedH       *float64 `json:"in_bed_h"`
	LightH       *float64 `json:"light_h"`
	SWSH         *float64 `json:"sws_h"`
	REMH         *float64 `json:"rem_h"`
	AwakeH       *float64 `json:"awake_h"`
	Efficiency   *float64 `json:"efficiency"`
	PerfPct      *float64 `json:"perf_pct"`
	Consistency  *float64 `json:"consistency"`
	Disturbances *float64 `json:"disturbances"`
	RespRate     *float64 `json:"resp_rate"`
}

type whoopRecovery struct {
	Day         string   `json:"day"`
	SleepID     string   `json:"sleep_id"`
	RecoveryPct *float64 `json:"recovery_pct"`
	RHR         *float64 `json:"rhr"`
	HRVMilli    *float64 `json:"hrv_milli"`
	SpO2        *float64 `json:"spo2"`
	SkinTempC   *float64 `json:"skin_temp_c"`
}

// Leaf names are appended to "sleep." for a night and "nap." for a nap: two
// sessions of one day are two series, not one series fighting over a slot.
var sleepFields = []struct {
	leaf string
	unit string
	pick func(whoopSleep) *float64
}{
	{"in_bed", "h", func(s whoopSleep) *float64 { return s.InBedH }},
	{"light", "h", func(s whoopSleep) *float64 { return s.LightH }},
	{"sws", "h", func(s whoopSleep) *float64 { return s.SWSH }},
	{"rem", "h", func(s whoopSleep) *float64 { return s.REMH }},
	{"awake", "h", func(s whoopSleep) *float64 { return s.AwakeH }},
	{"efficiency", "%", func(s whoopSleep) *float64 { return s.Efficiency }},
	{"performance", "%", func(s whoopSleep) *float64 { return s.PerfPct }},
	{"consistency", "%", func(s whoopSleep) *float64 { return s.Consistency }},
	{"disturbances", "count", func(s whoopSleep) *float64 { return s.Disturbances }},
	{"resp_rate", "brpm", func(s whoopSleep) *float64 { return s.RespRate }},
}

var recoveryFields = []struct {
	leaf string
	unit string
	pick func(whoopRecovery) *float64
}{
	{"score", "%", func(r whoopRecovery) *float64 { return r.RecoveryPct }},
	{"rhr", "bpm", func(r whoopRecovery) *float64 { return r.RHR }},
	{"hrv", "ms", func(r whoopRecovery) *float64 { return r.HRVMilli }},
	{"spo2", "%", func(r whoopRecovery) *float64 { return r.SpO2 }},
	{"skin_temp", "degC", func(r whoopRecovery) *float64 { return r.SkinTempC }},
}

// Whoop reads a whoop dump into readings under the sleep., nap. and recovery.
// families.
//
// Every value of a session is anchored at end_utc, the instant of waking.
// WHOOP's `day` is the local wake date, so the wake instant is the only one
// that always falls on it: of the 425 sessions in the reference dump 19 began
// on the day before the one they are filed under. Anchoring on the date alone
// would be worse than imprecise - 41 days hold more than one session and 9
// hold more than one recovery, and each of those pairs would collapse onto a
// single (metric, created_at) slot with the later row winning silently.
//
// A recovery is scored for the sleep it follows and carries no clock of its
// own, so it borrows that sleep's anchor.
func Whoop(path string) ([]Reading, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dump whoopDump
	if err := json.Unmarshal(raw, &dump); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if dump.Source != "" && dump.Source != "whoop" {
		return nil, fmt.Errorf("%s: source is %q, not whoop", path, dump.Source)
	}
	if len(dump.Sleeps) == 0 && len(dump.Recovery) == 0 {
		return nil, fmt.Errorf("%s: no sleeps and no recovery rows", path)
	}

	readings := make([]Reading, 0, len(dump.Sleeps)*len(sleepFields)+len(dump.Recovery)*len(recoveryFields))
	anchors := make(map[string]time.Time, len(dump.Sleeps))

	for _, s := range dump.Sleeps {
		at, err := time.Parse(time.RFC3339, s.EndUTC)
		if err != nil {
			return nil, fmt.Errorf("%s: sleep %s: end_utc %q: %w", path, s.ID, s.EndUTC, err)
		}
		anchors[s.ID] = at

		family := "sleep."
		if s.Nap != 0 {
			family = "nap."
		}
		for _, f := range sleepFields {
			if v := f.pick(s); v != nil {
				readings = append(readings, Reading{Metric: family + f.leaf, Unit: f.unit, Value: *v, At: at})
			}
		}
	}

	for _, r := range dump.Recovery {
		at, ok := anchors[r.SleepID]
		if !ok {
			return nil, fmt.Errorf("%s: recovery for %s: sleep %s is not in this dump", path, r.Day, r.SleepID)
		}
		for _, f := range recoveryFields {
			if v := f.pick(r); v != nil {
				readings = append(readings, Reading{Metric: "recovery." + f.leaf, Unit: f.unit, Value: *v, At: at})
			}
		}
	}
	return readings, nil
}

// Values turns readings into storable values once ids exist. resolve is
// called once per distinct name, so an importer of 6000 readings creates its
// two dozen metrics and no more.
func Values(readings []Reading, resolve func(name, unit string) (int, error)) ([]*meth.MetricValue, error) {
	ids := make(map[string]int)
	mvs := make([]*meth.MetricValue, 0, len(readings))
	for _, r := range readings {
		id, ok := ids[r.Metric]
		if !ok {
			var err error
			if id, err = resolve(r.Metric, r.Unit); err != nil {
				return nil, err
			}
			ids[r.Metric] = id
		}
		mvs = append(mvs, &meth.MetricValue{MetricID: id, Value: r.Value, CreatedAt: r.At})
	}
	return mvs, nil
}
