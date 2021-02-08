/*
Package Kernel : The kernel for Kaguya

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
	Kernel "github.com/star-inc/kaguya_kernel"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	client       *mongo.Client
	database     *mongo.Database
	queryTimeout time.Duration
}

const MessagesCollection = "messages"

func NewData() *Data {
	var err error
	data := new(Data)
	data.queryTimeout = 50 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), data.queryTimeout)
	data.client, err = mongo.Connect(ctx, options.Client().ApplyURI(Kernel.Config.Database.Host))
	if err != nil {
		panic(err)
	}
	data.database = data.client.Database(Kernel.Config.Database.Name)
	return data
}

func (data Data) FetchMessage(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), data.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := data.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (data Data) SyncMessageBox(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), data.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := data.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (data Data) GetMessageBox(identity string, target string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), data.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := data.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (data Data) GetMessage(identity string) interface{} {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), data.queryTimeout)
	defer cancel()
	filter := bson.M{"$or": []bson.M{{"target": identity}, {"origin": identity}}}
	cursor, _ := data.database.Collection(MessagesCollection).Find(ctx, filter)
	_ = cursor.All(ctx, &result)
	fmt.Println(result)
	return result
}

func (data Data) SaveMessage(message Message) {
	ctx, cancel := context.WithTimeout(context.Background(), data.queryTimeout)
	defer cancel()
	one, _ := bson.Marshal(&message)
	_, err := data.database.Collection(MessagesCollection).InsertOne(ctx, one)
	if err != nil {
		panic(err)
	}
}
