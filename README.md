# Next ideas
- [x] `list` - the verb that shows all metrics that exist inside of the service
- [x] `delete` - the verb that deletes metric by id and maybe the name (if we impose UNIQUE constraint on that bitch)
- [x] how should a process of adding a value work? `add-value` - CLI verb. Interface `meth add-value -name=[metric-name] -value`
- [ ] `add-values` from CSV. - this should be just `add` verb with path to the file provided
- [x] structure subcommands using flags
- [ ] tests? fuzzy testing with AI or something.
- [x] `value list` - the command to display the recent values (`limit, id` - flags that needed for this command)
- [x] `value update` - what is the best interface for this? we want to update using the id?
- [ ] `value delete` - using the `id` or using the `date` (this should be easy), what happens when values with given ids don't exist anymore
    - by `id` done (`value delete --vid`). missing id is not an error, it prints `rows affected: 0`
    - by `date` not done yet: exact minute or whole day?
- [ ] graphs? - it would be nice to display the metric along the period of time
- [ ] web server to use the API - this is the most interesting part for me. because i don't completely understand how to serve the API over the endpoint. this looks messy for me.
- [ ] general idea includes stuff like automatic execution but i'm not sure that this belongs in the scope of this project


# Can it be built with the existing code?
| task | existing code? | what is new |
|---|---|---|
| `value delete` by date | yes, same as delete by id with another `WHERE` | decide: exact minute or whole day (whole day needs a range) |
| `add` from CSV | no | `os.Open`, `encoding/csv`, what to do when row 50 of 100 fails (transactions) |
| tests | no | `testing` package. store takes `*sql.DB`, so an in-memory sqlite works without refactoring |
| graphs | partly | bars with `strings.Repeat("#", n)` are just `Printf`. a real chart needs a library or html |
| web server | no | `net/http`, handlers, json. `store/` can be reused as is |
| automatic execution | no | launchd / cron |

# Open questions
- ~~what happens when we delete the metric with a certain id, do we move other metrics?~~ no, ids never move, the gap stays
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

