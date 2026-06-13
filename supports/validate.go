package supports

import "regexp"

var (
	Email = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	Phone = regexp.MustCompile(`^\d{10,15}$`)
)

func IsEmail(s string) bool { return Email.MatchString(s) }
func IsPhone(s string) bool { return Phone.MatchString(s) }
