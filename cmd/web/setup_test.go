package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Before starting the tests, do smth inside this function,
	// then run the tests, then exit
	os.Exit(m.Run())
}

type SkeletonHandler struct{}

func (handler *SkeletonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = w
	_ = r
}
