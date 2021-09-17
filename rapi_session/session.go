package rapi_core

import (
	"sync"
	"time"

	"github.com/maldan/go-cmhp/cmhp_crypto"
)

var sessionList sync.Map

type Session struct {
	Id       string    `json:"id"`
	Key      string    `json:"key"`
	Duration int       `json:"duration"`
	Created  time.Time `json:"created"`
}

func Start() {

}

func Create(key string) string {
	session := Session{
		Id:       cmhp_crypto.UID(10),
		Key:      key,
		Duration: 60 * 60 * 24 * 7,
		Created:  time.Now(),
	}
	sessionList.Store(session.Id, session)
	return session.Id
}
