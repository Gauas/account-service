package service

//type TOTPSetup struct {
//	QRURL   string `json:"qr_code"`
//	Secret  string `json:"secret"`
//	Account string `json:"account"`
//	Issuer  string `json:"issuer"`
//}
//
//func (s *Service) GenerateTOTPSetup(ctx context.Context, userID uuid.UUID) (*TOTPSetup, error) {
//	user, err := s.GetUserByID(ctx, userID)
//	if err != nil {
//		return nil, err
//	}
//
//	currentMFA, err := s.GetMFAByType(ctx, userID, "totp")
//	if err != nil {
//		return nil, err
//	}
//	if currentMFA != nil && currentMFA.Enabled {
//		return nil, appError(http.StatusConflict, "TOTP is already enabled")
//	}
//
//	secret := make([]byte, 20)
//	if _, err := rand.Read(secret); err != nil {
//		return nil, appError(http.StatusInternalServerError, "failed to generate secret")
//	}
//	secretStr := base32.StdEncoding.EncodeToString(secret)
//
//	accountName := resolveAccountName(user, userID)
//
//	key, err := totp.Generate(totp.GenerateOpts{
//		Issuer:      "Gauas",
//		AccountName: accountName,
//		Secret:      secret,
//	})
//	if err != nil {
//		return nil, appError(http.StatusInternalServerError, "failed to generate TOTP key")
//	}
//
//	if currentMFA != nil {
//		currentMFA.Secret = &secretStr
//		if _, err := s.UpdateMFA(ctx, currentMFA); err != nil {
//			return nil, err
//		}
//	} else {
//		mfa := &model.UserMFA{ID: uuid.New(), UserID: userID, Type: "totp", Secret: &secretStr}
//		if _, err := s.CreateMFA(ctx, mfa); err != nil {
//			return nil, err
//		}
//	}
//
//	return &TOTPSetup{QRURL: key.URL(), Secret: secretStr, Account: accountName, Issuer: "Gauas"}, nil
//}
//
//func (s *Service) EnableTOTP(ctx context.Context, userID uuid.UUID, otpCode string) error {
//	mfa, err := s.GetMFAByType(ctx, userID, "totp")
//	if err != nil {
//		return err
//	}
//	if mfa == nil || mfa.Secret == nil {
//		return appError(http.StatusBadRequest, "no TOTP setup found, generate QR code first")
//	}
//	if mfa.Enabled {
//		return appError(http.StatusConflict, "TOTP already enabled")
//	}
//
//	if !totp.Validate(otpCode, *mfa.Secret) {
//		return appError(http.StatusBadRequest, "invalid OTP code")
//	}
//
//	now := time.Now()
//	mfa.Enabled = true
//	mfa.VerifiedAt = &now
//
//	_, err = s.UpdateMFA(ctx, mfa)
//	return err
//}
//
//func (s *Service) VerifyTOTP(ctx context.Context, userID uuid.UUID, otpCode, deviceID string) (string, string, time.Time, error) {
//	if deviceID == "" {
//		return "", "", time.Time{}, appError(http.StatusBadRequest, "device_id is required")
//	}
//
//	mfa, err := s.GetMFAByType(ctx, userID, "totp")
//	if err != nil {
//		return "", "", time.Time{}, err
//	}
//	if mfa == nil || !mfa.Enabled || mfa.Secret == nil {
//		return "", "", time.Time{}, appError(http.StatusBadRequest, "TOTP is not enabled")
//	}
//
//	if !totp.Validate(otpCode, *mfa.Secret) {
//		return "", "", time.Time{}, appError(http.StatusBadRequest, "invalid OTP code")
//	}
//
//	user, err := s.GetUserByID(ctx, userID)
//	if err != nil {
//		return "", "", time.Time{}, err
//	}
//
//	return s.CreateToken(ctx, user.UserID, user.Permission, deviceID)
//}
