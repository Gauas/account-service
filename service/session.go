package service

import "time"

type Session struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}
