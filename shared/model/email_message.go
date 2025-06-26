package sharedmodel

type EmailMessage struct {
	From        string
	To          []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string
}
