package main_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verdverm/frisby"
)

func TestSomething(t *testing.T) {
	errs := frisby.Create("Test ping").
		Get("http://localhost:5025/ping").
		Send().
		ExpectStatus(http.StatusOK).
		Errors()

	for _, err := range errs {
		assert.Nil(t, err)
	}
}
