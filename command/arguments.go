package command

import arg "github.com/alexflint/go-arg"

type arguments struct {
	Filename string `arg:"positional"`
}

func getArguments(argv []string) (arguments, error) {
	as := arguments{}
	p, err := arg.NewParser(arg.Config{Program: argv[0]}, &as)

	if err != nil {
		return arguments{}, err
	} else if err := p.Parse(argv[1:]); err != nil {
		return arguments{}, err
	}

	return as, nil
}
