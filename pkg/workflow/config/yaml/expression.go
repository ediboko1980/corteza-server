package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"gopkg.in/yaml.v3"
)

type (
	expression struct {
		Target string
		Source string
	}

	expressions []*expression
)

func (ee *expressions) UnmarshalYAML(n *yaml.Node) error {
	var (
		aux = expressions{}
	)

	for i := 0; i < len(n.Content); i += 2 {
		var (
			targetNode, exprNode = n.Content[i], n.Content[i+1]
			exp                  = &expression{Target: targetNode.Value}
		)

		// Basic check with node kind
		//
		// only scalar types are supported for now
		if exprNode.Kind != yaml.ScalarNode {
			return nodeErr(exprNode, "unexpected node kind, only scalar types supported for expression")
		}

		exp.SetSource(exprNode)
		aux = append(aux, exp)
	}

	*ee = aux
	return nil
}

func (ee expressions) Cast() []*workflow.Expr {
	var (
		oo = workflow.Expressions()
	)

	for _, e := range ee {
		oo.Push(workflow.NewExpr(e.Target, e.Source))
	}

	return oo
}

// SetSource sets value from yaml node and quotes it if needed
//
// It tries to understand what kind of input it got from the quote presence & style:
//  - single quotes indicate expression
//  - double quotes indicate string
//  - unquoted values are decoded and if possible, converted to expression
func (e *expression) SetSource(n *yaml.Node) {
	if n.Style == yaml.SingleQuotedStyle {
		// Already in the format we need it
		e.Source = n.Value
		return
	}

	if n.Style == yaml.DoubleQuotedStyle {
		// Re-quote it
		e.Source = fmt.Sprintf("%q", n.Value)
		return
	}

	// Try to parse the rest
	var (
		auxInt   int64
		auxFloat float64
		auxBool  bool
	)

	switch true {
	case nil == n.Decode(&auxFloat):
		e.Source = fmt.Sprintf("%f", auxFloat)
	case nil == n.Decode(&auxInt):
		e.Source = fmt.Sprintf("%d", auxInt)
	case nil == n.Decode(&auxBool):
		e.Source = fmt.Sprintf("%t", auxBool)
	default:
		e.Source = fmt.Sprintf("\"%v\"", n.Value)
	}
}
