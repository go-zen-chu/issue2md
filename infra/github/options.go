package github

type Option struct {
	baseURL string
	token   string
}

type OptionSetter func(*Option)

func BaseURL(baseURL string) OptionSetter {
	return func(o *Option) {
		o.baseURL = baseURL
	}
}

func Token(token string) OptionSetter {
	return func(o *Option) {
		o.token = token
	}
}
