/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"crypto/sha512"
	"encoding/base64"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (Handler *Handler) GetAccess(username string, password string) []byte {
	authorization := Handler.dataInterface.GetAccess(username, password)
	if authorization != nil {
		tokenSeed := uuid.New().String()
		tokenHandle := sha512.New()
		tokenHandle.Write([]byte(tokenSeed))
		authToken := tokenHandle.Sum(nil)
		func(authorization interface{}, authToken []byte) {
			queryResult := authorization.(primitive.D)
			encodedAuthToken := base64.StdEncoding.EncodeToString(authToken)
			Handler.dataInterface.RegisterAccess(queryResult.Map()["identity"].(string), encodedAuthToken)
		}(authorization, authToken)
		return authToken
	}
	return []byte{}
}

func (Handler *Handler) VerfiyAccess(authToken string) {
	data := NewDataInterface()
	verified := data.VerfiyAccess(authToken)
	if verified == nil {
		Handler.identity = ""
		return
	}
	verifiedData := verified.(primitive.D)
	Handler.identity = verifiedData.Map()["identity"].(string)
}

func (Handler *Handler) RegisterUser(displayName string, username string, password string) bool {
	data := NewDataInterface()
	var user User
	user.Identity = uuid.New().String()
	user.DisplayName = displayName
	user.Username = username
	user.Password = password
	return data.RegisterUser(user)
}
