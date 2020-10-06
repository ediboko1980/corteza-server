package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/workflow"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

type (
	resolver func(workflow.Node, steps) (workflow.Node, error)

	steps []*step

	step struct {
		wfn   workflow.Node
		ref   string
		next  string
		fn    resolver
		steps steps
	}
)

var (
	ErrDepsMissing = fmt.Errorf("waiting for deps to be resolved")
)

const (
	EndEventToken = "(FINAL)"
)

func (ss steps) resolve() (workflow.Node, error) {
	var (
		err  error
		next workflow.Node
		fss  = ss.flatten()
		deps = len(fss)
		max  = deps
	)

	for m := 0; m < max; m++ {
		println("=================================================================")
		println("iteration", m+1, "with", deps, "deps")
		//spew.Dump(fss)
		// going back from the last step and resolve it
		for i := len(fss) - 1; i >= 0; i-- {
			println("proc", i, "deps:", deps, "ref:", fss[i].ref)
			if fss[i].wfn != nil {
				println("  => workflow node already set")
				continue
			}

			if fss[i].fn == nil {
				println("  => resolve fn missing")
				continue
			}

			// find next node and send it to resolver
			next = nil
			if i < max-1 {
				next = fss[i+1].wfn
			}

			//spew.Dump(i, next)
			fss[i].wfn, err = fss[i].fn(next, fss)
			if err == nil {
				println("  => resolved")
				// resolved!!
				deps--
				continue
			}

			if err == ErrDepsMissing {
				println("  => dep missing")
				// offload to next iteration
				continue
			}
		}

		if deps == 0 {
			break
		}
	}

	if deps > 0 {
		return nil, fmt.Errorf("could not resolve all deps (%d)", deps)
	}

	for _, s := range fss {
		if s.wfn != nil {
			return s.wfn, nil
		}
	}

	return nil, nil
}

func (ss steps) flatten() steps {
	var out = ss

	// @todo check for unique refs.

	for _, s := range ss {
		out = append(out, s.steps.flatten()...)
	}

	for _, s := range ss {
		s.steps = nil
	}

	return out
}

func (ss steps) nodeLookupByRef(ref string) workflow.Node {
	for _, s := range ss {
		if s.ref == ref {
			return s.wfn
		}
	}

	return nil
}

func (ss *steps) UnmarshalYAML(n *yaml.Node) error {
	var (
		s *step
	)

	*ss = steps{}
	_ = spew.Dump

	switch n.Kind {
	case yaml.MappingNode:
		//spew.Dump(n.Content)
		for i := 0; i < len(n.Content); i += 2 {
			s = &step{ref: n.Content[i].Value}
			if s.ref == EndEventToken {
				s.fn = func(node workflow.Node, s steps) (workflow.Node, error) {
					// setting resolve fn instead of workflow node directly this
					// triggers counter dec and ensures resolve() works properly
					return workflow.Final(), nil
				}
				*ss = append(*ss, s)
				continue
			}

			if err := n.Content[i+1].Decode(s); err != nil {
				return err
			}

			if i > 0 {
				// link previous step to this one
				(*ss)[len(*ss)-1].next = s.ref
			}

			*ss = append(*ss, s)
		}

	case yaml.SequenceNode:
		for i, n := range n.Content {
			s = &step{}
			if err := n.Decode(s); err != nil {
				return err
			}

			if i > 0 {
				// link previous step to this one
				(*ss)[len(*ss)-1].next = s.ref
			}

			*ss = append(*ss, s)
		}

	default:
		return nodeErr(n, "expecting sequence or mapping node")

	}

	return nil
}

func (s *step) UnmarshalYAML(n *yaml.Node) error {
	//spew.Dump(n)
	if n.Kind == yaml.ScalarNode && n.Value == EndEventToken {
		s.fn = func(node workflow.Node, s steps) (workflow.Node, error) {
			return workflow.Final(), nil
		}

		return nil
	}

	if n.Kind != yaml.MappingNode {
		return fmt.Errorf("expecting mapping node")
	}

	var (
		err error
		nm  = mapNodes(n.Content...)

		next = func() {
			if nm.has("next") {
				s.next = nm.vNode("next").Value
				delete(nm, "next")
			}
		}
	)

	if nm.has("ref") {
		if s.ref != "" {
			return nodeErr(nm.vNode("ref"), "reference ID already set with mapping node")
		}

		s.ref = nm.vNode("ref").Value
		delete(nm, "ref")
	}

	if s.ref == "" {
		return nodeErr(n, "missing workflow step reference ID")
	}

	switch true {
	case nm.has("do"):
		if knode := nm.extra("do"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for step container", knode.Value)
		}

		if err = nm.vNode("do").Decode(&s.steps); err != nil {
			return err
		}

		s.next = s.steps[0].ref

	case nm.has("join-gateway"):
		if knode := nm.extra("join-gateway", "next"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for join-gateway", knode.Value)
		}

		if s.steps, s.fn, err = makeJoinGatewayResolver(nm.vNode("join-gateway")); err != nil {
			return err
		}

		next()

	case nm.has("fork-gateway"):
		if knode := nm.extra("fork-gateway"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for fork-gateway", knode.Value)
		}

		if s.steps, s.fn, err = makeForkGatewayResolver(nm.vNode("fork-gateway")); err != nil {
			return err
		}

		next()

	case nm.has("inclusive-gateway"):
		if knode := nm.extra("inclusive-gateway"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for inclusve-gateway", knode.Value)
		}

		if s.steps, s.fn, err = makeInclusiveGatewayResolver(nm.vNode("inclusive-gateway")); err != nil {
			return err
		}

		next()

	case nm.has("exclusive-gateway"):
		if knode := nm.extra("exclusive-gateway"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for exclusive-gateway", knode.Value)
		}

		if s.steps, s.fn, err = makeExclusiveGatewayResolver(nm.vNode("exclusive-gateway")); err != nil {
			return err
		}

		next()

	case nm.has("exec"):
		// @todo we need to access list of "registered" executables
		if knode := nm.extra("exec", "set", "params", "next"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for exec step", knode.Value)
		}

		// @todo translate params into executable's args
		// @todo translate exec's results to scope via set

		next()

	case nm.has("set"):
		if knode := nm.extra("set", "next"); knode != nil {
			return nodeErr(knode, "unexpected key %q found for set step", knode.Value)
		}

		if s.fn, err = makeSetResolver(nm.vNode("set")); err != nil {
			return err
		}

		next()

	default:
		return nodeErr(n, "unrecognized step configuration")

	}

	return nil
}

func makeSetResolver(set *yaml.Node) (resolver, error) {
	var (
		expr = &expressions{}
		err  = expr.UnmarshalYAML(set)
	)

	if err != nil {
		return nil, err
	}

	return func(n workflow.Node, ss steps) (workflow.Node, error) {
		if n == nil {
			return nil, ErrDepsMissing
		}

		return workflow.NewSetActivity(n, expr.Cast()...)
	}, nil
}
