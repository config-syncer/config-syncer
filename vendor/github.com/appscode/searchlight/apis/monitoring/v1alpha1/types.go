package v1alpha1

type Receiver struct {
	// For which state notification will be sent
	State string `json:"state,omitempty"`

	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}
