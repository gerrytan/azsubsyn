package rp

import "fmt"

type ResourceProvider struct {
	Id                 string
	Namespace          string
	RegistrationState  string
	RegistrationPolicy string
}

func (rp *ResourceProvider) String() string {
	return fmt.Sprintf("ResourceProvider{Id: %s, Namespace: %s, RegistrationState: %s, RegistrationPolicy: %s}",
		rp.Id, rp.Namespace, rp.RegistrationState, rp.RegistrationPolicy)
}
