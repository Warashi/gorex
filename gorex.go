package gorex

import (
	"regexp/syntax"

	"github.com/pkg/errors"
)

type node struct {
	s  string
	pc uint32
}

type nodeStack []node

func (s *nodeStack) push(n node) {
	*s = append(*s, n)
}

func (s *nodeStack) pop() node {
	n := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return n
}

type gorex struct {
	prog *syntax.Prog
}

type Gorex interface {
	Expand() ([]string, error)
}

func New(pattern string) (Gorex, error) {
	expr, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return nil, errors.Wrap(err, "parse failed")
	}

	prog, err := syntax.Compile(expr.Simplify())
	if err != nil {
		return nil, errors.Wrap(err, "compile failed")
	}

	return gorex{prog: prog}, nil
}

func (g gorex) Expand() ([]string, error) {
	went := make(map[uint32]struct{})

	var result []string
	var s nodeStack
	s.push(node{"", uint32(g.prog.Start)})

	for len(s) > 0 {
		n := s.pop()

		if _, ok := went[n.pc]; ok {
			return nil, errors.New("infinite loop exists")
		}
		went[n.pc] = struct{}{}

		inst := g.prog.Inst[n.pc]
		switch inst.Op {
		case syntax.InstMatch, syntax.InstFail:
			result = append(result, n.s)

		case syntax.InstAlt:
			s.push(node{n.s, inst.Arg})
			s.push(node{n.s, inst.Out})

		case syntax.InstCapture, syntax.InstEmptyWidth:
			s.push(node{n.s, inst.Out})
		}
	}
	return result, nil
}
