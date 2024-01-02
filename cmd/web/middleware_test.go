package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var skeletonHandler SkeletonHandler

	handler := NoSurf(&skeletonHandler)

	switch typeReceived := handler.(type) {
	case http.Handler:
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", typeReceived))
	}
}

func TestSessionLoad(t *testing.T) {
	var skeletonHandler SkeletonHandler

	handler := SessionLoad(&skeletonHandler)

	switch typeReceived := handler.(type) {
	case http.Handler:
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", typeReceived))
	}
}
