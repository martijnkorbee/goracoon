package render

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Secure          bool
	Port            string
	ServerName      string
}

func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.Port = c.Port
	td.CSRFToken = nosurf.Token(r)
	// user authenticated
	if c.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
	}

	return td
}

func (c *Render) Page(
	w http.ResponseWriter,
	r *http.Request,
	view string,
	variables interface{},
	data interface{},
) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:
	}

	return errors.New("no rendering engine specified")
}

// GoPage renders a template using the GO templating engine
func (c *Render) GoPage(
	w http.ResponseWriter,
	r *http.Request,
	view string,
	data interface{},
) error {
	tmpl, err := template.ParseFiles(c.RootPath + "/views/" + view + ".page.tmpl")
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	if err = tmpl.Execute(w, &td); err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet templating engine
func (c *Render) JetPage(
	w http.ResponseWriter,
	r *http.Request,
	view string,
	variables interface{},
	data interface{},
) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	// add default template data
	td = c.defaultData(td, r)

	t, err := c.JetViews.GetTemplate(view + ".jet")
	if err != nil {
		return err
	}

	if err = t.Execute(w, vars, &td); err != nil {
		return err
	}

	return nil
}
