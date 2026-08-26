# Product idea for the application
There are many metrics that I'm tracking currently - weight, my PRs in the gym, my sleep score etc.
It would be cool to store them in one place and share them with friends. Also friends inputting the metrics would be fun too.
Metric is just a number with a date that this number was generated - that's about it.


# Main types 
`Metric` - aggregate over the metric values. values could be aggregated in different ways. Max value, min value, average (storage)
`MetricValue` - singular metric value. Value of the metric in concrete point in time.
`MetricService` - the `interface` that tells us how to interact with a metric, everything that implements it could work with a metric.

# Stages of the program that i see currently
## MVP
MVP - input one metric that is weight and aggregate weight inside of the day. Representation layer is just a table with avg value of the weight in one day.

## Further work
For now we are doing the MVP and not thinking too much inside what would be possible in the further versions of the program.

# MVP types
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
    UpdatedAt timestamp // when was last this value updated - should the value be updated? of course it should, for example we input an incorrect value by mistake.
    }
```

## Metric
```go
type Metric struct {
    ID int // field to identify this metric

    // timefields that are aggregate over the metric values, metric is not updated directly
    CreatedAt timestamp // when was the first MetricValue that belongs to this metric created
    UpdatedAt timestamp // when was the last MetricValue that belongs to this metric updated, inferred, never input directly

    Name string // human readable name of the metric

    AggValues *[]MetricValues // values of the metric aggregated in some manner
    Values *[]MetricValues // all of the values that belong to this metric

    // ?? not sure maybe i want to store the rule to agg near the metric and create some predefined rules of aggregation

    AggRule AggRule // aggregation function to apply to the metric values
    AggInterval AggInterval // aggregation interval to use when aggregating metric values
    }

// supporting types
type AggRule int // aggregation rule to apply to the metric

const (
    AggMax AggRule = iota
    AggMin
    AggAvg
)

type AggInterval int // aggregation interval that applies to the metric

const (
    DayInterval AggInterval = iota
    HourInterval
    MinuteInterval
)
```

## MetricService
```go
type MetricService interface {
    FindMetricByID(id int) (*Metric, error) // search metric by id, returns error if metric is not found
    FindMetricByName(name string) (*Metric, error) // search metric by name, returns error if metric is not found
    CreateMetric(metric *Metric) error // validates the name of the metric 
    UpdateMetric(id int, upd MetricUpdate) error // updates some params inside of the metric
    DeleteMetricById(id int) error // tries to delete the metric, returns ERRNOTFOUND if there is no metric with the following id
    DeleteMetricByName(name string) error // tries to delete the metric, returns ERRNOTFOUND if there is no metric with the following name
    }
```


## MetricUpdate
```go
type MetricUpdate struct {
        Name // New name for the metric
        AggRule AggRule // New aggregation function to use for the metric
        AggInterval AggInterval // New aggregation interval to use with the metric
    }
```
