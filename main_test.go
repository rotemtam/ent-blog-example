package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rotemtam/ent-blog-example/ent/enttest"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	err := seed(context.Background(), client)
	require.NoError(t, err)

	srv := newServer(client)
	r := newRouter(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Hello, World!")
}

func TestAdd(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	err := seed(context.Background(), client)
	require.NoError(t, err)

	srv := newServer(client)
	r := newRouter(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := ts.Client().PostForm(ts.URL+"/add", map[string][]string{
		"title": {"Testing, one, two."},
		"body":  {"This is a test"},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "This is a test")
}
