package models

import (
	"github.com/golang/crypto/bcrypt"
	"github.com/martini-contrib/sessionauth"
	"github.com/usercenter/usercenter/settings"
	"github.com/wayn3h0/go-uuid/random"
	"time"
	// "log"
	// "reflect"
	// "fmt"
)

type User struct {
	Id            int64 `db:"id"`
	Created       int64
	LastLogin     int64
	Username      string `form:"Username"`
	Password      string `form:"Password"`
	Email         string `form:"Email"`
	FirstName     string `form:"FirstName"`
	LastName      string `form:"LastName"`
	NickName      string `form:"NickName"`
	authenticated bool   `form:"-" db:"-"`
	IsActive      bool
}

type Registration struct {
	Id             int64 `db:"id"`
	Created        int64
	ActivationCode string
	Uid            int64
}

type PasswordReset struct {
	Id        int64 `db:"id"`
	Created   int64
	ResetCode string
	Uid       int64
}

func NewPasswordReset(uid int64) PasswordReset {
	now := time.Now().Unix()
	code, _ := random.New()
	return PasswordReset{
		Created:   now,
		ResetCode: string(code),

		Uid: uid,
	}
}
func (r *PasswordReset) Reset(password string) error {
	db, err := InitDB(settings.Driver, settings.Source, settings.Dialect)
	if err != nil {
		return err
	}
	defer db.Db.Close()
	u := User{}
	err = db.SelectOne(&u, "SELECT * FROM users WHERE id = ?", r.Uid)
	if err != nil {
		return err
	}
	// passwordBytes:=[]byte(password)
	// hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, 10)
	// u.Password =string(hashedPassword)
	u.SetPassword(password)
	_, err1 := db.Update(u)
	if err1 != nil {
		return err1
	}

	return nil

}

func (r *Registration) Active() error {
	db, err := InitDB(settings.Driver, settings.Source, settings.Dialect)
	if err != nil {
		return err
	}
	defer db.Db.Close()
	u := User{}
	err = db.SelectOne(&u, "SELECT * FROM users WHERE id = ?", r.Uid)
	if err != nil {
		return err
	}

	u.IsActive = true
	_, err1 := db.Update(u)
	if err1 != nil {
		return err1
	}

	return nil

}

func NewRegistration(uid int64) Registration {
	now := time.Now().Unix()
	code, _ := random.New()
	return Registration{
		Created:        now,
		ActivationCode: string(code),
		Uid:            uid,
	}
}
func (u *User) GetById(id interface{}) error {
	db, err := InitDB(settings.Driver, settings.Source, settings.Dialect)
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
func NewUser(email string) User {
	now := time.Now().Unix()
	return User{
		//Created: time.Now().UnixNano(),

		Created:   now,
		LastLogin: now,
		Email:     email,
	}
}
func (u *User) UniqueId() interface{} {
	return u.Id
}

func (u *User) Authenticate(password string) error {
	hashedPassword := []byte(u.Password)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPassword, passwordBytes)
}

func (u *User) SetPassword(password string) {
	passwordBytes := []byte(password)
	// Hashing the password with the cost of 10
	hashedPassword, _ := bcrypt.GenerateFromPassword(passwordBytes, 10)
	u.Password = string(hashedPassword)

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
