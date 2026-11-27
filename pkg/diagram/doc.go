// Package diagram builds and lays out a render-independent Scene from a
// terraform.Plan.
//
// The Scene holds containers (providers and module instances), nodes (one per
// drawn resource instance), dependency edges and integer coordinates. It imports
// pkg/terraform only and has no SVG or XML knowledge, so several renderers can
// share one layout.
package diagram
