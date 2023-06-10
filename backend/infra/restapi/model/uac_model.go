package model

type QrCodeVerifyInput struct {
	QrCodeKey string `json:"qrCodeKey"`

	QrCodeValue string `json:"qrCodeValue"`

	LoginClientIp string `json:"loginClientIp"`

	OriginSystemCode string `json:"originSystemCode"`

	LoginSystemCode string `json:"loginSystemCode"`

	VerifyCode string `json:"verifyCode"`
}

type AuthCode struct {
	Code  string `json:"code"`
	EnMsg string `json:"enMsg"`
	Msg   string `json:"msg"`
}

type OtherMsg struct {
	Message string `json:"message"`
}

type QrCodeVerifyOutput struct {
	Bo    AuthCode          `json:"bo"`
	Code  TokenResponseCode `json:"code"`
	Other OtherMsg          `json:"other"`
}
