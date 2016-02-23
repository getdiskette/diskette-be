package main_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verdverm/frisby"
)

func TestPing(t *testing.T) {
	errs := frisby.Create("Test ping").
		Get("http://localhost:5025/ping").
		Send().
		ExpectStatus(http.StatusOK).
		ExpectContent("pong").
		Errors()

	for _, err := range errs {
		assert.Nil(t, err)
	}
}
