/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"github.com/google/uuid"
)

const (
	ContentType_Text = 0
	TargetType_Contact = 0
)

type Message struct {
	ContentType int
	TargetType  int
	Origin      string
	Target      string
	Content     []byte
}

type User struct {
	UserID      string
	DisplayName string
}

func NewUserInfo(displayName string) User {
	return User{UserID: uuid.New().String(), DisplayName: displayName}
}
