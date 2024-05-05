package main

import (
	"os"
	"strings"
	"testing"
)

const testRoot = "tmp/test"

func split(s string) []string {
	return strings.Split(s, " ")
}

func setup() {
	os.RemoveAll(testRoot)
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}
