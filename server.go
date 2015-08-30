package main

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"github.com/usercenter/csrf"
	"github.com/usercenter/usercenter/forms"
	"github.com/usercenter/usercenter/models"
	"github.com/usercenter/usercenter/settings"
	"log"
	"net/http"

	"image/jpeg"
	"io"
)

const (
	CAPTCHA_LENGTH  = 6
	CAPTCHA_WIDTH   = 180
	CAPTCHA_HEIGHT  = 50
	CAPTCHA_QUALITY = 90
)

func DB() martini.Handler {
	return func(c martini.Context) {
		db, err := models.InitDB(settings.Driver, settings.Source, settings.Dialect)
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

/*var funcMap = template.FuncMap{
 "csrfToken" :csrfToken,
}*/

func main() {
	m := martini.Classic()
	m.Use(render.Renderer( /*render.Options{
	  Funcs: []template.FuncMap{
	     funcMap,
	     },
	     }*/))
	db, err := models.InitDB(settings.Driver, settings.Source, settings.Dialect)
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
		SetHeader:  true,
		SetCookie:  true,
		// Custom error response.
		ErrorFunc: func(w http.ResponseWriter) {
			http.Error(w, "CSRF token validation failed", http.StatusBadRequest)
		},
	}))
	m.Use(func(s sessions.Session, res http.ResponseWriter, req *http.Request) {
		s.Set("userID", "123456")
	})

	m.Post("/signup", csrf.Validate, binding.Bind(forms.SignupForm{}), func(user sessionauth.User, res http.ResponseWriter, req *http.Request, signupForm forms.SignupForm, r render.Render, db *models.DB) {

		if !captcha.VerifyString(req.FormValue("captchaId"), req.FormValue("captchaSolution")) {
			io.WriteString(res, "Wrong captcha solution! No robots allowed!\n")
		} else {
			io.WriteString(res, "Great job, human! You solved the captcha.\n")
		}

		if signupForm.Password1 != signupForm.Password2 {
			panic("two password should be matched")
		}
		// password :=  []byte(signupForm.Password1)
		// Hashing the password with the cost of 10
		// hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
		// user.(*models.User)
		if err != nil {
			panic(err)
		}

		u1 := models.NewUser(signupForm.Email)

		log.Println(u1)

		err = db.Insert(&u1)
		checkErr(err, "Insert failed")

		newmap := map[string]interface{}{"metatitle": "created user", "user": u1}
		r.HTML(200, "user", newmap)
	})

	m.Get("/login", func(s sessions.Session, r render.Render, x csrf.CSRF) {

		d := struct {
			CaptchaId string
			Token     string
		}{x.GetToken(),
			captcha.New(),
		}

		r.HTML(200, "login", d)
	})

	m.Post("/login", csrf.Validate, binding.Bind(forms.LoginForm{}), func(session sessions.Session, loginForm forms.LoginForm, r render.Render, req *http.Request, db *models.DB) {
		user := models.User{
			Email: loginForm.Email,
		}
		err := db.SelectOne(&user, "select * from users where email=?", user.Email)
		checkErr(err, "SelectOne failed")
		if err != nil {
			r.Redirect(sessionauth.RedirectUrl)
			return
		}
		// hashedPassword := []byte(user.Password)
		// password := []byte(loginForm.Password)
		// err = bcrypt.CompareHashAndPassword(hashedPassword, password)
		err = user.Authenticate(loginForm.Password)
		checkErr(err, "Password match failed")
		err = sessionauth.AuthenticateSession(session, &user)
		if err != nil {
			r.JSON(500, err)
		}
		params := req.URL.Query()
		redirect := params.Get(sessionauth.RedirectParam)
		r.Redirect(redirect)
		return
		r.HTML(200, "success", "")
	})

	m.Get("/", func(r render.Render) {
		r.HTML(200, "hello", "jeremy")
	})

	m.Get("/captcha/:id", func(res http.ResponseWriter, req *http.Request) {
		digits := captcha.RandomDigits(CAPTCHA_LENGTH)
		value := ""
		for _, d := range digits {
			value += fmt.Sprintf("%v", d)
		}
		image := captcha.NewImage("", digits, CAPTCHA_WIDTH, CAPTCHA_HEIGHT)
		err := jpeg.Encode(res, image, &jpeg.Options{Quality: CAPTCHA_QUALITY})
		if err != nil {
			res.WriteHeader(500)
		} else {
			res.WriteHeader(200)
		}
	})

	m.Run()
}
