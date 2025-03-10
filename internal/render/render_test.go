package render

import (
	"github.com/AugustSerenity/booking/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "000")
	result := AddDefaultData(&td, r)

	if result.Flash != "000" {
		t.Error("flash value of 000 not found in session")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/something-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
