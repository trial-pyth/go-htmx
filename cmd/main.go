package main

import (
    "html/template"
    "io"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

type Template struct {
    tmpl *template.Template
}

func newTemplate() *Template {
    return &Template{
        tmpl: template.Must(template.ParseGlob("views/*.html")),
    }
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}

type Contact struct {
    Name  string
    Email string
}

type Data struct {
    Contacts []Contact
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() FormData {
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

func (d *Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if contact.Email == email {
            return true
        }
    }

    return false
}

func NewData() Data {
    return Data{
        Contacts: []Contact{
            NewContact("john", "jo@gmail.com"),
            NewContact("clara", "cl@gmail.com"),
        },
    }
}

func NewContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}


type Page struct {
    Data Data
    Form FormData
}

func newPage() Page {
    return Page{
        Data: NewData(),
        Form: newFormData(),
    }
}

func main() {

    e := echo.New()

    page := newPage()

    e.Renderer = newTemplate()
    e.Use(middleware.RequestLogger())

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index.html", page)
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists"

            return c.Render(422, "form", formData)
        }

        contact := NewContact(name, email)
        page.Data.Contacts = append(page.Data.Contacts, contact)

        c.Render(200, "form", newFormData())
        return c.Render(200, "display", page.Data)
    })

    e.Logger.Fatal(e.Start(":42069"))
}