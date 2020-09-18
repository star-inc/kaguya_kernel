/*
Package kaguya : The library for kaguya

    Copyright 2020 Star Inc.(https://starinc.xyz)

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
package kaguya

import (
	"context"
	"fmt"
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

func (dataInterface DataInterface) CheckUserExisted(identity string) bool {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"identity": identity}
	_ = dataInterface.database.Collection(Collection_Users).FindOne(ctx, filter).Decode(&result)
	if result != nil {
		return true
	}
	return false
}

func (dataInterface DataInterface) LogMessage(message *Message) {
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	data, _ := bson.Marshal(&message)
	_, err := dataInterface.database.Collection(Collection_Messages).InsertOne(ctx, data)
	DeBug("LogMessage", err)
}

func (dataInterface DataInterface) GetMessageBox(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{bson.M{"target": identity}, bson.M{"origin": identity}}}
	cursor, _ := dataInterface.database.Collection(Collection_Messages).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}
