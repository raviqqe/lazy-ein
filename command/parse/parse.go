package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
	"github.com/raviqqe/parcom"
)

type sign string

const (
	bindSign               sign = "="
	typeDefinitionSign          = ":"
	functionSign                = "->"
	openParenthesisSign         = "("
	closeParenthesisSign        = ")"
	additionOperator            = sign(ast.Add)
	subtractionOperator         = sign(ast.Subtract)
	multiplicationOperator      = sign(ast.Multiply)
	divisionOperator            = sign(ast.Divide)
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
		s.HeteroBlock(
			s.WithPosition(s.And(s.identifier(), s.sign(typeDefinitionSign), s.typ())),
			s.WithPosition(s.And(s.identifier(), s.arguments(), s.sign(bindSign), s.expression())),
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

			e := zs[3].(ast.Expression)

			if as := zs[1].([]string); len(as) != 0 {
				e = ast.NewLambda(as, e)
			}

			return ast.NewBind(zs[0].(string), ys[2].(types.Type), e), nil
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
	return s.Lazy(func() parcom.Parser {
		return s.expressionWithOptions(true, true)
	})
}

func (s *state) expressionWithOptions(b, a bool) parcom.Parser {
	ps := []parcom.Parser{}

	if b {
		ps = append(ps, s.binaryOperatorTerm())
	}

	if a {
		ps = append(ps, s.application())
	}

	return s.Or(
		append(
			ps,
			s.numberLiteral(),
			s.let(),
			s.variable(),
			s.parenthesesed(s.expression()),
		)...,
	)
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

func (s *state) variable() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			return ast.NewVariable(x.(string)), nil
		},
		s.identifier(),
	)
}

func (s *state) application() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			ys := xs[1].([]interface{})
			es := make([]ast.Expression, 0, len(ys))

			for _, y := range ys {
				es = append(es, y.(ast.Expression))
			}

			return ast.NewApplication(xs[0].(ast.Expression), es), nil
		},
		s.And(s.expressionWithOptions(false, false), s.Many1(s.expressionWithOptions(false, false))),
	)
}

func (s *state) let() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			ys := xs[0].([]interface{})
			bs := make([]ast.Bind, 0, len(ys))

			for _, y := range ys {
				bs = append(bs, y.(ast.Bind))
			}

			return ast.NewLet(bs, xs[2].(ast.Expression)), nil
		},
		s.And(
			s.WithBlock1(
				s.keyword(letKeyword),
				s.untypedBind(),
			),
			s.keyword(inKeyword),
			s.expression(),
		),
	)
}

func (s *state) binaryOperatorTerm() parcom.Parser {
	return s.App(
		func(x interface{}) (interface{}, error) {
			xs := x.([]interface{})
			os := []ast.BinaryOperator{}
			es := []ast.Expression{xs[0].(ast.Expression)}

			for _, y := range xs[1].([]interface{}) {
				ys := y.([]interface{})

				os = append(os, ast.BinaryOperator(ys[0].(string)))
				es = append(es, ys[1].(ast.Expression))
			}

			for len(es) > 1 {
				es, os = reduceBinaryOperatorTerms(es, os)
			}

			return es[0], nil
		},
		s.And(
			s.expressionWithOptions(false, true),
			s.Many1(
				s.And(
					s.Or(
						s.sign(additionOperator),
						s.sign(subtractionOperator),
						s.sign(multiplicationOperator),
						s.sign(divisionOperator),
					),
					s.expressionWithOptions(false, true),
				),
			),
		),
	)
}

func reduceBinaryOperatorTerms(
	es []ast.Expression,
	os []ast.BinaryOperator,
) ([]ast.Expression, []ast.BinaryOperator) {
	if len(es) == 1 {
		return es, nil
	} else if len(es) == 2 || os[0].Priority() >= os[1].Priority() {
		return append(
			[]ast.Expression{ast.NewBinaryOperatorTerm(os[0], es[0], es[1])},
			es[2:]...,
		), os[1:]
	}

	ees, oos := reduceBinaryOperatorTerms(es[1:], os[1:])

	return append([]ast.Expression{es[0]}, ees...), append([]ast.BinaryOperator{os[0]}, oos...)
}

func (s *state) untypedBind() parcom.Parser {
	return s.withDebugInformation(
		s.WithPosition(s.And(s.identifier(), s.arguments(), s.sign(bindSign), s.expression())),
		func(x interface{}, i *debug.Information) (interface{}, error) {
			xs := x.([]interface{})
			e := xs[3].(ast.Expression)

			if as := xs[1].([]string); len(as) != 0 {
				e = ast.NewLambda(as, e)
			}

			return ast.NewBind(xs[0].(string), types.NewUnknown(i), e), nil
		},
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
				return nil, newError(fmt.Sprintf("'%v' is a keyword", x), i)
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

func (s *state) keyword(k keyword) parcom.Parser {
	return s.token(s.Str(string(k)))
}

func (s *state) sign(sg sign) parcom.Parser {
	return s.token(s.Str(string(sg)))
}

func (s *state) token(p parcom.Parser) parcom.Parser {
	return s.Suffix(s.SameLineOrIndent(p), s.Many(s.Chars(" \t\n\r")))
}
