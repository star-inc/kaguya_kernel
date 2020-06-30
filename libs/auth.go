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

	"github.com/google/uuid"
)

func (handle *Handle) GetAccess(username string, password string) []byte {
	authorization := handle.dataInterface.GetAccess(username, password)
	if authorization != nil {
		tokenSeed := uuid.New().String()
		tokenHandle := sha512.New()
		tokenHandle.Write([]byte(tokenSeed))
		return tokenHandle.Sum(nil)
	}
	return []byte{}
}

func (handle *Handle) RegisterUser(displayName string, username string, password string) bool {
	data := NewDataInterface()
	var user User
	user.Identity = uuid.New().String()
	user.DisplayName = displayName
	user.Username = username
	user.Password = password
	return data.RegisterUser(user)
}
