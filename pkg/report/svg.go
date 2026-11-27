package report

import (
	"fmt"
	"io"

	"github.com/germainlefebvre4/tf2report/pkg/diagram"
	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// SVGOptions configures the SVG formatter.
type SVGOptions struct {
	// IconsDir is an optional directory of icons; empty means built-in glyphs only.
	IconsDir string
	// TypeFiltered and ActionFiltered tell whether the user gave a type or an
	// action filter; they decide no-op visibility and dependency context.
	TypeFiltered   bool
	ActionFiltered bool
	// OnSkip, when set, is called for each icon file that is skipped. The
	// formatter itself never writes to standard error.
	OnSkip func(file, reason string)
}

// SVGFormatter draws the resources of a plan as an SVG architecture diagram.
type SVGFormatter struct {
	opts SVGOptions
}

// NewSVGFormatter creates an SVG formatter with the given options.
func NewSVGFormatter(opts SVGOptions) *SVGFormatter {
	return &SVGFormatter{opts: opts}
}

// Format builds, lays out and renders the Scene, then writes the document once.
// The summary is not used. An invalid icons directory fails before any output.
func (f *SVGFormatter) Format(plan *terraform.Plan, _ *terraform.Summary, w io.Writer) error {
	icons, err := newIconSet(f.opts.IconsDir, f.opts.OnSkip)
	if err != nil {
		return err
	}
	scene, err := diagram.Build(plan, diagram.Options{
		TypeFiltered:   f.opts.TypeFiltered,
		ActionFiltered: f.opts.ActionFiltered,
	})
	if err != nil {
		return fmt.Errorf("failed to build diagram: %w", err)
	}
	diagram.Layout(scene)

	doc := renderScene(scene, icons)
	if _, err := io.WriteString(w, doc); err != nil {
		return fmt.Errorf("failed to write svg: %w", err)
	}
	return nil
}
