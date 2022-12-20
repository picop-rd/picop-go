package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// _ "github.com/go-sql-driver/mysql" bcopmysql内でimportされるので不要
	"github.com/hiroyaonoe/bcop-go/contrib/github.com/go-sql-driver/mysql/bcopmysql"
	bcopprop "github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
	"go.opentelemetry.io/otel/baggage"
	otelprop "go.opentelemetry.io/otel/propagation"
)

func main() {
	// 伝播されたContextを用意
	bag := TestBaggage()
	h := header.NewV1(bag.String())
	ctx := otelprop.Baggage{}.Extract(context.Background(), bcopprop.NewBCoPCarrier(h))

	bcopmysql.RegisterDialContext("tcp", otelprop.Baggage{})

	db, err := sql.Open("mysql", "root:@tcp(localhost:9000)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(0)

	_, err = db.ExecContext(ctx, "INSERT INTO books(id, name) VALUES (1, \"test1\")")
	if err != nil {
		fmt.Printf("exec error: %s\n", err.Error())
		return
	}

	rows, err := db.QueryContext(ctx, "SELECT id, name FROM books")
	if err != nil {
		fmt.Printf("query error: %s\n", err.Error())
		return
	}

	for rows.Next() {
		id := 0
		name := ""
		if err := rows.Scan(&id, &name); err != nil {
			fmt.Printf("scan error: %s\n", err.Error())
			continue
		}
		fmt.Printf("book{ id: %d, name: %s }\n", id, name)
	}
}

func TestBaggage() baggage.Baggage {
	m1p1, _ := baggage.NewKeyProperty("p1Key")
	m1p2, _ := baggage.NewKeyValueProperty("p2Key", "p2Value")
	m1, _ := baggage.NewMember("m1Key", "m1Value", m1p1, m1p2)
	m2, _ := baggage.NewMember("m2Key", "m2Value")
	m3, _ := baggage.NewMember("env-id", "aaaaa")
	b, _ := baggage.New(m1, m2, m3)
	return b
}
