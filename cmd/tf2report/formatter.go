package main

import (
	"fmt"
	"io"

	"github.com/germainlefebvre4/tf2report/pkg/config"
	"github.com/germainlefebvre4/tf2report/pkg/report"
)

// newFormatter builds the formatter for the resolved format. The svg format needs
// to know which filters were given and where the optional icon directory is; every
// other format goes through report.NewFormatter. With a format other than svg,
// --icons-dir is ignored and, in verbose mode, a notice says so.
func newFormatter(format report.Format, cfg *config.Config, iconsDir string, verbose bool, stderr io.Writer) (report.Formatter, error) {
	if format != report.FormatSVG {
		if iconsDir != "" && verbose {
			// A failed write to standard error has no recovery.
			_, _ = fmt.Fprintln(stderr, "--icons-dir is ignored: it only applies to the svg format")
		}
		return report.NewFormatter(format)
	}
	opts := report.SVGOptions{
		IconsDir:       iconsDir,
		TypeFiltered:   len(cfg.Filters.ResourceTypes) > 0,
		ActionFiltered: len(cfg.Filters.Actions) > 0,
	}
	if verbose {
		opts.OnSkip = func(file, reason string) {
			_, _ = fmt.Fprintf(stderr, "skipped icon %s: %s\n", file, reason)
		}
	}
	return report.NewSVGFormatter(opts), nil
}
