package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"gopkg.in/yaml.v3"
)

type (
	nodeMapper map[string][2]*yaml.Node
)

func mapNodes(nn ...*yaml.Node) nodeMapper {
	var (
		nm = nodeMapper{}
	)

	for i := 0; i < len(nn); i += 2 {
		nm.add(nn[i], nn[i+1])
	}

	return nm
}

func (nm nodeMapper) add(k, v *yaml.Node) {
	nm[k.Value] = [2]*yaml.Node{k, v}
}

func (nm nodeMapper) has(k string) bool {
	_, set := nm[k]
	return set
}

func (nm nodeMapper) vNode(k string) *yaml.Node {
	if nn, set := nm[k]; set {
		return nn[1]
	}
	return nil
}

func (nm nodeMapper) kNode(k string) *yaml.Node {
	if nn, set := nm[k]; set {
		return nn[0]
	}
	return nil
}

// returns first extra key found
func (nm nodeMapper) extra(kk ...string) *yaml.Node {
	ikk := slice.ToStringBoolMap(kk)
	for k := range nm {
		if !ikk[k] {
			return nm[k][0]
		}
	}
	return nil
}
