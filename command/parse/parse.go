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
	bindSign             sign = "="
	typeDefinitionSign        = ":"
	functionSign              = "->"
	openParenthesisSign       = "("
	closeParenthesisSign      = ")"
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
		s.Exhaust(s.Block(s.bind())),
	)
}

func (s *state) bind() parcom.Parser {
	return s.withDebugInformation(
		s.WithPosition(
			s.And(
				s.SameColumn(
					s.WithPosition(s.And(s.identifier(), s.sign(typeDefinitionSign), s.typ())),
				),
				s.SameColumn(
					s.WithPosition(s.And(s.identifier(), s.arguments(), s.sign(bindSign), s.expression())),
				),
			),
		),
		func(x interface{}, i *debug.Information) (interface{}, error) {
			xs := x.([]interface{})
			ys := xs[0].([]interface{})
			zs := xs[1].([]interface{})

			if ys[0].(string) != zs[0].(string) {
				return nil, newError(
					"inconsistent identifiers between a type declaration and a definition",
					i,
				)
			}

			return ast.NewBind(
				zs[0].(string),
				zs[1].([]string),
				ys[2].(types.Type),
				zs[3].(ast.Expression),
			), nil
		},
	)
}

func (s *state) arguments() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			ss := make([]string, 0, len(xs))

			for _, x := range xs {
				ss = append(ss, x.(string))
			}

			return ss, nil
		},
		s.Many(s.identifier()),
	)
}

func (s *state) expression() parcom.Parser {
	return s.numberLiteral()
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
		s.token(
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
		),
	)
}

func (s *state) typ() parcom.Parser {
	return s.Lazy(func() parcom.Parser { return s.strictType() })
}

func (s *state) strictType() parcom.Parser {
	return s.Or(
		s.functionType(),
		s.scalarType(),
	)
}

func (s *state) argumentType() parcom.Parser {
	return s.Or(
		s.scalarType(),
		s.parenthesesed(s.functionType()),
	)
}

func (s *state) scalarType() parcom.Parser {
	return s.withDebugInformation(
		s.token(s.Str("Number")),
		func(_ interface{}, i *debug.Information) (interface{}, error) {
			return types.NewNumber(i), nil
		},
	)
}

func (s *state) functionType() parcom.Parser {
	return s.Lazy(func() parcom.Parser { return s.strictFunctionType() })
}

func (s *state) strictFunctionType() parcom.Parser {
	return s.withDebugInformation(
		s.And(s.argumentType(), s.sign(functionSign), s.typ()),
		func(x interface{}, i *debug.Information) (interface{}, error) {
			xs := x.([]interface{})
			return types.NewFunction(xs[0].(types.Type), xs[2].(types.Type), i), nil
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

func (s *state) parenthesesed(p parcom.Parser) parcom.Parser {
	return s.Wrap(s.sign(openParenthesisSign), p, s.sign(closeParenthesisSign))
}

func (s *state) identifier() parcom.Parser {
	return s.withDebugInformation(
		s.token(s.Stringify(s.And(s.alphabet(), s.Many(s.Or(s.alphabet(), s.number()))))),
		func(x interface{}, i *debug.Information) (interface{}, error) {
			if _, ok := keywords[keyword(x.(string))]; ok {
				return nil, newError(fmt.Sprintf("%#v is a keyword", x), i)
			}

			return x, nil
		},
	)
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
	return s.token(s.Str(string(sg)))
}

func (s *state) token(p parcom.Parser) parcom.Parser {
	return s.Wrap(s.blanks(), s.SameLineOrIndent(p), s.blanks())
}

func (s *state) blanks() parcom.Parser {
	return s.Void(s.Many(s.Chars(" \t\n\r")))
}
