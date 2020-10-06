package yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"gopkg.in/yaml.v3"
	"io"
)

type (
	root struct {
		Scope expressions
		Steps steps
	}
)

func Load(r io.Reader) (workflow.Node, workflow.Variables, error) {
	var (
		d   = yaml.NewDecoder(r)
		cfg = &root{}
		err = d.Decode(cfg)
	)

	if err != nil {
		return nil, nil, err
	}

	return cfg.Convert()
}

func (r *root) Convert() (start workflow.Node, scope workflow.Variables, err error) {
	var (
		expr = workflow.Expressions(r.Scope.Cast()...)
	)

	if err = expr.Init(); err != nil {
		return nil, nil, err
	}

	if scope, err = expr.Run(context.Background()); err != nil {
		return nil, nil, err
	}

	if start, err = r.Steps.resolve(); err != nil {
		return nil, nil, err
	}

	return
}

func nodeErr(n *yaml.Node, format string, aa ...interface{}) error {
	format += " (%d:%d)"
	aa = append(aa, n.Line, n.Column)
	return fmt.Errorf(format, aa...)
}
