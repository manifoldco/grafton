package acceptance

import (
	"sort"
)

type node struct {
	f          *FeatureImpl
	children   []*node
	isTeardown bool
}

func (n *node) addChildren(features []*FeatureImpl) []*FeatureImpl {
	label := ""
	if n.f != nil {
		label = n.f.label
	}

	var leftover []*FeatureImpl
	for _, f := range features {
		if f.inside == label {
			n.children = append(n.children, &node{f: f})
			n.children = append(n.children, &node{f: f, isTeardown: true})
		} else {
			leftover = append(leftover, f)
		}
	}

	for _, c := range n.children {
		if !c.isTeardown {
			leftover = c.addChildren(leftover)
		}
	}

	sort.Stable(peerOrder(n.children))

	return leftover
}

func buildGraph(features []*FeatureImpl) *node {
	root := &node{}
	root.addChildren(features)

	return root
}

type peerOrder []*node

func (p peerOrder) Len() int      { return len(p) }
func (p peerOrder) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p peerOrder) Less(i, j int) bool {
	if p[i].f.label == p[j].f.label {
		return p[j].isTeardown
	}

	if p[i].f.before == p[j].f.label {
		return true
	}

	return false
}
