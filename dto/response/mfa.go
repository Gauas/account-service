package response

type TOTPSetupResponse struct {
	Email   string `json:"email"`
	QRURL   string `json:"qr_code"`
	Secret  string `json:"secret"`
	Account string `json:"account"`
	Issuer  string `json:"issuer"`
}
