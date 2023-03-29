package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// _ "github.com/go-sql-driver/mysql" picopmysql内でimportされるので不要
	"github.com/picop-rd/picop-go/contrib/github.com/go-sql-driver/mysql/picopmysql"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

func main() {
	// 伝播されたContextを用意
	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, "aaaaa")
	ctx := propagation.EnvID{}.Extract(context.Background(), propagation.NewPiCoPCarrier(h))

	picopmysql.RegisterDialContext("tcp", propagation.EnvID{})

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
