package service

//func (s *Service) sendVerification(user *model.User, email, fullName string) {
//	if email == "" {
//		return
//	}
//	go func() {
//		token, err := s.issueToken(context.Background(), user.UserID.String(), email)
//		if err != nil {
//			return
//		}
//		link := fmt.Sprintf("https://%s/v1/account/verify-email/%s", s.config.DomainName, token)
//		name := fullName
//		if name == "" {
//			name = supports.Val(user.Username)
//		}
//		content := fmt.Sprintf(
//			"Hi %s,\n\nThank you for registering at Gauas!\n\nPlease verify your email by clicking the link below.\n\nThe link expires in 24 hours.",
//			name,
//		)
//		_ = s.sendEmail(context.Background(), "confirmation", email, name, content, link)
//	}()
//}
