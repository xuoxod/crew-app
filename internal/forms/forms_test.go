package forms

import (
	"net/http/httptest"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("expected true but received false")
	}
}
