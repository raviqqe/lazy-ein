package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/raviqqe/parcom"
)

type sign string

const (
	bindSign           sign = "="
	typeDefinitionSign      = ":"
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
		s.Exhaust(s.Prefix(s.blanks(), s.Many(s.bind()))),
	)
}

func (s *state) bind() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			return ast.NewBind(xs[4].(string), xs[2].(types.Type), xs[6].(ast.Expression)), nil
		},
		s.And(
			s.identifier(), s.sign(typeDefinitionSign), s.typ(), s.lineBreak(),
			s.identifier(), s.sign(bindSign), s.expression(), s.lineBreak(),
		),
	)
}

func (s *state) expression() parcom.Parser {
	return s.trim(s.numberLiteral())
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

func (s *state) typ() parcom.Parser {
	return s.trim(s.Or(s.numberType()))
}

func (s *state) numberType() parcom.Parser {
	return s.withDebugInformation(
		s.Str("Number"),
		func(_ interface{}, i *debug.Information) (interface{}, error) {
			return types.NewNumber(i), nil
		},
	)
}

func (s *state) withDebugInformation(
	p parcom.Parser,
	f func(interface{}, *debug.Information) (interface{}, error),
) parcom.Parser {
	return func() (interface{}, error) {
		i := s.debugInformation()
		x, err := p()

		if err != nil {
			return nil, err
		}

		return f(x, i)
	}
}

func (s *state) identifier() parcom.Parser {
	p := s.trim(s.Stringify(s.And(s.alphabet(), s.Many(s.Or(s.alphabet(), s.number())))))

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

func (s *state) sign(sg sign) parcom.Parser {
	return s.trim(s.Str(string(sg)))
}

func (s *state) trim(p parcom.Parser) parcom.Parser {
	return s.Wrap(s.spaces(), p, s.spaces())
}

func (s *state) lineBreak() parcom.Parser {
	return s.And(s.spaces(), s.newLine(), s.blanks())
}

func (s *state) blanks() parcom.Parser {
	return s.Void(s.Many(s.Or(s.space(), s.newLine())))
}

func (s *state) spaces() parcom.Parser {
	return s.Void(s.Many(s.space()))
}

func (s *state) space() parcom.Parser {
	return s.Void(s.Chars(" \t\r"))
}

func (s *state) newLine() parcom.Parser {
	return s.Void(s.Char('\n'))
}
