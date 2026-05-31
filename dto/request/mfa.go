package request

type GenerateTOTPRequest struct{}

type EnableTOTPRequest struct {
	OTPCode string `json:"otp_code"`
}

type VerifyTOTPRequest struct {
	OTPCode string `json:"otp_code"`
}
