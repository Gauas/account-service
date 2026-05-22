package supports

import "regexp"

var (
	Email = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	Phone = regexp.MustCompile(`^\d{10,15}$`)
)

func Match(s string, re *regexp.Regexp) bool {
	return re.MatchString(s)
}

func IsEmail(s string) bool { return Match(s, Email) }
func IsPhone(s string) bool { return Match(s, Phone) }
