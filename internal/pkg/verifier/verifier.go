package verifier

import (
	"regexp"
)

var (
	nicknameRegexp = regexp.MustCompile(`^[\w.]+$`)
	emailRegexp    = regexp.MustCompile(`^[\w\-.]+@[\w\-.]+\.[a-zA-Z]{2,6}$`)
	slugRegexp     = regexp.MustCompile(`^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$`)
)

func IsEmpty(value string) bool {
	return value == ""
}

func IsNickname(nickname string, omitempty bool) bool {
	return (omitempty && IsEmpty(nickname)) || nicknameRegexp.MatchString(nickname)
}

func IsEmail(email string, omitempty bool) bool {
	return (omitempty && IsEmpty(email)) || emailRegexp.MatchString(email)
}

func IsFullName(fullname string, omitempty bool) bool {
	return omitempty || !IsEmpty(fullname)
}

func IsSlug(slug string, omitempty bool) bool {
	return (omitempty && IsEmpty(slug)) || slugRegexp.MatchString(slug)
}

func IsTitle(title string, omitempty bool) bool {
	return omitempty || !IsEmpty(title)
}
