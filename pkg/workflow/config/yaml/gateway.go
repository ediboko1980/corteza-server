package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

type (
	paths []path
	path  struct {
		condition string
		ref       string
		def       steps
		fn        resolver
	}
)

func makeJoinGatewayResolver(inPaths *yaml.Node) (steps, resolver, error) {
	return nil, nil, nil
}

func makeForkGatewayResolver(outPaths *yaml.Node) (steps, resolver, error) {
	return nil, nil, nil
}

func makeInclusiveGatewayResolver(outPaths *yaml.Node) (steps, resolver, error) {
	return nil, nil, nil
}

func makeExclusiveGatewayResolver(outPaths *yaml.Node) (steps, resolver, error) {
	var (
		pp  = paths{}
		err = outPaths.Decode(&pp)
	)

	if err != nil {
		return nil, nil, err
	}

	fn := func(n workflow.Node, ss steps) (workflow.Node, error) {
		for _, p := range pp {
			if ss.nodeLookupByRef(p.ref) == nil {
				spew.Dump(p)
				println("dep missing :D", p.ref)
				return nil, ErrDepsMissing
			}
		}

		egw, _ := workflow.NewExclGateway()

		for _, p := range pp {
			pn := ss.nodeLookupByRef(p.ref)
			if p.condition == "" {
				err = egw.AddPaths(workflow.NewGatewayNoCondition(pn))
			} else {
				err = egw.AddPaths(workflow.NewGatewayCondition(p.condition, pn))
			}

			if err != nil {
				return nil, err
			}
		}

		return egw, nil
	}

	return pp.steps(), fn, nil
}

func (p *path) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return fmt.Errorf("expecting mapping node")
	}

	var (
		err error
		nm  = mapNodes(n.Content...)
	)

	if knode := nm.extra("if", "next"); knode != nil {
		return nodeErr(knode, "unexpected key %q found for gateway path", knode.Value)
	}

	if nm.has("if") {
		p.condition = nm.vNode("if").Value
	}

	if next := nm.vNode("next"); next != nil {
		switch next.Kind {
		case yaml.ScalarNode:
			// string reference to next node
			p.ref = next.Value
		default:
			if err = next.Decode(&p.def); err != nil {
				return err
			}
		}
	}

	return nil
}

// steps collects steps from all paths and returns them
//
// this is needed because sub-steps can be defined as part of the gateway path def
func (pp paths) steps() steps {
	var ss = steps{}

	for _, p := range pp {
		ss = append(ss, p.def...)
	}

	return ss
}
