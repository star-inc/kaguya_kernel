/*
Package kaguya : The library for kaguya

    Copyright 2020 Star Inc.(https://starinc.xyz)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
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
