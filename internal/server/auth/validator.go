package auth

type Credentials interface {
}

type Validator interface {
	Validate(credentials Credentials) bool
}
