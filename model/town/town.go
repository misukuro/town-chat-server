// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthtown

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"

	l4g "github.com/alecthomas/log4go"
)

type TownProvider struct {
}

type TownUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &TownProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_TOWN, provider)
}

func userFromTownUser(glu *TownUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	userId := strconv.FormatInt(glu.Id, 10)
	user.AuthData = &userId
	user.AuthService = model.USER_AUTH_SERVICE_TOWN

	return user
}

func townUserFromJson(data io.Reader) *TownUser {
	decoder := json.NewDecoder(data)
	var glu TownUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
    l4g.Error("Town json error err=%v", err.Error())
		return nil
	}
}

func (glu *TownUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *TownUser) IsValid() bool {
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *TownUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *TownProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_TOWN
}

func (m *TownProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := townUserFromJson(data)
	if glu.IsValid() {
		return userFromTownUser(glu)
	}

	return &model.User{}
}

func (m *TownProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := townUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
