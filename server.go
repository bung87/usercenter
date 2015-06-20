package main

import (
    "net/http"
	// "database/sql"
	// "github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
    "github.com/martini-contrib/csrf"
     "github.com/martini-contrib/render"
   "github.com/bung87/usercenter/settings"
        "github.com/martini-contrib/sessionauth"
    "github.com/martini-contrib/sessions"
    "github.com/bung87/usercenter/models"
    "github.com/bung87/usercenter/forms"
    // "fmt"
    _ "github.com/go-sql-driver/mysql"
    "html/template"
        "log"
    // "time"
          "github.com/martini-contrib/binding"
    "golang.org/x/crypto/bcrypt"
    // "net/http"
)

/*func initDB() (*models.DB, error) {
    return models.InitDB("mysql", "go:go@/go", gorp.SqliteDialect{})
}*/

func DB() martini.Handler {
    return func(c martini.Context) {
        db ,err := models.InitDB(settings.Driver,settings.Source,settings.Dialect)
        if err != nil {
            panic(err)
        }
        defer db.Db.Close()
        c.Map(db)
        // c.Map(models.DbMap)
        c.Next()
    }
}




func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
func csrfToken() string {
                    return ""
                }
var funcMap = template.FuncMap{
         "csrfToken" :csrfToken,
        
}

func main() {
  m := martini.Classic()
  m.Use(render.Renderer(render.Options{
     Funcs: []template.FuncMap{
        funcMap,
     },
    }))
    db ,err := models.InitDB(settings.Driver,settings.Source,settings.Dialect)
    if err != nil {
        panic(err)
    }
    db.InitSchema()
    db.Db.Close()
    m.Use(DB())
    store := sessions.NewCookieStore([]byte("secret123"))
    // Default our store to use Session cookies, so we don't leave logged in
    // users roaming around
    store.Options(sessions.Options{
        MaxAge: 0,
    })
    m.Use(sessions.Sessions("my_session", store))
    m.Use(sessionauth.SessionUser(models.GenerateAnonymousUser))

    sessionauth.RedirectUrl = "/new-login"
    sessionauth.RedirectParam = "new-next"
     m.Use(csrf.Generate(&csrf.Options{
        Secret:     "token123",
        SessionKey: "userID",
        SetHeader:true,
        SetCookie:true,
        // Custom error response.
        ErrorFunc: func(w http.ResponseWriter) {
            http.Error(w, "CSRF token validation failed", http.StatusBadRequest)
        },
    }))
     m.Use(func(s sessions.Session,res http.ResponseWriter, req *http.Request) {
        s.Set("userID", "123456")
        })

 m.Post("/signup", csrf.Validate, binding.Bind(forms.SignupForm{}), func(signupForm forms.SignupForm, r render.Render,db *models.DB) {
    if signupForm.Password1 != signupForm.Password2 {
        panic("two password should be matched")
    }
    password :=  []byte(signupForm.Password1)
    // Hashing the password with the cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
    if err != nil {
        panic(err)
    }
   
        u1 := models.NewUser(signupForm.Email, string(hashedPassword))
        
        log.Println(u1)

        err = db.Insert(&u1)
        checkErr(err, "Insert failed")
        
        newmap := map[string]interface{}{"metatitle": "created user", "user": u1}
        r.HTML(200, "user", newmap)
    })

 m.Get("/login",func( s sessions.Session,r render.Render, x csrf.CSRF){
    // 

     log.Println(x.GetToken())
     // err = bcrypt.CompareHashAndPassword(hashedPassword, password)
 	r.HTML(200, "login",x.GetToken())
 	})
 m.Post("/login", csrf.Validate,binding.Bind(forms.LoginForm{}), func(session sessions.Session, loginForm forms.LoginForm, r render.Render,req *http.Request,db *models.DB){
    user := models.User{
        Email:loginForm.Email,
    }
     err := db.SelectOne(&user, "select * from users where email=?", user.Email)
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