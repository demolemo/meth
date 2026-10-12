package meth

import "errors"

// ErrNotFound is what a lookup returns when the row is not there, wrapped with
// which thing was missing: fmt.Errorf("metric %q: %w", name, meth.ErrNotFound).
// Callers test with errors.Is rather than comparing the result against nil, so
// a store can never report a miss by handing back a zero value nobody checks.
var ErrNotFound = errors.New("not found")
