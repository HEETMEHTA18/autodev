package tui

// line is referenced by the legacy renderer while the professional TUI is
// migrated to use a local row variable. Keep this compatibility symbol until
// that migration is complete.
var line string

var _ = line
