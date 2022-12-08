package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"go_page", "go", "home", false, "error rendering Go template"},
	{"go_page_no_template", "go", "no-file", true, "error rendering non-existent Go template"},
	{"jet_page", "jet", "home", false, "error rendering Jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "error rendering non-existent Jet template"},
	{"invalid_render_engine", "foo", "no-file", true, "no error returned while rendering with invalid engine specified"},
}

func TestRender_Page(t *testing.T) {

	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/some-url", nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()

		testRenderer.RootPath = "./testdata"
		testRenderer.Renderer = e.renderer

		err = testRenderer.Page(w, r, e.template, nil, nil)
		if e.errorExpected {
			if err == nil {
				t.Errorf("%s: %s", e.name, e.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", e.name, e.errorMessage, err.Error())
			}
		}

	}
}

func TestRender_GoPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	testRenderer.RootPath = "./testdata"
	testRenderer.Renderer = "go"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering Go template.", err)
	}

	err = testRenderer.Page(w, r, "no-file", nil, nil)
	if err == nil {
		t.Error("error rendering non-existent go template", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	testRenderer.RootPath = "./testdata"
	testRenderer.Renderer = "jet"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("error rendering Jet template", err)
	}

	err = testRenderer.Page(w, r, "no-file", nil, nil)
	if err == nil {
		t.Error("error rendering non-existent go template", err)
	}
}
