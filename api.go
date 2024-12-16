package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Api struct {
	logger *slog.Logger
	router *http.ServeMux
}

func newApi(logger *slog.Logger) *Api {
	mux := http.NewServeMux()

	return &Api{
		logger: logger,
		router: mux,
	}
}

func (a *Api) Start(ctx context.Context) error {

	a.PostsRouter(ctx)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     a.router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	fmt.Printf("Starting server on port %d...\n", port)

	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

type ModelPost struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

func (a *Api) PostsRouter(ctx context.Context) {
	a.router.HandleFunc("GET /posts", a.FindPosts)
}

func (a *Api) FindPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := a.logger.With("method", "FindPosts")

	resp := struct {
		Results []*ModelPost `json:"results"`
	}{
		Results: []*ModelPost{
			{
				Title:  "Blog 1",
				Body:   "Text 1",
				Author: "Author 1",
			},
			{
				Title:  "Blog 2",
				Body:   "Text 2",
				Author: "Author 2",
			},
			{
				Title:  "Blog 3",
				Body:   "Text 3",
				Author: "Author 3",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.ErrorContext(ctx, "error happened in JSON write", "error", err)
		w.Write([]byte("Error"))
	}
	if _, err := w.Write(jsonResp); err != nil {
		log.ErrorContext(ctx, "error writing response", "error", err)
		w.Write([]byte("Error"))
		return
	}

	log.InfoContext(ctx,
		"success find posts",
		"number_of_posts",
		len(resp.Results),
	)
	return

}
