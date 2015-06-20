package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	// "github.com/bung87/usercenter/settings"
)
// var DbMap *gorp.DbMap
// func init(){
// 	InitDB(settings.Driver,settings.Source,settings.Dialect)
// }
type DB struct {
	*gorp.DbMap
}

func InitDB(driver string, source string, dialect gorp.Dialect) (*DB, error){
	db,err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	dbMap := &gorp.DbMap{
		Db:      db,
		Dialect: dialect,
	}
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	return &DB{dbMap}, nil
}

func (db *DB) InitSchema() error {
	if err := db.CreateTablesIfNotExists(); err != nil {
		return err
	}
	return nil
}