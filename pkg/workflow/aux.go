package workflow

import "context"

type (
	final struct{}
)

var (
	_ Finalizer = &final{}
)

func Final() *final                                            { return &final{} }
func (t *final) Finalize(_ context.Context, s Variables) error { return nil }
func (final) NodeRef() string                                  { return "(FINAL)" }
