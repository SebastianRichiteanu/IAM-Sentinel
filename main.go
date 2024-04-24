package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	db, err := InitializeDB(ctx, Neo4jConfig{
		uri:      "bolt://localhost:7687",
		username: "neo4j",
		password: "sentinel",
		database: "neo4j"}) // TODO: move this to env params or config or something else :)
	if err != nil {
		panic(err)
	}
	defer db.Stop(ctx)

	if err := db.PurgeDB(ctx); err != nil {
		panic(err)
	}

	what, err := db.ExecuteQueryWrite(ctx, "CREATE (charlie:Person:Actor {name: 'Charlie Sheen'}), (oliver:Person:Director {name: 'Oliver Stone'})", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", what)

	count, err := db.ExecuteQueryRead(ctx, "match(n) return count(n)", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", count[0])
}
