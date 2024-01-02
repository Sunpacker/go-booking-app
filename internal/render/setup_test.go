package render

import (
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = false

	testApp.Session = session
	testApp.UseCache = true
	app = &testApp

	os.Exit(m.Run())
}
