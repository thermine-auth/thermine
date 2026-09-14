package setup

// setupRequest is the form the first administrator is made from.
//
// There is no username: the address is the account. Asking for one more thing
// to remember, at the one moment someone is trying to get in for the first
// time, buys nothing.
type setupRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=10,max=128"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}
