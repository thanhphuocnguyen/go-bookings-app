package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var testCases = []struct {
	name             string
	url              string
	method           string
	params           []postData
	expectStatusCode int
}{
	{"home", "/", "GET", []postData{}, 200},
	{"about", "/about", "GET", []postData{}, 200},
	{"contact", "/contact", "GET", []postData{}, 200},
	{"search-availability", "/search-availability", "GET", []postData{}, 200},
	{"show-login", "/login", "GET", []postData{}, 200},
	{"show-registration", "/register", "GET", []postData{}, 200},
}

func TestHandlers(t *testing.T) {
	// arrange
	routes := getRoutes()
	// httptest.NewTLSServer is used to create a new test server with TLS
	ts := httptest.NewTLSServer(routes)

	// defer is used to close the server after the test is done
	defer ts.Close()

	for _, tt := range testCases {
		// t.Run is used to run subtests and sub-benchmarks in a single test
		t.Run(tt.name, func(t *testing.T) {
			// act
			var res *http.Response
			var err error
			// ts.Client() is used to create a new client to make requests to the server
			client := ts.Client()
			if tt.method == "GET" {
				res, err = client.Get(ts.URL + tt.url)
			} else if tt.method == "POST" {
				values := url.Values{}
				for _, v := range tt.params {
					values.Add(v.key, v.value)
				}
				res, err = client.PostForm(ts.URL+tt.url, values)
			} else {
				t.Fatal("Unsupported method")
			}
			if err != nil {
				t.Fatal(err)
			}
			// assert
			if res.StatusCode != tt.expectStatusCode {
				t.Errorf("for %s, expected %d but got %d", tt.name, tt.expectStatusCode, res.StatusCode)
			}
		})
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		UserId: 1,
		RoomId: 1,
		Room: models.Room{
			Name:        "General's Quarters",
			ID:          2,
			Price:       123,
			Description: "description",
			Slug:        "generals-quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// arrange
	rr := httptest.NewRecorder()

	appConfig.Session.Put(ctx, "reservation", reservation)

	// act
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	reservation.RoomId = 2
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	appConfig.Session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_CreateReservation(t *testing.T) {
	startDate, _ := time.Parse(layout, "2021-01-01")
	endDate, _ := time.Parse(layout, "2021-01-02")
	reservation := models.Reservation{
		RoomId:    1,
		StartDate: startDate,
		EndDate:   endDate,
	}
	postData := url.Values{}
	postData.Add("first_name", "Thanh Phuoc")
	postData.Add("last_name", "Nguyen")
	postData.Add("email", "testing@example.com")
	postData.Add("phone", "123456789123")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	appConfig.Session.Put(ctx, "reservation", reservation)
	// arrange
	rr := httptest.NewRecorder()

	// act
	handler := http.HandlerFunc(Repo.CreateReservation)
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusSeeOther {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// missing body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	appConfig.Session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// missing reservation in session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid email, missing phone, missing last name, missing first name
	postData.Set("email", "invalid")
	postData.Del("phone")
	postData.Del("last_name")
	postData.Del("first_name")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	appConfig.Session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid room id
	postData.Set("email", "testing@xample.com")
	postData.Set("phone", "1234567893")
	postData.Set("last_name", "Nguyen")
	postData.Set("first_name", "Thanh Phuoc")

	reservation.RoomId = 2

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	appConfig.Session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// cannot insert room restriction

	reservation.RoomId = 1000

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	appConfig.Session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("CreateReservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	postData := url.Values{}
	postData.Add("start", "2021-01-01")
	postData.Add("end", "2021-01-02")
	postData.Add("room_id", "1")

	// test case for valid room id
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, want %d", responseRecorder.Code, http.StatusOK)
	}

	// test case for invalid room id
	postData.Set("room_id", "2")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)
	var j jsonResponse

	err := json.Unmarshal(responseRecorder.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("AvailabilityJSON returned OK when it should not have")
	}

	// form invalid
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, want %d", responseRecorder.Code, http.StatusBadRequest)
	}

	if j.OK {
		t.Error("AvailabilityJSON returned OK when it should not have")
	}

	// missing start date
	postData.Del("start")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, want %d", responseRecorder.Code, http.StatusBadRequest)
	}

	if j.Message != "Cannot parse dates" {
		t.Error("AvailabilityJSON returned wrong message")
	}

	// missing end date
	postData.Set("start", "2021-01-01")
	postData.Del("end")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, want %d", responseRecorder.Code, http.StatusBadRequest)
	}

	if j.Message != "Cannot parse dates" {
		t.Error("AvailabilityJSON returned wrong message")
	}

	// missing room id
	postData.Set("end", "2021-01-02")
	postData.Del("room_id")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, req)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &j)

	if err != nil {
		t.Error("failed to parse json")
	}

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("AvailabilityJSON handler returned wrong response code: got %d, want %d", responseRecorder.Code, http.StatusBadRequest)
	}

	if j.Message != "Cannot parse room id" {
		t.Error("AvailabilityJSON returned wrong message")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		ID:        1,
		UserId:    1,
		RoomId:    1,
		FirstName: "Thanh Phuoc",
		LastName:  "Nguyen",
		Email:     "testing@example.com",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 1),
		Phone:     "123123123123",
		Room: models.Room{
			Name:        "General's Quarters",
			ID:          2,
			Price:       123,
			Description: "description",
			Slug:        "generals-quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// arrange
	rr := httptest.NewRecorder()

	appConfig.Session.Put(ctx, "reservation", reservation)

	// act
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		ID:        1,
		UserId:    1,
		FirstName: "Thanh Phuoc",
		LastName:  "Nguyen",
		Email:     "testing@example.com",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 1),
		Phone:     "123123123123",
		Room: models.Room{
			Name:        "General's Quarters",
			ID:          2,
			Price:       123,
			Description: "description",
			Slug:        "generals-quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	appConfig.Session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, want %d", rr.Code, http.StatusOK)
	}

	// missing reservation in session
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// missing room id
	req, _ = http.NewRequest("GET", "/choose-room/abc", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := appConfig.Session.Load(req.Context(), req.Header.Get("X-Session"))

	if err != nil {
		return ctx
	}

	return ctx
}
