package settings

import (
	// "database/sql"
	"github.com/coopernurse/gorp"
)

const (
    Driver  = "mysql" 
    Source = "go:go@/go" 
    
)

var Dialect gorp.Dialect   = gorp.SqliteDialect{} 