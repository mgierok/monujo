package console

type Console struct{}

func NewConsole() (*Console, error) {
	return &Console{}, nil
}
