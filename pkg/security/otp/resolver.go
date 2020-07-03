package otp

type Response struct {
	Otp bool `json:"otp"`
}

type Resolver interface {
	RequestOtp(url, organization string) (Response, error)
}