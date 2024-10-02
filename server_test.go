package main

import (
	"os"
	"testing"
)

func TestContactStore(t *testing.T) {
	const filename = "testdata/db.json"
	s := OpenContactStore(filename)
	err := s.SaveContact(&Contact{FirstName: "John"})
	if err != nil {
		t.Fatal(err)
	}
	// reload db from disk
	s = OpenContactStore(filename)
	c := s.Contacts[s.ContactSeqID]
	if c.FirstName != "John" {
		if err != nil {
			t.Fatalf("error db.FirstName[%s]!=John", c.FirstName)
		}
	}
	os.Remove(filename)
}
