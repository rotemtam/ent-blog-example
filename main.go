package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/rotemtam/ent-blog-example/ent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rotemtam/ent-blog-example/ent/user"
)

func main() {
	// Read the connection string to the database from a CLI flag.
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "database DSN")
	flag.Parse()

	// Instantiate the Ent client.
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed connecting to mysql: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	// If we don't have any posts yet, seed the database.
	if client.Post.Query().ExistX(ctx) {
		if err := seed(ctx, client); err != nil {
			log.Fatalf("failed seeding the database: %v", err)
		}
	}
	// ... Continue with server start.
}

func seed(ctx context.Context, client *ent.Client) error {
	// Check if the user "rotemtam" already exists.
	r, err := client.User.Query().
		Where(
			user.Name("rotemtam"),
		).
		Only(ctx)
	switch {
	// If not, create the user.
	case ent.IsNotFound(err):
		r, err = client.User.Create().
			SetName("rotemtam").
			SetEmail("r@hello.world").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating user: %v", err)
		}
	case err != nil:
		return fmt.Errorf("failed querying user: %v", err)
	}
	// Finally, create a "Hello, world" blogpost.
	return client.Post.Create().
		SetTitle("Hello, World!").
		SetBody("This is my first post").
		SetAuthor(r).
		Exec(ctx)
}
