package notify

type ByEmail interface {
	From(from string) ByEmail
	WithSubject(subject string) ByEmail
	WithBody(body string) ByEmail
	WithTag(tag string) ByEmail
	To(to string, cc ...string) ByEmail
	Send() error
	SendHtml() error
}

type BySMS interface {
	From(from string) BySMS
	WithBody(body string) BySMS
	To(to string, cc ...string) BySMS
	Send() error
}

type ByChat interface {
	WithBody(body string) ByChat
	To(to string, cc ...string) ByChat
	Send() error
}
