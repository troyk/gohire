package main

import (
	"context"
	"fmt"
	"gohire/proto/gen/api"
	"log"
	"net/http"
	"os"
	"text/template"

	"connectrpc.com/connect"
	"github.com/tailscale/sqlite/cgosqlite"
	"github.com/tailscale/sqlite/sqliteh"
)

func main() {
	dbc, err := cgosqlite.Open("db/gohire.db", sqliteh.SQLITE_OPEN_READWRITE, "")
	if err != nil {
		panic(err)
	}
	templates, err := template.New("").ParseGlob("./web/templates/*.html")
	if err != nil {
		panic(err)
	}

	websvr := &WebServer{
		mux:       http.NewServeMux(),
		api:       NewAPIServer(dbc),
		templates: templates,
	}

	websvr.HandleFunc("GET /", websvr.UsersIndex)

	httpServer := &http.Server{
		Addr:    ":3001",
		Handler: websvr.mux,
	}
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}

	//http://localhost:3001

}

type WebServer struct {
	mux       *http.ServeMux
	api       *APIServer
	templates *template.Template
}

func (s *WebServer) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *WebServer) UsersIndex(w http.ResponseWriter, r *http.Request) {
	apireq := &connect.Request[api.GetUsersRequest]{
		Msg: &api.GetUsersRequest{},
	}
	users, err := s.api.GetUsers(r.Context(), apireq)
	err = s.templates.ExecuteTemplate(w, "index", users.Msg)
	if err != nil {
		panic(err)
	}
}

func NewAPIServer(db sqliteh.DB) *APIServer {
	return &APIServer{db: db}
}

type APIServer struct {
	db sqliteh.DB
}

func (s *APIServer) GetUsers(
	ctx context.Context,
	req *connect.Request[api.GetUsersRequest],
) (*connect.Response[api.GetUsersResponse], error) {
	//log.Println("Request headers: ", req.Header())
	stmt, rq, err := s.db.Prepare("select * from users order by username", sqliteh.SQLITE_PREPARE_NORMALIZE)
	log.Println("select id,username,first_name,last_name from users order by username", rq, err)
	if err != nil {
		return nil, err
	}
	defer stmt.Finalize()
	users := make([]*api.User, 0)
	for {
		ok, err := stmt.Step(nil)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		u := &api.User{
			Id:        stmt.ColumnText(1),
			Username:  stmt.ColumnText(2),
			FirstName: stmt.ColumnText(3),
			LastName:  stmt.ColumnText(4),
		}
		users = append(users, u)
	}

	res := connect.NewResponse(&api.GetUsersResponse{
		Users: users,
	})
	return res, nil
}
