package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rotemtam/ent-blog-example/ent"
	"github.com/rotemtam/ent-blog-example/ent/user"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "database DSN")
	flag.Parse()

	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed connecting to mysql: %v", err)
	}
	defer client.Close()
	ctx := context.Background()

	if client.Post.Query().CountX(ctx) == 0 {
		if err := seed(ctx, client); err != nil {
			log.Fatalf("failed seeding the database: %v", err)
		}
	}

	// ... Continue with server start.
}

func seed(ctx context.Context, client *ent.Client) error {
	r, err := client.User.Query().
		Where(
			user.Name("rotemtam"),
		).
		Only(ctx)

	if ent.IsNotFound(err) {
		r, err = client.User.Create().
			SetName("rotemtam").
			SetEmail("r@hello.world").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating user: %v", err)
		}
	}
	return client.Post.Create().
		SetTitle("Hello, World!").
		SetBody("This is my first post").
		SetAuthor(r).
		Exec(ctx)
}
