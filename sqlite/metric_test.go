package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/sqlite"
)

func TestCreateMetricReturnsID(t *testing.T) {
	cases := []struct {
		name string
		id   int
		want int64
	}{
		{name: "assigned", id: 0, want: 1},
		{name: "explicit", id: 42, want: 42},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, err := sqlite.Open(filepath.Join(t.TempDir(), "t.db"))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			ms := sqlite.NewMetricService(db)

			got, err := ms.CreateMetric(&meth.Metric{ID: c.id, Name: "weight"})
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Errorf("CreateMetric id = %d, want %d", got, c.want)
			}

			m, err := ms.GetMetricByID(int(got))
			if err != nil {
				t.Fatalf("GetMetricByID(%d): %v", got, err)
			}
			if m.Name != "weight" {
				t.Errorf("name = %q, want %q", m.Name, "weight")
			}
		})
	}
}
