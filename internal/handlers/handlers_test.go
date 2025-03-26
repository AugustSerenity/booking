package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AugustSerenity/booking/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"m-s", "/majors-suite", "GET", http.StatusOK},
	{"s-a", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// {"post-search-av", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2025-01-01"},
	// 	{key: "end", value: "2025-01-11"},
	// }, http.StatusOK},
	// {"post-search-av-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2025-01-01"},
	// 	{key: "end", value: "2025-01-11"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "JoJo"},
	// 	{key: "last_name", value: "Onimechkin"},
	// 	{key: "email", value: "Anime@NerYlit.com"},
	// 	{key: "phone", value: "123-21434-1136"},
	// }, http.StatusOK},
}

func TestNewHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with none-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 99
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2030-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@zhuwd.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@zhu.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid end date
	reqBody = "start_date=2030-10-12"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@zhu.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid room id
	reqBody = "start_date=2030-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@zhu.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid room_id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid data
	reqBody = "start_date=2030-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=M")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@ztw.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid data
	reqBody = "start_date=2030-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=P")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zhu@zhu.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid data
	reqBody = "start_date=2030-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Mili")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Pops")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=777777")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
