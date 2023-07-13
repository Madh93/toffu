package toffu

import (
	"github.com/Madh93/toffu/internal/woffuapi"
)

type Signs struct {
	signs woffuapi.Signs
}

func (s Signs) hasAlreadyClockedIn() bool {
	if len(s.signs) > 0 {
		return s.signs[len(s.signs)-1].SignIn
	}
	return false
}
