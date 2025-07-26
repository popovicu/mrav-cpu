package system

import (
	"log/slog"
)

type SystemOpts struct {
	Logger  *slog.Logger
	Verbose bool
}
