package console

type Console struct{}

func New() (*Console, error) {
	return &Console{}, nil
}
