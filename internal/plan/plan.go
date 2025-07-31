package plan

type Plan struct {
	RpRegistrations []RpRegistration `json:"rpRegistrations"`
	PreviewFeatures []PreviewFeature `json:"previewFeatures"`
}

type RpRegistration struct {
	Namespace string `json:"namespace"` // eg: "Microsoft.Cache"
	Reason    string `json:"reason"`    // NotRegisteredInTarget | NotFoundInTarget
}

type PreviewFeature struct {
	Key       string `json:"key"`       // eg: "Dev"
	Namespace string `json:"namespace"` // eg: "Microsoft.DevAI"
	Reason    string `json:"reason"`    // NotRegisteredInTarget | NotFoundInTarget
}
