package gorex_test

import (
	"regexp"
	"testing"

	"github.com/Warashi/gorex"
)

func TestExpand(t *testing.T) {
	patterns := []string{
		`a`,
	}

	for _, p := range patterns {
		g, err := gorex.New(p)
		if err != nil {
			t.Errorf("error occured at pattern: %s", err.Error())
		}

		e := g.Expand()
		if len(e) == 0 {
			t.Errorf("%s does not expanded", p)
		}

		for _, s := range e {
			if m, _ := regexp.MatchString(p, s); !m {
				t.Errorf("%s does not match %s", s, p)
			}
		}
	}
}
func TestExpandEmpty(t *testing.T) {
	g, err := gorex.New("")
	if err != nil {
		t.Errorf("error occured at empty pattern: %s", err.Error())
	}
	s := g.Expand()
	if len(s) > 0 {
		t.Errorf("empty pattern expanded to %#v", s)
	}
}
