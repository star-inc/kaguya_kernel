/*
Package KaguyaKernel : The kernel for Kaguya

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
package talk

import (
	"context"
	"fmt"
	KaguyaKernel "github.com/star-inc/kaguya_kernel"
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

const MessagesCollection = "messages"

func NewDataInterface() *DataInterface {
	var err error
	dataInterface := new(DataInterface)
	dataInterface.queryTimeout = 50 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	dataInterface.client, err = mongo.Connect(ctx, options.Client().ApplyURI(KaguyaKernel.Config.Database.Host))
	if err != nil {
		panic(err)
	}
	dataInterface.database = dataInterface.client.Database(KaguyaKernel.Config.Database.Name)
	return dataInterface
}

func (dataInterface DataInterface) FetchMessage(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := dataInterface.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (dataInterface DataInterface) SyncMessageBox(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := dataInterface.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (dataInterface DataInterface) GetMessageBox(identity string, target string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := dataInterface.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (dataInterface DataInterface) GetMessage(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := dataInterface.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (dataInterface DataInterface) SaveMessage(message Message) {
	ctx, cancel := context.WithTimeout(context.Background(), dataInterface.queryTimeout)
	defer cancel()
	data, _ := bson.Marshal(&message)
	_, err := dataInterface.database.Collection(MessagesCollection).InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
}
