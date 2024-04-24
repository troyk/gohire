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
	websvr.HandleFunc("POST /api/users", websvr.PostUsersHandler)

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

func (s *WebServer) PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	password := r.Form.Get("password")

  postUsersRequest := &api.PostUsersRequest{
		Username:     username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}

	apireq := &connect.Request[api.PostUsersRequest]{Msg: postUsersRequest}

  _, err = s.api.PostUsers(r.Context(), apireq)
	if err != nil {
    http.Error(w, "Failed to add user", http.StatusInternalServerError)
    return
	}
  //fmt.Fprint(w, "successfully added user")
  http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (s *APIServer) PostUsers(
	ctx context.Context,
	req *connect.Request[api.PostUsersRequest],
) (*connect.Response[api.PostUsersResponse], error) {
  userData := req.Msg

  stmt, _, err := s.db.Prepare("INSERT INTO users (username, first_name, last_name, password) VALUES (?, ?, ?, ?)", sqliteh.SQLITE_PREPARE_NORMALIZE)
	if err != nil {
		return nil, err
	}
	defer stmt.Finalize()

	if err := stmt.BindText64(1, userData.Username); err != nil {
		return nil, err
	}
  if err := stmt.BindText64(2, userData.FirstName); err != nil {
		return nil, err
	}
	if err := stmt.BindText64(3, userData.LastName); err != nil {
		return nil, err
	}
	if err := stmt.BindText64(4, userData.Password); err != nil {
		return nil, err
	}

	if _, err := stmt.Step(nil); err != nil {
		return nil, err
	}

  res := connect.NewResponse(&api.PostUsersResponse{})
	return res, nil
}
