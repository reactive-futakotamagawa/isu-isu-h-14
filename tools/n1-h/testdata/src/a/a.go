package a

import (
	"context"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var tx *sqlx.Tx

type rcv struct{}

func (rcv) dbMethod() {
	db.Get(nil, "SELECT 1")
}

func dbFunc() {
	db.Get(nil, "SELECT 1")
}

func txFunc() {
	tx.Get(nil, "SELECT 1")
}

func dbFuncWrap() {
	dbFunc()
}

func dbSelect() {
	db.Select(nil, "SELECT 1")
}

func dbCtx() {
	db.GetContext(context.Background(), nil, "SELECT 1")

}

func dbDirectInFor() {
	var dst int
	err := db.Get(&dst, "SELECT 1")
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		var dst int
		err := db.Get(&dst, "SELECT ?", i) // want `^maybe N\+1`
		if err != nil {
			return
		}
	}

	for i := range 10 {
		var dst int
		err := db.Get(&dst, "SELECT ?", i) // want `^maybe N\+1`
		if err != nil {
			return
		}
		db.MustBegin().Get(&dst, "SELECT 2") // want `^maybe N\+1`
	}

	for range 10 {
		var dst int
		err := db.Get(&dst, "SELECT 1") // want `^maybe N\+1`
		if err != nil {
			return
		}
	}

	for range 10 {
		var dst int
		err := tx.Get(&dst, "SELECT 1") // want `^maybe N\+1`
		if err != nil {
			return
		}
	}

	for range 10 {
		_, err := db.Exec("SELECT 1") // want `^maybe N\+1`
		if err != nil {
			return
		}
	}
}

func dbFuncInfor() {
	for range 10 {
		dbFunc() // want `^maybe N\+1`
	}

	for i := 0; i < 10; i++ {
		dbDirectInFor() // want `^maybe N\+1`
	}
	var r rcv
	for i := 0; i < 10; i++ {
		r.dbMethod() // want `^maybe N\+1`
	}

	for range 10 {
		dbFuncWrap() // want `^maybe N\+1`
	}

	for range 10 {
		txFunc() // want `^maybe N\+1`
	}

	dbDirectInFor()

	for {
		dbDirectInFor() // want `^maybe N\+1`
		dbFuncInfor()   // want `^maybe N\+1`
	}
}
