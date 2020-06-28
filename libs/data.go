/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DataInterface struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewDataInterface() *DataInterface {
	var err error
	dataInterface := new(DataInterface)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	dataInterface.client, err = mongo.Connect(ctx, options.Client().ApplyURI(Config.DBhost))
	DeBug("NewDataInterface", err)
	dataInterface.database = dataInterface.client.Database(Config.DBname)
	return dataInterface
}

func (dataInterface DataInterface) LogMessage(message []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	_, err := dataInterface.database.Collection("kaguya").InsertOne(ctx, bson.M{"name": "message"})
	DeBug("LogMessage", err)
}
