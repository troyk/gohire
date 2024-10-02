package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	tpl, err := template.New("views").ParseFS(templateFS,
		"templates/header.go.html",
		"templates/footer.go.html",
		"templates/edit.go.html",
	)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Lookup("edit.go.html").Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	log.Println("http://0.0.0.0:8089")
	log.Fatal(http.ListenAndServe(":8089", nil))

}

func OpenContactStore(filename string) *ContactStore {
	s := &ContactStore{
		filename: filename,
		Contacts: make(map[int64]Contact),
	}
	err := s.Load()
	if err != nil {
		panic(err)
	}
	return s
}

//go:embed templates/*.go.html
var templateFS embed.FS

type Contact struct {
	ID        int64
	FirstName string
}

type ContactStore struct {
	filename     string
	Contacts     map[int64]Contact
	ContactSeqID int64
	Lock         sync.Mutex
}

func (s *ContactStore) GetContacts() []Contact {
	//TODO
	return nil
}

func (s *ContactStore) SaveContact(c *Contact) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if c.ID == 0 {
		s.ContactSeqID++
		c.ID = s.ContactSeqID
	}
	s.Contacts[c.ID] = *c
	return s.Save()
}

func (s *ContactStore) Load() error {
	f, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return (err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(s)
	if err != nil && err.Error() != "EOF" {
		return (err)
	}
	return nil
}

func (s *ContactStore) Save() error {
	f, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return (err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(s)
	if err != nil && err.Error() != "EOF" {
		return (err)
	}
	return nil
}
