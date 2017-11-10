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

type runeRanges []rune

func (rr runeRanges) strs() []string {
	var l int32
	for i := 0; i < len(rr); i = i + 2 {
		l += rr[i+1] - rr[i] + 1
	}
	result := make([]string, 0, l)

	for i := 0; i < len(rr); i = i + 2 {
		for r := rr[i]; r <= rr[i+1]; r++ {
			result = append(result, string(r))
		}
	}

	return result
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
		case syntax.InstNop:
			return nil, nil

		case syntax.InstMatch, syntax.InstFail:
			result = append(result, n.s)

		case syntax.InstAlt:
			s.push(node{n.s, inst.Arg})
			s.push(node{n.s, inst.Out})

		case syntax.InstCapture, syntax.InstEmptyWidth:
			s.push(node{n.s, inst.Out})

		case syntax.InstRuneAny:
			rr := runeRanges{0, 1114111}
			for _, r := range rr.strs() {
				s.push(node{s: n.s + r, pc: inst.Out})
			}

		case syntax.InstRuneAnyNotNL:
			rr := runeRanges{0, 9, 11, 1114111}
			for _, r := range rr.strs() {
				s.push(node{s: n.s + r, pc: inst.Out})
			}

		case syntax.InstRune:
			for _, r := range runeRanges(inst.Rune).strs() {
				s.push(node{s: n.s + r, pc: inst.Out})
			}

		case syntax.InstRune1:
			s.push(node{s: n.s + string(inst.Rune[0]), pc: inst.Out})

		default:
			return nil, errors.New("not implemented")
		}

	}
	return result, nil
}
