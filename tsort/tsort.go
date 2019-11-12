// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package tsort

import (
	"bytes"
	"fmt"
	"sort"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type vertex struct {
	id      string
	depth   int
	lowlink int
	onstack bool
	adjlist []*vertex
}

func (v *vertex) reset() {
	v.depth = 0
	v.lowlink = 0
	v.onstack = false
}

func (v *vertex) String() string {
	var s string
	s = v.id + ": [ "
	for _, a := range v.adjlist {
		s += a.id + " "
	}
	s += "]"
	return s
}

func (v *vertex) Dot(unconstrained bool) string {
	var buf bytes.Buffer
	if len(v.adjlist) == 0 {
		buf.WriteByte('\t')
		buf.WriteByte('"')
		buf.WriteString(v.id)
		buf.WriteByte('"')
		buf.WriteByte(';')
		buf.WriteByte('\n')
	}
	for _, adj := range v.adjlist {
		buf.WriteByte('\t')
		buf.WriteByte('"')
		buf.WriteString(v.id)
		buf.WriteByte('"')
		buf.WriteString(" -> ")
		buf.WriteByte('"')
		buf.WriteString(adj.id)
		buf.WriteByte('"')
		if unconstrained {
			buf.WriteString(" [constraint=false]")
		}
		buf.WriteByte(';')
		buf.WriteByte('\n')
	}
	return buf.String()
}

func makevertex(id string) *vertex {
	return &vertex{id: id, adjlist: make([]*vertex, 0)}
}

type scc struct {
	vs []*vertex
}

func (s *scc) append(v *vertex) { s.vs = append(s.vs, v) }

func (s *scc) String() string {
	var str string
	str = "[ "
	for _, v := range s.vs {
		str += v.id + " "
	}
	str += "]"
	return str
}

type stack struct {
	s []*vertex
}

func (s *stack) Push(v *vertex) {
	s.s = append(s.s, v)
	v.onstack = true
}

func (s *stack) Pop() (v *vertex) {
	v = s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	v.onstack = false
	return v
}

type Graph struct {
	vs map[string]*vertex
}

func (g *Graph) keys() []string {
	keys := make([]string, 0, len(g.vs))
	for k, _ := range g.vs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func New() *Graph {
	return &Graph{
		vs: make(map[string]*vertex),
	}
}

func (g *Graph) String() string {
	var s string
	for _, v := range g.vs {
		s += v.String() + "\n"
	}
	return s
}

func (g *Graph) reset() {
	for _, v := range g.vs {
		v.reset()
	}
}

func (g *Graph) AddEdge(v, w string) {
	v1 := g.AddVertex(v)
	v2 := g.AddVertex(w)
	v1.adjlist = append(v1.adjlist, v2)
}

func (g *Graph) AddVertex(v string) *vertex {
	v1, ok := g.vs[v]
	if !ok {
		v1 = makevertex(v)
		g.vs[v] = v1
	}
	return v1
}

func (g *Graph) HasEdge(v, w string) bool {
	v1, ok := g.vs[v]
	if !ok {
		return false
	}
	for _, adj := range v1.adjlist {
		if adj.id == w {
			return true
		}
	}
	return false
}

/*
 * Tarjan's algorithm
 * Tarjan, R. E. (1972), "Depth-first search and linear graph algorithms",
 * SIAM Journal on Computing 1 (2): 146–160, doi:10.1137/0201010
 */
func (g *Graph) tarjan() []*scc {
	//TODO(jhs): Clean up the implementation, directly out of paper here

	var strongconnect func(v *vertex)

	sccs := make([]*scc, 0, len(g.vs))
	depth := 0
	s := stack{s: make([]*vertex, 0)}

	strongconnect = func(v *vertex) {
		depth++
		v.depth = depth
		v.lowlink = depth
		s.Push(v)
		for _, w := range v.adjlist {
			if w.depth == 0 {
				/*(v,w) is a tree arc; recursively visit w*/
				strongconnect(w)
				v.lowlink = min(v.lowlink, w.lowlink)
			} else if w.onstack {
				/*
				 * (v,w) is a frond or cross-link
				 * i.e. it is in the current SCC
				 * set the link id to the lower depth
				 */
				v.lowlink = min(v.lowlink, w.depth)
			}
		}
		if v.lowlink == v.depth {
			/*
			 * v is the root of a component
			 * build a new component and append it
			 * to the list of componenets
			 */
			sc := &scc{vs: make([]*vertex, 0, 1)}
			for {
				w := s.Pop()
				sc.append(w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, sc)
		}
	}

	/*sort the hash keys to make the tsort stable*/
	keys := g.keys()

	/*loop over the vertices of the graph and find the strongly connected components*/
	for _, k := range keys {
		v := g.vs[k]
		if v.depth != 0 {
			continue
		}
		strongconnect(v)
	}
	return sccs
}

func (g *Graph) sort() ([]string, *scc) {
	defer g.reset()
	sccs := g.tarjan()
	out := make([]string, 0, len(sccs))
	for _, sc := range sccs {
		if len(sc.vs) > 1 {
			return nil, sc
		}
		out = append(out, sc.vs[0].id)
	}
	return out, nil

}

func (g *Graph) Sort() ([]string, error) {
	out, scc := g.sort()
	if scc != nil {
		return nil, fmt.Errorf("cycle detected %s", scc)
	}
	return out, nil
}

func (g *Graph) SortDot() string {
	var buf bytes.Buffer
	buf.WriteString("digraph \"graph\" {\n")
	sort, scc := g.sort()
	if scc != nil {
		buf.WriteString("\tranksep=\"5.0 equally\";")
		for _, v := range g.keys() {
			buf.WriteString(g.vs[v].Dot(false))
		}
		buf.WriteByte('\t')
		hd, tl := scc.vs[0], scc.vs[1:]
		buf.WriteByte('"')
		buf.WriteString(hd.id)
		buf.WriteByte('"')
		for _, s := range tl {
			buf.WriteString(" -> ")
			buf.WriteByte('"')
			buf.WriteString(s.id)
			buf.WriteByte('"')
		}
		buf.WriteString(" -> ")
		buf.WriteByte('"')
		buf.WriteString(hd.id)
		buf.WriteByte('"')
		buf.WriteString(" [color=red,style=bold,weight=5];\n")
		buf.WriteString("}\n")
		return buf.String()
	}
	for _, v := range g.keys() {
		buf.WriteString(g.vs[v].Dot(true))
	}
	buf.WriteByte('\t')
	hd, tl := sort[0], sort[1:]
	buf.WriteByte('"')
	buf.WriteString(hd)
	buf.WriteByte('"')
	for _, s := range tl {
		buf.WriteString(" -> ")
		buf.WriteByte('"')
		buf.WriteString(s)
		buf.WriteByte('"')
	}
	buf.WriteString(" [color=blue,style=bold,weight=5];\n")
	buf.WriteString("}\n")
	return buf.String()
}

func (g *Graph) Dot() string {
	var buf bytes.Buffer
	buf.WriteString("digraph \"yang module graph\" {\n")
	buf.WriteString("\tranksep=\"5.0 equally\";")
	for _, v := range g.keys() {
		buf.WriteString(g.vs[v].Dot(false))
	}
	buf.WriteString("}\n")
	return buf.String()
}
