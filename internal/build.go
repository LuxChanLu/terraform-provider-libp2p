package internal

import (
	_ "embed"
)

var BuildCommit string  // nolint:gochecknoglobals // Filled at compilation time.
var BuildVersion string // nolint:gochecknoglobals // Filled at compilation time.
