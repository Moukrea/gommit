package main

import (
	"strings"
	"testing"
)

func TestSuccessArt(t *testing.T) {
	if len(successArt) == 0 {
		t.Error("successArt is empty")
	}
	if !strings.Contains(successArt, "██") {
		t.Error("successArt does not contain expected ASCII characters")
	}
	if !strings.Contains(successArt, "███") {
		t.Error("successArt does not contain the expected happy face pattern")
	}
	if strings.Count(successArt, "\n") < 10 {
		t.Error("successArt seems too short, expected at least 10 lines")
	}
}

func TestFailureArt(t *testing.T) {
	if len(failureArt) == 0 {
		t.Error("failureArt is empty")
	}
	if !strings.Contains(failureArt, "██") {
		t.Error("failureArt does not contain expected ASCII characters")
	}
	if !strings.Contains(failureArt, "████") {
		t.Error("failureArt does not contain the expected sad face pattern")
	}
	if strings.Count(failureArt, "\n") < 10 {
		t.Error("failureArt seems too short, expected at least 10 lines")
	}
}
