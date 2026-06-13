package mfa

import (
	"bytes"
	"fmt"
	"image/png"
	"net/url"
	"strconv"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	ISSUER      = "Gauas"
	SECRET_SZ   = 20
	TOTP_PERIOD = 30
	TOTP_SKEW   = 1
)

var TOTP_DIGITS = otp.DigitsSix

func BuildKey(account string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      ISSUER,
		AccountName: account,
		SecretSize:  SECRET_SZ,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func KeyURL(account, secret string) string {
	keyURL := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   "/" + url.PathEscape(fmt.Sprintf("%s:%s", ISSUER, account)),
	}

	query := keyURL.Query()
	query.Set("algorithm", "SHA1")
	query.Set("digits", strconv.Itoa(int(TOTP_DIGITS.Length())))
	query.Set("issuer", ISSUER)
	query.Set("period", strconv.Itoa(TOTP_PERIOD))
	query.Set("secret", secret)
	keyURL.RawQuery = query.Encode()

	return keyURL.String()
}

func BuildEnrollment(account string) (string, string, []byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      ISSUER,
		AccountName: account,
		SecretSize:  SECRET_SZ,
	})
	if err != nil {
		return "", "", nil, err
	}

	img, err := key.Image(256, 256)
	if err != nil {
		return "", "", nil, err
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return "", "", nil, err
	}

	return key.Secret(), key.URL(), out.Bytes(), nil
}

func Verify(code, secret string, now time.Time) bool {
	ok, err := totp.ValidateCustom(code, secret, now, totp.ValidateOpts{
		Period:    TOTP_PERIOD,
		Skew:      TOTP_SKEW,
		Digits:    TOTP_DIGITS,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}

	return ok
}

func AccountName(userKey string) string {
	return fmt.Sprintf("user-%s", userKey)
}
