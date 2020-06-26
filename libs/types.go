/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

type Message struct {
	ContentType int
	TargetType int
	Origin string
	Target string
	Content []byte
}

type User struct {
	UserID string
	DisplayName string
}
