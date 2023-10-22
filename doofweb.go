package doofweb

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type W2GServer struct {
	DBUrl   string
	Addr    string
	Resp404 string
	Paths   map[string]interface{}
	dbpool  *pgxpool.Pool
}

type W2GViewData struct {
	Resp   http.ResponseWriter
	Req    *http.Request
	DBPool *pgxpool.Pool
}

func (w2gs W2GServer) RunServer() {
	var err error
	fmt.Println("Setting up server")
	if w2gs.DBUrl != "" {
		w2gs.dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
	}
	hndl := w2gs
	s := &http.Server{
		Addr:           ":8080",
		Handler:        hndl,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Listening on :8080")

	log.Fatal(s.ListenAndServe())
}

func (w2gs W2GServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	viewdata := W2GViewData{Resp: resp, Req: req, DBPool: w2gs.dbpool}
	fmt.Print("Recieved request for ")
	fmt.Print(req.URL.Path)
	fmt.Println()
	viewfunc, ok := w2gs.Paths[req.URL.Path]
	if ok {
		viewfunc.(func(*W2GViewData))(&viewdata)
	} else {
		do404(resp, req)
	}
}

func do404(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(404)
	resp.Write([]byte("Error 404: Not Found..."))
}
