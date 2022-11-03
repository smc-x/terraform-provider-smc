package cli

type Cli struct{}

func NewCli() (*Cli, error) {
	return &Cli{}, nil
}
