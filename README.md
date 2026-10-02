# Next ideas
- [x] `list` - the verb that shows all metrics that exist inside of the service
- [x] `delete` - the verb that deletes metric by id and maybe the name (if we impose UNIQUE constraint on that bitch)
- [x] how should a process of adding a value work? `add-value` - CLI verb. Interface `meth add-value -name=[metric-name] -value`
- [x] add a convention for adding metric by `name`, names are case insensitive and cannot be empty
- [ ] `import` from CSV. - this should be just `add` verb with path to the file provided
- [x] structure subcommands using flags
- [ ] tests? fuzzy testing with AI or something.
- [x] `value list` - the command to display the recent values (`limit, id` - flags that needed for this command)
- [x] `value update` - what is the best interface for this? we want to update using the id?
- [x] `value delete` - using the `id` or using the `date` (this should be easy), what happens when values with given ids don't exist anymore - for now we delete values only by `id`, values with given ids just collapse. that leads to some values going not after each other.
- [ ] graphs? - it would be nice to display the metric along the period of time. this is additional functionality.
- [ ] web server to use the API - this is the most interesting part for me. because i don't completely understand how to serve the API over the endpoint. this looks messy for me.
- [ ] general idea includes stuff like automatic execution but i'm not sure that this belongs in the scope of this project. what did i mean here but automatic execution. i think it was 


# Open questions
- what happens when we delete the metric with a certain id, do we move other metrics?
- how do we attach values to metrics? by name prolly (name should be the main method of interaction, not `ID` field, too hard to keep track of)
- if we delete the metric, what happens to the values? should we delete the values first?

# Storing values and interacting with them
```go
type MetricValue struct {
    ID int
    MetricID int // or maybe store the pointer
    Value float
    CreatedAt time.Time // maybe store the date instead of datetime
    UpdatedAt time.Time
}
```

- **verb:** `add-value`
- **usage:** `add-value -id=[ID] -value=[value] -name=[name] -created-at=[datetime|string]` - we can pass metric either by the `ID` or by the `name`
- `add-values` with the following interface: `add-values -path=[filepath] -id=[ID] -name=[name]` and the csv should be one and only format that exists for now. CSV format: `value, created-at`

