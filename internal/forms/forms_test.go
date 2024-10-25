package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("got valid when required fields are missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r = httptest.NewRequest("POST", "/whatever", nil)

	r.PostForm = postData

	form = New(r.PostForm)

	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("got invalid when required fields are provided")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	form.MinLength("x", 4, r)

	if form.Valid() {
		t.Error("Expected form to be invalid")
	}

	postData := url.Values{}
	postData.Add("x", "test")

	r.PostForm = postData
	form = New(r.PostForm)

	if !form.Valid() {
		t.Error("Expected form to be valid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("got valid when should have been invalid")
	}

	err := form.Errors.Get("email")
	if err == "" {
		t.Error("should have an error, but did not get one")
	}

	postData := url.Values{}
	postData.Add("email", "testing@example.com")

	r = httptest.NewRequest("POST", "/whatever", nil)

	r.PostForm = postData

	form = New(r.PostForm)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}

	err = form.Errors.Get("email")
	if err != "" {
		t.Error("should not have an error, but got one")
	}
}
