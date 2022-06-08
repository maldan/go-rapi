package rapi_user

import (
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_string"
	"github.com/maldan/go-rapi/rapi_core"
	"strings"
)

type UserApi struct {
}

type ArgsSignUp struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type ArgsSignIn struct {
	Kind     string `json:"kind"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PostSignUp Sign Up
func (r UserApi) PostSignUp(args ArgsSignUp) map[string]string {
	// Check email
	email := strings.Trim(strings.ToLower(args.Email), " ")
	if !cmhp_string.IsEmailValid(email) {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "email", Description: "Incorrect email"})
	}

	// Get user
	user := GetUserByUID(email)

	// Check password
	if args.Password1 != args.Password2 {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "password1", Description: "Passwords mismatch!"})
	}
	if len(args.Password1) < 6 {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "password1", Description: "Password must contain at least 6 characters!"})
	}

	// User not found
	if user == nil {
		user = CreateUser(email)
		user.Password = cmhp_crypto.Sha256(args.Password1)
		user.Role = "user"
		if email == "tagir@europe.com" || email == "tagir@botzzup.com" || email == "blackwanted@yandex.ru" {
			user.Role = "admin"
		}
		SaveUserNow(email)

		return map[string]string{
			"accessToken": SaveSession(email),
		}
	} else {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "email", Description: "User already exists"})
	}
	return nil
}

// PostSignIn Sign In
func (r UserApi) PostSignIn(args ArgsSignIn) map[string]string {
	email := strings.Trim(strings.ToLower(args.Email), " ")
	user := GetUserByUID(email)

	if user == nil {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "email", Description: "User not found"})
	}

	if user.UID == email && user.Password == cmhp_crypto.Sha256(args.Password) {
		return map[string]string{
			"accessToken": SaveSession(email),
		}
	} else {
		rapi_core.Fatal(rapi_core.Error{Code: 500, Field: "email", Description: "Incorrect email or password"})
	}
	return nil
}
