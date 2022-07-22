package rapi_user

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-rapi/rapi_core"
	"github.com/maldan/go-rapi/rapi_log"
	"sync"
	"time"
)

var sessionList sync.Map
var userList sync.Map
var saveQueue sync.Map
var dbPath string

var OnUserCreated func(uid string)

func StartSession(db string) {
	dbPath = db

	rapi_log.Info("LOAD SESSION LIST")
	// Load session list
	slist := make([]Session, 0)
	cmhp_file.ReadJSON(dbPath+"/session.json", &slist)
	for _, x := range slist {
		sessionList.Store(x.Id, x)
	}
	rapi_log.Info(fmt.Sprintf("Loaded %v sessions", len(slist)))

	// Session write schedule
	go func() {
		rapi_log.Info("SESSION START SCHEDULER")

		for {
			out := make([]Session, 0)
			sessionList.Range(func(key, value any) bool {
				out = append(out, value.(Session))
				return true
			})
			cmhp_file.Write(dbPath+"/session.json", &out)
			time.Sleep(time.Second * 10)
		}
	}()

	// Saving schedule
	go func() {
		rapi_log.Info("SAVE SCHEDULER")

		for {
			saveQueue.Range(func(key, value any) bool {
				// Check user
				item, ok := userList.Load(key)
				if !ok || item == nil {
					saveQueue.Delete(key)
					return true
				}

				// Save user
				user := item.(*User)
				user.Mu.Lock()
				err := cmhp_file.Write(dbPath+"/user/"+key.(string)+"/info.json", user)
				if err != nil {
					rapi_log.Error(fmt.Sprintf("Saving Error %v", err.Error()))
				} else {
					rapi_log.Info("Saved ")
				}

				user.Mu.Unlock()

				saveQueue.Delete(key)
				return true
			})

			time.Sleep(time.Second * 3)
		}
	}()
}

func SaveUser(uid string) {
	saveQueue.Store(uid, true)
	rapi_log.Info("Saving... " + uid)
}

func SaveUserNow(uid string) {
	v, ok := userList.Load(uid)
	if ok {
		err := cmhp_file.Write(dbPath+"/user/"+uid+"/info.json", v.(*User))
		if err != nil {
			rapi_log.Error(err.Error())
		}
	}
}

func SaveSession(uid string) string {
	session := Session{
		Id:       cmhp_crypto.UID(10),
		UID:      uid,
		Duration: 60 * 60 * 24 * 7,
		Created:  time.Now(),
	}
	sessionList.Store(session.Id, session)
	return session.Id
}

func CreateUser(uid string) *User {
	if cmhp_file.Exists(dbPath + "/user/" + uid + "/info.json") {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Type: "accessDenied", Field: "email", Description: "User already exists"})
	}

	user := User{}
	user.UID = uid
	user.Created = time.Now()
	userList.Store(uid, &user)

	if OnUserCreated != nil {
		OnUserCreated(uid)
	}

	return &user
}

func GetUserByUID(uid string) *User {
	v, ok := userList.Load(uid)
	if ok {
		return v.(*User)
	}

	user := User{}
	err := cmhp_file.ReadJSON(dbPath+"/user/"+uid+"/info.json", &user)
	rapi_log.Info("Load user from file " + uid)

	if err != nil {
		return nil
	}
	userList.Store(uid, &user)
	return &user
}

func GetUserBySession(accessToken string) *User {
	session, ok := sessionList.Load(accessToken)
	if ok {
		user := GetUserByUID(session.(Session).UID)
		if user == nil {
			rapi_core.Fatal(rapi_core.Error{
				Code:        401,
				Type:        "accessDenied",
				Field:       "accessToken",
				Description: fmt.Sprintf("User '%v' not found", session.(Session).UID),
			})
		} else {
			return user
		}
	} else {
		rapi_core.Fatal(rapi_core.Error{
			Code:        401,
			Type:        "accessDenied",
			Field:       "accessToken",
			Description: fmt.Sprintf("Session '%v' not found", accessToken),
		})
	}

	return nil
}

func GetUserData(uid string, kind string, v any) error {
	return cmhp_file.ReadJSON(dbPath+"/user/"+uid+"/"+kind+".json", v)
}

func GetUserDataDir(uid string) string {
	return dbPath + "/user/" + uid
}

func SaveUserData(uid string, kind string, v any) error {
	return cmhp_file.Write(dbPath+"/user/"+uid+"/"+kind+".json", v)
}
