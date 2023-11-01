package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	A  int                `bson:"a"`
	B  string             `bson:"b"`
}

func main() {
	url := "mongodb://localhost:27017"
	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	c := client.Database("picop-example").Collection("example")

	res, err := c.InsertOne(ctx, &Data{A: 123, B: "example"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)

	cur, err := c.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	data := []Data{}
	err = cur.All(ctx, &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
