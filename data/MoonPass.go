// Copyright 2021 Star Inc.(https://starinc.xyz)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"github.com/cloudflare/circl/dh/sidh"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"io"
	"log"
)

type Moonpass struct {
	ID        uint8
	RNG       io.Reader
	PublicKey []byte
}

func NewMoonpass() Interface {
	instance := new(Moonpass)
	return instance
}

func (m Moonpass) CheckReady() bool {
	panic("implement me")
}

func (m Moonpass) Load(source KernelSource.Interface, filter ...interface{}) error {
	panic("implement me")
}

func (m Moonpass) Reload(source KernelSource.Interface) error {
	panic("implement me")
}

func (m Moonpass) Create(source KernelSource.Interface) error {
	panic("implement me")
}

func (m Moonpass) Replace(source KernelSource.Interface) error {
	panic("implement me")
}

func (m Moonpass) Destroy(source KernelSource.Interface) error {
	panic("implement me")
}

func (m Moonpass) Encrypt(data []byte) []byte {
	var output []byte
	publicKey := sidh.NewPublicKey(m.ID, sidh.KeyVariantSike)
	if err := publicKey.Import(m.PublicKey); err != nil {
		log.Panicln(err)
	}
	pair := sidh.NewSike751(m.RNG)
	if err := pair.Encapsulate(data, output, publicKey); err != nil {
		log.Panicln(err)
	}
	return output
}
