package main

import (
    "net/http"
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
     "github.com/martini-contrib/render"
     "github.com/martini-contrib/binding"
        "github.com/martini-contrib/sessionauth"
    "github.com/martini-contrib/sessions"
    // "fmt"
    _ "github.com/go-sql-driver/mysql"
    // "html/template"
        "log"
    "time"
    "golang.org/x/crypto/bcrypt"
)

var dbmap *gorp.DbMap

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
func (u *User) UniqueId() interface{} {
    return u.Id
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *User) GetById(id interface{}) error {
    err := dbmap.SelectOne(u, "SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        return err
    }

    return nil
}
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

type SignupForm struct{
	Password1 string `form:"Password1" binding:"required"`
	Password2 string `form:"Password2" binding:"required"`
    Email   string `form:"Email" binding:"required"`
}

type LoginForm struct{
    Password string `form:"Password" binding:"required"`
    Email   string `form:"Email" binding:"required"`
}
func (cf SignupForm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
    /*if strings.Contains(cf.Email, "Go needs generics") {
        errors = append(errors, binding.Error{
            FieldNames:     []string{"Email"},
            Classification: "ComplaintError",
            Message:        "Go has generics. They're called interfaces.",
        })
    }*/
    return errors
}
func newUser(email, password string) User {
	now := time.Now().Unix()
    return User{
        //Created: time.Now().UnixNano(),
        Password:password,//need to be crypted
        Created: now,
        LastLogin: now,
        Email: email,
    }
}
func initDb() *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    
    //db, err := sql.Open("sqlite3", "/tmp/post_db.bin")
    db, err := sql.Open("mysql", "go:go@/go")
    checkErr(err, "sql.Open failed")

    // construct a gorp DbMap
    // dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())
   // initialize the DbMap
    dbmap := initDb()
    defer dbmap.Db.Close()
    store := sessions.NewCookieStore([]byte("secret123"))
    // Default our store to use Session cookies, so we don't leave logged in
    // users roaming around
    store.Options(sessions.Options{
        MaxAge: 0,
    })
    m.Use(sessions.Sessions("my_session", store))
    m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
    sessionauth.RedirectUrl = "/new-login"
    sessionauth.RedirectParam = "new-next"
 m.Post("/signup", binding.Bind(SignupForm{}), func(signupForm SignupForm, r render.Render) {
    if signupForm.Password1 != signupForm.Password2 {
        panic("two password should be matched")
    }
    password :=  []byte(signupForm.Password1)
    // Hashing the password with the cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
    if err != nil {
        panic(err)
    }
   
        u1 := newUser(signupForm.Email, string(hashedPassword))
        
        log.Println(u1)

        err = dbmap.Insert(&u1)
        checkErr(err, "Insert failed")
        
        newmap := map[string]interface{}{"metatitle": "created user", "user": u1}
        r.HTML(200, "user", newmap)
    })

 m.Get("/login",func( r render.Render){
     // err = bcrypt.CompareHashAndPassword(hashedPassword, password)
 	r.HTML(200, "login","233")
 	})
 m.Post("/login",binding.Bind(LoginForm{}), func(session sessions.Session, loginForm LoginForm, r render.Render,req *http.Request){
    user := User{
        Email:loginForm.Email,
    }
     err := dbmap.SelectOne(&user, "select * from users where email=?", user.Email)
    checkErr(err, "SelectOne failed")
    if err != nil {
            r.Redirect(sessionauth.RedirectUrl)
            return
        } 
    hashedPassword := []byte(user.Password)
    password := []byte(loginForm.Password)
    err = bcrypt.CompareHashAndPassword(hashedPassword, password)
    checkErr(err, "Password match failed")
    err = sessionauth.AuthenticateSession(session, &user)
            if err != nil {
                r.JSON(500, err)
            }
            params := req.URL.Query()
            redirect := params.Get(sessionauth.RedirectParam)
            r.Redirect(redirect)
            return
    r.HTML(200,"success","")
    })
  m.Get("/", func(r render.Render)  {
      r.HTML(200, "hello", "jeremy")
  })
  m.Run()
}