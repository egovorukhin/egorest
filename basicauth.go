package egorest

type BasicAuth struct {
	Name     string
	Password string
}

func SetBasicAuth(name, password string) *BasicAuth {
	return &BasicAuth{
		Name:     name,
		Password: password,
	}
}
