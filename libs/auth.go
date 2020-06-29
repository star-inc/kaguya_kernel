/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

func GetAccess(username string, password string) bool {
	return true
}

func RegisterUser(username string, password string) bool {
	data := NewDataInterface()
	user := NewUserInfo(username)
	return 	data.RegisterUser(user)
}