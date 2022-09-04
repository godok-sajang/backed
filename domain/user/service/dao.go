package service

import (
	"godok/db"
)

var dao *db.Queries

func Init() {
	dao = db.Connect()
}
