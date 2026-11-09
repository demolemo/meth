package tests

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/sqlite"
)

func TestNotFoundIsCatchable(t *testing.T) {
	db, err := sqlite.Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ms := sqlite.NewMetricService(db)
	vs := sqlite.NewValueService(db)

	cases := map[string]func() error{
		"GetMetricByID":      func() error { _, err := ms.GetMetricByID(999); return err },
		"GetMetricByName":    func() error { _, err := ms.GetMetricByName("nope"); return err },
		"ResolveMetricID/id": func() error { _, err := ms.ResolveMetricID(999, ""); return err },
		"ResolveMetricID/nm": func() error { _, err := ms.ResolveMetricID(0, "nope"); return err },
		"UpdateMetricByID":   func() error { return ms.UpdateMetricByID(&meth.MetricUpdate{ID: 999, Name: "x"}) },
		"DeleteMetricByID":   func() error { return ms.DeleteMetricByID(999) },
		"DeleteMetricByName": func() error { return ms.DeleteMetricByName("nope") },
		"GetMetricValue":     func() error { _, err := vs.GetMetricValue(999); return err },
		"UpdateMetricValue":  func() error { _, err := vs.UpdateMetricValue(&meth.MetricValue{ID: 999}); return err },
		"DeleteMetricValue":  func() error { _, err := vs.DeleteMetricValue(999); return err },
	}

	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			err := call()
			if err == nil {
				t.Fatal("want an error, got nil")
			}
			if !errors.Is(err, meth.ErrNotFound) {
				t.Fatalf("errors.Is(err, meth.ErrNotFound) = false, err = %v", err)
			}
			t.Logf("ok: %v", err)
		})
	}
}
