package forms

import (
  "net/http"
      "github.com/martini-contrib/binding"
)

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

