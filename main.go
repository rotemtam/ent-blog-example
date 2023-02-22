package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rotemtam/ent-blog-example/ent"
	"github.com/rotemtam/ent-blog-example/ent/user"
)

var (
	//go:embed templates/*
	resources embed.FS
	tmpl      = template.Must(template.ParseFS(resources, "templates/*"))
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
	srv := newServer(client)
	r := newRouter(srv)
	log.Fatal(http.ListenAndServe(":8080", r))
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

type server struct {
	client *ent.Client
}

func newServer(client *ent.Client) *server {
	return &server{client: client}
}

// index serves the blog home page
func (s *server) index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := s.client.Post.
		Query().
		WithAuthor().
		All(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// newRouter creates a new router with the blog handlers mounted.
func newRouter(srv *server) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", srv.index)
	return r
}
