package plan

type Plan struct {
	RpRegistrations []RpRegistration `json:"rpRegistrations"`
}

type RpRegistration struct {
	Namespace string `json:"namespace"`
	// NotRegisteredInTarget | NotFoundInTarget
	Reason string `json:"reason"`
}
