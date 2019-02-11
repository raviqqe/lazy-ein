package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarListCases(t *testing.T) {
	for _, es := range [][2]ast.Expression{
		// Single constant elements
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"x",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.x"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.x",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.x"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{
										ast.NewAlternative(
											ast.NewNumber(42),
											ast.NewCase(
												ast.NewVariable("$list-case.tail-0"),
												types.NewUnknown(nil),
												[]ast.Alternative{
													ast.NewAlternative(
														ast.NewList(types.NewUnknown(nil), nil),
														ast.NewNumber(42),
													),
												},
												d,
											),
										),
									},
									d,
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// Multiple constant elements
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{
								ast.NewListArgument(ast.NewNumber(42), false),
								ast.NewListArgument(ast.NewNumber(42), false),
							},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"x",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.x"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.x",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.x"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{
										ast.NewAlternative(
											ast.NewNumber(42),
											ast.NewCase(
												ast.NewVariable("$list-case.tail-0"),
												types.NewUnknown(nil),
												[]ast.Alternative{
													ast.NewAlternative(
														ast.NewList(
															types.NewUnknown(nil),
															[]ast.ListArgument{
																ast.NewListArgument(ast.NewVariable("$list-case.head-1"), false),
																ast.NewListArgument(ast.NewVariable("$list-case.tail-1"), true),
															},
														),
														ast.NewCase(
															ast.NewVariable("$list-case.head-1"),
															types.NewUnknown(nil),
															[]ast.Alternative{
																ast.NewAlternative(
																	ast.NewNumber(42),
																	ast.NewCase(
																		ast.NewVariable("$list-case.tail-1"),
																		types.NewUnknown(nil),
																		[]ast.Alternative{
																			ast.NewAlternative(
																				ast.NewList(
																					types.NewUnknown(nil),
																					nil,
																				),
																				ast.NewNumber(42),
																			),
																		},
																		d,
																	),
																),
															},
															d,
														),
													),
												},
												d,
											),
										),
									},
									d,
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// Hidden default alternatives
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("xs"), true)},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{},
				ast.NewDefaultAlternative(
					"xs",
					ast.NewNumber(42),
				),
			),
		},
		// Hidden default alternatives with other alternatives
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{
								ast.NewListArgument(ast.NewNumber(42), false),
							},
						),
						ast.NewNumber(42),
					),
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("xs"), true)},
						),
						ast.NewNumber(13),
					),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"xs",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.xs"),
							),
						},
						ast.NewNumber(13),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.xs",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.xs"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{
										ast.NewAlternative(
											ast.NewNumber(42),
											ast.NewCase(
												ast.NewVariable("$list-case.tail-0"),
												types.NewUnknown(nil),
												[]ast.Alternative{
													ast.NewAlternative(
														ast.NewList(types.NewUnknown(nil), nil),
														ast.NewNumber(42),
													),
												},
												d,
											),
										),
									},
									d,
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// Hidden default alternatives with the same usual default alternatives
		{
			ast.NewCaseWithoutDefault(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{
								ast.NewListArgument(ast.NewNumber(42), false),
							},
						),
						ast.NewNumber(42),
					),
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("xs"), true)},
						),
						ast.NewNumber(13),
					),
				},
			),
			newListCaseDesugarer().Desugar(
				ast.NewCase(
					ast.NewVariable("argument"),
					types.NewUnknown(nil),
					[]ast.Alternative{
						ast.NewAlternative(
							ast.NewList(
								types.NewUnknown(nil),
								[]ast.ListArgument{
									ast.NewListArgument(ast.NewNumber(42), false),
								},
							),
							ast.NewNumber(42),
						),
					},
					ast.NewDefaultAlternative("xs", ast.NewNumber(13)),
				),
			),
		},
		// Variable elements
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("x"), false)},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("y", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.y"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.y",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.y"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{},
									ast.NewDefaultAlternative(
										"",
										ast.NewCase(
											ast.NewVariable("$list-case.tail-0"),
											types.NewUnknown(nil),
											[]ast.Alternative{
												ast.NewAlternative(
													ast.NewList(types.NewUnknown(nil), nil),
													ast.NewLet(
														[]ast.Bind{
															ast.NewBind(
																"x",
																types.NewUnknown(nil),
																ast.NewVariable("$list-case.head-0"),
															),
														},
														ast.NewNumber(42),
													),
												),
											},
											d,
										),
									),
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// Duplicate constant elements
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
						),
						ast.NewNumber(42),
					),
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{
								ast.NewListArgument(ast.NewNumber(42), false),
								ast.NewListArgument(ast.NewVariable("xs"), true),
							},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"x",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.x"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.x",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.x"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{
										ast.NewAlternative(
											ast.NewNumber(42),
											ast.NewCase(
												ast.NewVariable("$list-case.tail-0"),
												types.NewUnknown(nil),
												[]ast.Alternative{
													ast.NewAlternative(
														ast.NewList(types.NewUnknown(nil), nil),
														ast.NewNumber(42),
													),
												},
												ast.NewDefaultAlternative("xs", ast.NewNumber(42)),
											),
										),
									},
									d,
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// Empty and non-empty list patterns
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
						),
						ast.NewNumber(42),
					),
					ast.NewAlternative(ast.NewList(types.NewUnknown(nil), nil), ast.NewNumber(42)),
				},
				ast.NewDefaultAlternative("x", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"x",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.x"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.x",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.x"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{
										ast.NewAlternative(
											ast.NewNumber(42),
											ast.NewCase(
												ast.NewVariable("$list-case.tail-0"),
												types.NewUnknown(nil),
												[]ast.Alternative{
													ast.NewAlternative(
														ast.NewList(types.NewUnknown(nil), nil),
														ast.NewNumber(42),
													),
												},
												d,
											),
										),
									},
									d,
								),
							),
							ast.NewAlternative(ast.NewList(types.NewUnknown(nil), nil), ast.NewNumber(42)),
						},
						d,
					),
				)
			}(),
		},
		// Constant elements after variable elements
		{
			ast.NewCase(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("x"), false)},
						),
						ast.NewNumber(42),
					),
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
						),
						ast.NewNumber(42),
					),
				},
				ast.NewDefaultAlternative("y", ast.NewNumber(42)),
			),
			func() ast.Expression {
				d := ast.NewDefaultAlternative(
					"",
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								types.NewUnknown(nil),
								ast.NewVariable("$default-alternative.y"),
							),
						},
						ast.NewNumber(42),
					),
				)

				return ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"$default-alternative.y",
							types.NewUnknown(nil),
							ast.NewVariable("argument"),
						),
					},
					ast.NewCase(
						ast.NewVariable("$default-alternative.y"),
						types.NewUnknown(nil),
						[]ast.Alternative{
							ast.NewAlternative(
								ast.NewList(
									types.NewUnknown(nil),
									[]ast.ListArgument{
										ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
										ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
									},
								),
								ast.NewCase(
									ast.NewVariable("$list-case.head-0"),
									types.NewUnknown(nil),
									[]ast.Alternative{},
									ast.NewDefaultAlternative(
										"",
										ast.NewCase(
											ast.NewVariable("$list-case.tail-0"),
											types.NewUnknown(nil),
											[]ast.Alternative{
												ast.NewAlternative(
													ast.NewList(types.NewUnknown(nil), nil),
													ast.NewLet(
														[]ast.Bind{
															ast.NewBind(
																"x",
																types.NewUnknown(nil),
																ast.NewVariable("$list-case.head-0"),
															),
														},
														ast.NewNumber(42),
													),
												),
											},
											ast.NewDefaultAlternative(
												"",
												ast.NewCase(
													ast.NewVariable("$list-case.head-0"),
													types.NewUnknown(nil),
													[]ast.Alternative{
														ast.NewAlternative(
															ast.NewNumber(42),
															ast.NewCase(
																ast.NewVariable("$list-case.tail-0"),
																types.NewUnknown(nil),
																[]ast.Alternative{
																	ast.NewAlternative(
																		ast.NewList(types.NewUnknown(nil), nil),
																		ast.NewNumber(42),
																	),
																},
																d,
															),
														),
													},
													d,
												),
											),
										),
									),
								),
							),
						},
						d,
					),
				)
			}(),
		},
		// No default alternatives
		{
			ast.NewCaseWithoutDefault(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
						),
						ast.NewNumber(42),
					),
				},
			),
			ast.NewCaseWithoutDefault(
				ast.NewVariable("argument"),
				types.NewUnknown(nil),
				[]ast.Alternative{
					ast.NewAlternative(
						ast.NewList(
							types.NewUnknown(nil),
							[]ast.ListArgument{
								ast.NewListArgument(ast.NewVariable("$list-case.head-0"), false),
								ast.NewListArgument(ast.NewVariable("$list-case.tail-0"), true),
							},
						),
						ast.NewCaseWithoutDefault(
							ast.NewVariable("$list-case.head-0"),
							types.NewUnknown(nil),
							[]ast.Alternative{
								ast.NewAlternative(
									ast.NewNumber(42),
									ast.NewCaseWithoutDefault(
										ast.NewVariable("$list-case.tail-0"),
										types.NewUnknown(nil),
										[]ast.Alternative{
											ast.NewAlternative(
												ast.NewList(types.NewUnknown(nil), nil),
												ast.NewNumber(42),
											),
										},
									),
								),
							},
						),
					),
				},
			),
		},
	} {
		assert.Equal(
			t,
			newModuleWithExpression(es[1]),
			desugarListCases(newModuleWithExpression(es[0])),
		)
	}
}

func TestDesugarListCasesWithDesugaredCases(t *testing.T) {
	for _, c := range []ast.Case{
		// Non-list case expressions
		ast.NewCase(
			ast.NewNumber(42),
			types.NewUnknown(nil),
			[]ast.Alternative{ast.NewAlternative(ast.NewNumber(42), ast.NewNumber(42))},
			ast.NewDefaultAlternative("x", ast.NewNumber(42)),
		),
		// No alternatives
		ast.NewCase(
			ast.NewList(
				types.NewUnknown(nil),
				[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
			),
			types.NewUnknown(nil),
			[]ast.Alternative{},
			ast.NewDefaultAlternative("x", ast.NewNumber(42)),
		),
		// Empty list patterns with default alternatives
		ast.NewCase(
			ast.NewList(
				types.NewUnknown(nil),
				[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
			),
			types.NewUnknown(nil),
			[]ast.Alternative{
				ast.NewAlternative(
					ast.NewList(types.NewUnknown(nil), nil),
					ast.NewNumber(42),
				),
			},
			ast.NewDefaultAlternative("x", ast.NewNumber(42)),
		),
		// Empty list patterns without default alternatives
		ast.NewCaseWithoutDefault(
			ast.NewList(
				types.NewUnknown(nil),
				[]ast.ListArgument{ast.NewListArgument(ast.NewNumber(42), false)},
			),
			types.NewUnknown(nil),
			[]ast.Alternative{
				ast.NewAlternative(
					ast.NewList(types.NewUnknown(nil), nil),
					ast.NewNumber(42),
				),
			},
		),
	} {
		assert.Equal(t, newModuleWithExpression(c), desugarListCases(newModuleWithExpression(c)))
	}
}

func newModuleWithExpression(e ast.Expression) ast.Module {
	return ast.NewModule("", []ast.Bind{ast.NewBind("x", types.NewUnknown(nil), e)})
}
