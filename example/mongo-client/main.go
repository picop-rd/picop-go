package main

import (
	"context"
	"fmt"
	"log"

	"github.com/picop-rd/picop-go/contrib/go.mongodb.org/mongo-driver/mongo/picopmongo"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"

	// Please use github.com/picop-rd/mongo-go-driver instead of go.mongodb.org/mongo-driver by replacing in go.mod
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	A  int                `bson:"a"`
	B  string             `bson:"b"`
}

func main() {
	// Prepare propagated context
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewPiCoPCarrier(h))

	uri := "mongodb://localhost:27017"
	client := picopmongo.New(options.Client().ApplyURI(uri))

	conn, err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	c := conn.Database("picop-example").Collection("example")

	res, err := c.InsertOne(ctx, &Data{A: 123, B: "example"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.InsertedID)
	conn.Disconnect(ctx)

	conn, err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cur, err := c.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	conn.Disconnect(ctx)

	data := []Data{}
	err = cur.All(ctx, &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
