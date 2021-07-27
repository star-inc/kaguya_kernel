/*
Package KaguyaKernel: The kernel for Kaguya

    Copyright 2021 Star Inc.(https://starinc.xyz)

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
package data

// Message: Message is the data structure from client, to be created the message by including into Container.
type Message struct {
	Content     string `rethinkdb:"content" json:"content"`
	ContentType int    `rethinkdb:"contentType" json:"contentType"`
	Origin      string `rethinkdb:"origin" json:"origin"`
}