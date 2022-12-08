package session

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func TestSession_InitSession(t *testing.T) {
	s := &Session{
		CookieName:     "    goracoon",
		CookieLifeTime: "60",
		CookiePersist:  "true",
		CookieDomain:   "localhost",
		CookieSecure:   "true",
		SessionType:    "cookie",
	}

	var sm *scs.SessionManager
	session := s.InitSession()

	var sessionKind reflect.Kind
	var sessionType reflect.Type

	rv := reflect.ValueOf(session)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("For loop:", rv.Kind(), rv.Type(), rv)
		sessionKind = rv.Kind()
		sessionType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("invalid type or kind type:", rv.Type(), "kind:", rv.Kind())
	}

	if sessionKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned testing cookie session. Expected", reflect.ValueOf(sm).Kind(), "and got", sessionKind)
	}

	if sessionType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type returned testing cookie session. Expected", reflect.ValueOf(sm).Kind(), "and got", sessionType)
	}
}
