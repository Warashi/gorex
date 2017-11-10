package gorex

import (
	"regexp/syntax"

	"github.com/pkg/errors"
)

type gorex struct {
	prog *syntax.Prog
}

type Gorex interface {
	Expand() []string
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

func (g gorex) Expand() []string {
	return nil
}
