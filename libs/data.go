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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataInterface struct {
	client       *mongo.Client
	database     *mongo.Database
	queryTimeout time.Duration
}

const (
	Collection_Access   = "accesses"
	Collection_Users    = "users"
	Collection_Messages = "messages"
)

func NewDataInterface() *DataInterface {
	var err error
	dataInterface := new(DataInterface)
	dataInterface.queryTimeout = 50 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	dataInterface.client, err = mongo.Connect(ctx, options.Client().ApplyURI(Config.DBhost))
	DeBug("NewDataInterface", err)
	dataInterface.database = dataInterface.client.Database(Config.DBname)
	return dataInterface
}

func (dataInterface DataInterface) GetAccess(username string, password string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"username": username, "password": password}
	_ = dataInterface.database.Collection(Collection_Users).FindOne(ctx, filter).Decode(&result)
	return result
}

func (dataInterface DataInterface) RegisterAccess(identity string, authToken string) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	data := bson.M{"identity": identity, "authToken": authToken}
	_, err := dataInterface.database.Collection(Collection_Access).InsertOne(ctx, data)
	if err != nil {
		return false
	}
	return true
}

func (dataInterface DataInterface) VerfiyAccess(authToken string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"authToken": authToken}
	_ = dataInterface.database.Collection(Collection_Access).FindOne(ctx, filter).Decode(&result)
	return result
}

func (dataInterface DataInterface) RegisterUser(user User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	data, _ := bson.Marshal(user)
	_, err := dataInterface.database.Collection(Collection_Users).InsertOne(ctx, data)
	if err != nil {
		return false
	}
	return true
}

func (dataInterface DataInterface) LogMessage(message *Message) {
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	data, _ := bson.Marshal(&message)
	_, err := dataInterface.database.Collection(Collection_Messages).InsertOne(ctx, data)
	DeBug("LogMessage", err)
}
