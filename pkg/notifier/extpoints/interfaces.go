package extpoints

type Driver interface {
	Notify(string) error
	SetOptions(opts map[string]string) error
	Uid() string
}
