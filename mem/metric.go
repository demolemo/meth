package mem

import (
	"errors"
	"math"
	"time"

	"github.com/demolemo/meth"
)

type MetricService struct {
	curID       int                  // current counter. removed metrics are not filled
	metrics     map[int]*meth.Metric // from numeric ID to the metric
	metricNames map[string]int       // from name to numeric ID
}

func NewMetricService() *MetricService {
	metrics := make(map[int]*meth.Metric)
	metricNames := make(map[string]int)
	return &MetricService{
		curID: 0, metrics: metrics, metricNames: metricNames,
	}
}

func (ms *MetricService) FindMetricByID(id int) (*meth.Metric, error) {
	met, ok := ms.metrics[id]
	if !ok {
		return nil, errors.New("Metric with a given ID have not been found")
	}
	return met, nil
}

func (ms *MetricService) FindMetricByName(name string) (*meth.Metric, error) {
	metID, ok := ms.metricNames[name]
	if !ok {
		return nil, errors.New("Metric with a given name have not been found")
	}
	met, _ := ms.metrics[metID]
	return met, nil
}

func (ms *MetricService) CreateMetric(met *meth.Metric) error {
	// add the ID for the metric
	id := ms.curID
	currentTime := time.Now()

	// setting all internal fields for the metric
	met.ID = id
	met.CreatedAt = currentTime
	met.UpdatedAt = currentTime

	// update metric storage
	ms.metrics[id] = met
	ms.curID += 1
	return nil
}

// this issue arises again here, because now we need to explicitly pass the interval from outside in the parts of the code
// that don't necessirily want to know about inner structure of intervals
func (ms *MetricService) GenerateMetricReport(id int, aggRule meth.AggRule, aggInterval meth.AggInterval) (*meth.MetricValueReport, error) {
	// NOTE: naming of interval variable could be more straightforward for easier LSP access
	met, err := ms.FindMetricByID(id)
	if err != nil {
		return nil, err
	}

	// decide a splitting interval duration
	var groupDuration time.Duration
	switch aggInterval {
	case meth.MinuteInterval:
		groupDuration = time.Minute
	case meth.HourInterval:
		groupDuration = time.Hour
	case meth.DayInterval:
		groupDuration = 24 * time.Hour
	default:
		return nil, errors.New("unknown agg interval was provided")
	}

	// split values to the buckets according to the splitting rule
	buckets := make(map[int64][]meth.MetricValue)
	for _, v := range *met.Values {
		key := v.CreatedAt.Truncate(groupDuration).Unix()
		buckets[key] = append(buckets[key], v)
	}

	// do we make a switch expression on the agg cases or we just delegate it to some function?
	var records []meth.MetricValueRecord
	switch aggRule {
	case meth.AggMax:
		var maxVal float32 = -1.0
		for intTimestamp, bucket := range buckets {
			// aggregate all values inside of one time bucket
			for _, v := range bucket {
				if v.Value > maxVal {
					maxVal = v.Value
				}
			}
			// append the timestamp converted back from int to unix
			// and aggregated max value
			records = append(records, meth.MetricValueRecord{
				Timestamp: time.Unix(intTimestamp, 0),
				Value:     maxVal,
			})
		}

	case meth.AggMin:
		var minVal float32 = float32(math.Inf(1))
		for intTimestamp, bucket := range buckets {
			// aggregate all values inside of one time bucket
			for _, v := range bucket {
				if v.Value < minVal {
					minVal = v.Value
				}
			}
			// append the timestamp converted back from int to unix
			// and aggregated max value
			records = append(records, meth.MetricValueRecord{
				Timestamp: time.Unix(intTimestamp, 0),
				Value:     minVal,
			})
		}
	case meth.AggAvg:
		var sumVal float32 = 0.0
		var cnt float32 = 0

		for intTimestamp, bucket := range buckets {
			// calculate the sum and count for average
			for _, v := range bucket {
				sumVal += v.Value
				cnt += 1.0
			}
			// calculate average value and store it
			avgVal := sumVal / cnt
			records = append(records, meth.MetricValueRecord{
				Timestamp: time.Unix(intTimestamp, 0),
				Value:     avgVal,
			})
		}
	default:
		return nil, errors.New("There is such an aggregation rule, use Min/Max/Avg rules")
	}
	return &meth.MetricValueReport{
		Name:        met.Name,
		AggRule:     aggRule,
		AggInterval: aggInterval,
		Records:     records,
	}, nil
}

// NOTE: for future self
// this is a particular quirk that i don't like from this architecture
// i have to constantly remember which fields are stored in the update requests
// as well as for some responses and shit like that. def a minus of storing types that way
func (ms *MetricService) UpdateMetric(id int, upd meth.MetricUpdate) (*meth.Metric, error) {
	met, err := ms.FindMetricByID(id)
	if err != nil {
		return nil, err
	}

	// for now only the name is updated
	met.Name = upd.Name
	return met, nil
}

func (ms *MetricService) DeleteMetricByID(id int) error {
	met, err := ms.FindMetricByID(id)
	if err != nil {
		return err
	}
	// delete metric from both storages
	delete(ms.metricNames, met.Name)
	delete(ms.metrics, id)
	return nil
}

func (ms *MetricService) DeleteMetricByName(name string) error {
	met, err := ms.FindMetricByName(name)
	if err != nil {
		return err
	}
	// delete metric from both storages
	delete(ms.metricNames, met.Name)
	delete(ms.metrics, met.ID)
	return nil
}
