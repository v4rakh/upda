package auth

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type StaticUserPasswordValidator struct {
	permitted map[string]string
}

func NewStaticUserPasswordValidator(permitted map[string]string) Validator {
	return &StaticUserPasswordValidator{permitted: permitted}
}

func (v *StaticUserPasswordValidator) Validate(credentials Credentials) bool {
	userPwCredentials, ok := credentials.(*UserCredentials)

	if !ok {
		return false
	}

	value, found := v.permitted[userPwCredentials.Username]
	if !found {
		return false
	}

	return userPwCredentials.Password == value
}
