package service

import (
	"echo_sample/db"
)

var dao *db.Queries

func Init() {
	dao = db.Connect()
}
