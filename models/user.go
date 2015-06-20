package models

import (
"time"
"github.com/martini-contrib/sessionauth"
"github.com/bung87/usercenter/settings"
// "log"
// "reflect"
// "fmt"
)
    

type User struct {
    Id      int64 `db:"id"`
    Created int64
    LastLogin int64 
    Username string `form:"Username"`
    Password   string `form:"Password"`
    Email   string `form:"Email"`
    FirstName    string `form:"FirstName"`
    LastName    string `form:"LastName"`
    NickName    string `form:"NickName"`
    authenticated bool   `form:"-" db:"-"`
}
func (u *User) GetById(id interface{}) error {
    db ,err := InitDB(settings.Driver,settings.Source,settings.Dialect)
    if err != nil {
        return err
    }
    defer db.Db.Close()
    
    err = db.SelectOne(&u, "SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        return err
    }

    return nil
}
func NewUser(email, password string) User {
	now := time.Now().Unix()
    return User{
        //Created: time.Now().UnixNano(),
        Password:password,//need to be crypted
        Created: now,
        LastLogin: now,
        Email: email,
    }
}
func (u *User) UniqueId() interface{} {
    return u.Id
}

// GetById will populate a user object from a database model with
// a matching id.

func GenerateAnonymousUser() sessionauth.User {
    return &User{}
}
// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *User) Login() {
    // Update last login time
    // Add to logged-in user's list
    // etc ...
    u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *User) Logout() {
    // Remove from logged-in user's list
    // etc ...
    u.authenticated = false
}

func (u *User) IsAuthenticated() bool {
    return u.authenticated
}