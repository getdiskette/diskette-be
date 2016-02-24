package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/verdverm/frisby"
)

func TestPing(t *testing.T) {
	go start()
	time.Sleep(5 * time.Second)

	errs := frisby.Create("Test ping").
		Get("http://localhost:5025/ping").
		Send().
		ExpectStatus(http.StatusOK).
		ExpectContent("pong").
		AfterText(
		func(F *frisby.Frisby, text string, err error) {
			assert.Nil(t, err)
			assert.Equal(t, "pong", text)
		}).
		Errors()

	for _, err := range errs {
		assert.Nil(t, err)
	}
}
