package parse

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/parcom"
)

type sign string

const (
	bindSign sign = "="
)

type keyword string

const (
	caseKeyword   keyword = "case"
	exportKeyword         = "export"
	importKeyword         = "import"
	inKeyword             = "in"
	letKeyword            = "let"
	ofKeyword             = "of"
)

var keywords = map[keyword]struct{}{
	caseKeyword:   {},
	exportKeyword: {},
	importKeyword: {},
	inKeyword:     {},
	letKeyword:    {},
	ofKeyword:     {},
}

const blankCharacters = " \t\n\r"

// Parse parses a module into an AST.
func Parse(f, s string) (ast.Module, error) {
	x, err := newState(f, s).module(f)()

	if err != nil {
		return ast.Module{}, err
	}

	return x.(ast.Module), nil
}

func (s *state) module(f string) parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			bs := make([]ast.Bind, 0, len(xs))

			for _, x := range xs {
				bs = append(bs, x.(ast.Bind))
			}

			return ast.NewModule(f, bs), nil
		},
		s.Exhaust(
			s.Many(s.bind()),
			func(parcom.State) error {
				return errors.New("syntax error")
			},
		),
	)
}

func (s *state) bind() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			return ast.NewBind(xs[0].(string), xs[2].(ast.Expression)), nil
		},
		s.And(s.identifier(), s.sign(bindSign), s.expression()),
	)
}

func (s *state) expression() parcom.Parser {
	return s.strip(s.numberLiteral())
}

func (s *state) numberLiteral() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			n, err := strconv.ParseFloat(x.(string), 64)

			if err != nil {
				return nil, err
			}

			return ast.NewNumber(n), nil
		},
		s.Stringify(
			s.And(
				s.Maybe(s.Str("-")),
				s.Or(
					s.And(s.Chars("123456789"), s.Many(s.number())),
					s.Str("0"),
				),
				s.Maybe(s.And(s.Str("."), s.Many1(s.number()))),
			),
		),
	)
}

func (s *state) identifier() parcom.Parser {
	p := s.strip(s.Stringify(s.And(s.alphabet(), s.Many(s.Or(s.alphabet(), s.number())))))

	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		} else if _, ok := keywords[keyword(x.(string))]; ok {
			return nil, fmt.Errorf("%#v is a keyword", x)
		}

		return x, nil
	}
}

func (s *state) alphabet() parcom.Parser {
	as := ""

	for r := 'a'; r <= 'z'; r++ {
		as += string(r)
	}

	return s.Chars(as + strings.ToUpper(as))
}

func (s *state) number() parcom.Parser {
	ns := ""

	for i := 0; i <= 9; i++ {
		ns += strconv.Itoa(i)
	}

	return s.Chars(ns)
}

func (s *state) keyword(k keyword) parcom.Parser {
	return s.strip(s.Str(string(k)))
}

// TODO: Remove nolint attribute when signs other than = are implemented.
func (s *state) sign(sg sign) parcom.Parser { // nolint: unparam
	return s.strip(s.Str(string(sg)))
}

func (s *state) strip(p parcom.Parser) parcom.Parser {
	return s.Wrap(s.blank(), p, s.blank())
}

func (s *state) blank() parcom.Parser {
	return s.Void(s.Many(s.Chars(blankCharacters)))
}
