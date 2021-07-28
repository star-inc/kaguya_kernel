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

type Interface interface {
	Load(source *RethinkSource, filter ...interface{}) error
	Reload(source *RethinkSource) error
	Create(source *RethinkSource) error
	Update(source *RethinkSource) error
	Replace(source *RethinkSource) error
	Destroy(source *RethinkSource) error
}
