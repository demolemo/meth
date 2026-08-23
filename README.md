# Main types 
`Metric` - represents some kind of an aggregate over the metric values (storage)
`MetricValue` - singular metric value (more or less single)
`MetricService` - some kind of an interface that updates the metric, `interface` not a `type`

For the first version of the program
- all metrics are held locally in the sqlite database
- only one user is present, just store the metrics and update them via the command line


## MetricValue
```go
type MetricValue struct {
    Value float // float value of the metric

    ID int // id field to identify this value

    // parent fields
    MetricID int // parent metric id value so that metric storage has 1to1 correspondence
    Metric *Metric // pointer to the parent metric, not sure why we store it here, i think for parent correspondence

    // timefields that belong to this metric
    CreatedAt timestamp // when was first this value created
    UpdatedAt timestamp // when was last this value updated

    // in the further updates we will add some kind of ownership over the metric
    }
```

# Metric
- do we need some fields that validate the underlying values like `minMetricValue` & `maxMetricValue`?
- do we need some fun aggregation rules for metric values?
```go
type Metric struct {
    ID int // field to identify this metric
    
    // timefields that are aggregate over the metric values, metric is not updated directly
    CreatedAt timestamp // when was the first MetricValue that belongs to this metric created
    UpdatedAt timestamp // when was the last MetricValue that belongs to this metric updated

    Name string // human readable name of the metric

    Value float // the last value of the metric (any aggregate of the metric would serve this, we can apply different filtering conditions here)

    Values *[]MetricValues // all of the values that belong to this metric
    }
```

# MetricService
```go
type MetricService interface {
    FindMetricByID(id int) (*Metric, error) // search metric by id
    FindMetricByName(name string) (*Metric, error) // search metric by name
    CreateMetric(metric *metric) error // validates the name of the metric 
    UpdateMetric(metric *metric, id int, upd MetricUpdate) error // updates some params inside of the metric
    DeleteMetricById(id int) error // tries to delete the metric, returns ERRNOTFOUND if there is no metric with the following id
    DeleteMetricByName(name string) error // tries to delete the metric, returns ERRNOTFOUND if there is no metric with the following name
    AverageMetricValues(metric *metric) error // average all underlying values for the metric
    AverageMetricValuesN(metric *metric, n int) error // average all underlying values for the metric - time aggregation comes in handy here too, maybe we should do something with it
    }

```

```
