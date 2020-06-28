/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataInterface struct {
	client *mongo
}

func NewDataInterface() {
	dataInterface = new(DataInterface)
	dataInterface.client, err := mongo.NewClient(
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", Config.DBhost))
	)
	DeBug("NewDataInterface", err);
}
