package main

// 
// Example schema
// 
// CREATE TABLE benefactors (
//   id SERIAL PRIMARY KEY,
//   seized BOOLEAN
// );
// 


import (
  "context"
  "fmt"

  "github.com/jackc/pgx/v5"
  "github.com/nfb/doofweb"
)

var paths = map[string]doofweb.ViewFunc{
  "/":                         defaultView,
  "/api/v1/benefactor/create": newBenefactor,
  "/api/v1/benefactor/list":   listBenefactors,
}

func defaultView(view *doofweb.W2GViewData) error {
  view.Resp.Write([]byte("defaultresponse"))
  return nil
}

func newBenefactor(view *doofweb.W2GViewData) error {
  _, err := view.DBPool.Exec(context.Background(), "INSERT INTO benefactors (seized) VALUES (FALSE)")
  if err != nil{
    fmt.Println(err)
    return err
  }
  view.Resp.Write([]byte("new benefactor id <>"))
  return nil
}

func listBenefactors(view *doofweb.W2GViewData) error {
  var rows pgx.Rows
  var ben_id int
  rows, err := view.DBPool.Query(context.Background(), "select id from benefactors")

  if err != nil {
    fmt.Println(err)
    return err
  } else {
    pgx.ForEachRow(rows, []any{&ben_id}, func() error {
      fmt.Printf("%v\n", ben_id)
      return nil
    })
  }
  return nil
}

func main() {

  s := doofweb.W2GServer{
    DBUrl:   "postgresql://postgres:example@127.0.0.1:5432/postgres",
    Addr:    ":8080",
    Paths:   paths,
  }
  s.RunServer()

}
