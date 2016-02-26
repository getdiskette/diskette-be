package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/verdverm/frisby"
)

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	// teardown()
	os.Exit(retCode)
}

func setup() {
	go main()
	time.Sleep(1 * time.Second)
}

func TestPing(t *testing.T) {
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
