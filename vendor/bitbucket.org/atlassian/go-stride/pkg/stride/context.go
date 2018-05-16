package stride

// contextKey is a value used with with context.WithValue.
// context.WithValue(p, contextKey("something"), interface{})
type contextKey string

func (c contextKey) String() string {
	return "go-stride context key " + string(c)
}

var (
	contextKeyRequestContext = contextKey("request-context")
)
