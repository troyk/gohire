package main

import (
	"log"
	"testing"

	"github.com/matryer/is"
	"github.com/tailscale/sqlite/cgosqlite"
	"github.com/tailscale/sqlite/sqliteh"
)

func TestContactUpdate(t *testing.T) {
	is := is.New(t)
	dbc, err := cgosqlite.Open("db/gohire.db", sqliteh.SQLITE_OPEN_READWRITE, "")
	is.NoErr(err)
	defer dbc.Close()

	stmt, _, err := dbc.Prepare("update users set password=:password where id=@id", sqliteh.SQLITE_PREPARE_NORMALIZE)
	is.NoErr(err)
	defer stmt.Finalize()

	is.Equal(stmt.BindParameterIndexSearch("password"), 1)
	is.Equal(stmt.BindParameterIndexSearch("id"), 2)

	stmt.BindText64(1, "linux2024")
	stmt.BindText64(2, "018f094d3f9f77fcbc2b735444c5017b")

	row, lastInsertRowID, changes, dur, err := stmt.StepResult()
	log.Println(row, lastInsertRowID, changes, dur, err)

}
