package render

import (
	"github.com/Sunpacker/go-booking-app/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var templatedata models.TemplateData

	request, err := getRequestWithSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(request.Context(), "flash", "123")

	result := AddDefaultData(&templatedata, request)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}
func getRequestWithSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx, _ = session.Load(ctx, request.Header.Get("X-Session"))
	request = request.WithContext(ctx)

	return request, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestRenderTemplate(t *testing.T) {
	templatesFormat = "../../templates/%s.tmpl"

	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templateCache

	request, err := getRequestWithSession()
	if err != nil {
		t.Error(err)
	}

	var writer skeletonWriter

	err = Template(&writer, request, "home", &models.TemplateData{})
	if err != nil {
		t.Error("[TestRenderTemplate] home rendering error:", err)
	}
	err = Template(&writer, request, "non-existent", &models.TemplateData{})
	if err == nil {
		t.Error("[TestRenderTemplate] non-existent template has rendered")
	}
}

type skeletonWriter struct{}

func (w *skeletonWriter) Header() http.Header {
	var header http.Header
	return header
}
func (w *skeletonWriter) WriteHeader(statusCode int) {
	_ = statusCode
}
func (w *skeletonWriter) Write(bytes []byte) (int, error) {
	length := len(bytes)
	return length, nil
}

func TestCreateTemplateCache(t *testing.T) {
	templatesFormat = "../../templates/%s.tmpl"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
