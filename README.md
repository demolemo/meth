# meth

smol CLI for tracking METrics related to Health. Every metric and value lives
in one sqlite file.

## Install

Needs Go 1.23+ and a C compiler. `go-sqlite3` is a cgo package, and with
`CGO_ENABLED=0` it compiles fine and then fails at runtime with *"Binary was
compiled with 'CGO_ENABLED=0' ... This is a stub"*, so the Makefile pins it.

```sh
git clone git@github.com:demolemo/meth.git
cd meth
make install
```

That builds `bin/meth` and copies it to `~/.local/bin/meth`. Somewhere else:

```sh
make install PREFIX=/usr/local      # -> /usr/local/bin/meth
make uninstall                      # same PREFIX you installed with
```

`make install` warns if the target directory is not on your `PATH`.

Straight from the module, without cloning:

```sh
go install github.com/demolemo/meth/cmd/meth@latest
```

That lands in `$(go env GOPATH)/bin`, which has to be on your `PATH` yourself.

Other targets: `make build`, `make test`, `make vet`, `make fmt`, `make clean`.

## Where the data lives

| what | path | override |
|------|------|----------|
| database | `~/.local/share/meth/meth.db` | `$XDG_DATA_HOME` |
| config | `~/.config/meth/config.toml` | `$XDG_CONFIG_HOME` |

The database file and its directory are created by the first command that
touches them, and `sqlite/schema.sql` is applied on every open with
`CREATE TABLE IF NOT EXISTS`. There is no migrate step and no setup SQL to run
by hand.

Database path, first hit wins:

1. `meth --db ./scratch.db ...`
2. `db = "..."` in the config file
3. the default above

```sh
meth config         # print the effective paths and whether the db exists yet
meth config init    # write the config file, refuses to clobber an existing one
```

An unknown key in the config file is an error rather than a silent default.

## Usage

```sh
meth metric create --id 1 --name weight
meth metric list
meth metric find -m weight
meth metric update --id 1 --name bodyweight
meth metric delete -m weight

meth value add -m weight -v 81.5                      # now
meth value add -m weight -v 81.5 --at "2026-10-08 09:30"
meth value list -m weight --limit 10                  # newest first, -1 for all
meth value update --id 3 -v 82.0
meth value delete --id 3
meth value import -m weight -f export.csv             # value,created_at columns
meth value import -m weight -f health-export.json     # Health Export Kit, bodyMass in kg
```

Every verb takes `--id` or `-m/--metric` to name the metric, never both.

A metric or value that is not there is an error with exit 1, never a
`rows affected: 0` that exits 0. Every lookup and every write reports a miss by
wrapping `meth.ErrNotFound`, so a scheduled importer can tell "nothing to do"
from "it worked".

`value import` is idempotent: a reading is identified by
`(metric_id, created_at)`, so re-importing an overlapping export changes
nothing unless the number itself changed. It reports `new / updated / already
there`.

## Layout

```
meth/            domain: Metric, MetricValue, the service interfaces. stdlib only
config/          where the db and config files live, TOML loading
parse/           file -> []*MetricValue. knows formats, knows nothing about storage
sqlite/          the store. schema.sql is embedded and applied on Open
cmd/meth/        the cobra CLI
```

`sqlite` asserts `var _ meth.MetricService = (*MetricService)(nil)` (and the
same for values), so renaming a method without updating the interface is a
build failure instead of silent drift.

---

# Next ideas
- [x] `list` - the verb that shows all metrics that exist inside of the service
- [x] `delete` - the verb that deletes metric by id and maybe the name (if we impose UNIQUE constraint on that bitch)
- [x] how should a process of adding a value work? `add-value` - CLI verb. Interface `meth add-value -name=[metric-name] -value`
- [x] add a convention for adding metric by `name`, names are case insensitive and cannot be empty
- [x] `import` from CSV. - landed as `value import -f`, which also reads a Health Export `.json`
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

