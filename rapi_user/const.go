package rapi_user

import (
	"sync"
	"time"
)

type Session struct {
	Id       string    `json:"id"`
	UID      string    `json:"uid"`
	Duration int       `json:"duration"`
	Created  time.Time `json:"created"`
}

type User struct {
	Mu sync.Mutex `json:"-"`

	// Locked
	UID      string `json:"uid"`
	Password string `json:"password"`
	Role     string `json:"role"`

	Created time.Time `json:"created"`
}
