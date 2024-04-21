package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	testDataDir = "./test_dst/data"
	testLogDir  = "./test_dst/log"
	testbin     = "./bin/vault-cli-test"
)

func TestMain(m *testing.M) {
	// remove test data
	exec.Command("rm", "-rf", testDataDir).Run()
	exec.Command("rm", "-rf", testLogDir).Run()

	// set env
	os.Setenv("VAULT_DIR", testDataDir)
	os.Setenv("VAULT_LOG_DIR", testLogDir)
	retCode := m.Run()
	os.Exit(retCode)
}

func TestHelp(t *testing.T) {
	output, err := exec.Command(testbin, "help").Output()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(output))
}

func testInfo(t *testing.T, name string) {
	output, err := info(name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(output))
}

func TestDefaultName(t *testing.T) {
	name := ""
	testInfo(t, name)
	testInit(t, name)
	testSetAndGet(t, name)
}

func TestNamedVault(t *testing.T) {
	name := "test"
	testInfo(t, name)
	testInit(t, name)
	testSetAndGet(t, name)
}

func initCommand(name string) ([]byte, error) {
	if name == "" {
		return exec.Command(testbin, "init").Output()
	} else {
		return exec.Command(testbin, "init", name).Output()
	}
}

func testInit(t *testing.T, name string) {
	checkInfo(t, false, name)
	t.Log("must be success")
	output, err := initCommand(name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(output))
	t.Log("must be error: already exists")
	output, err = initCommand(name)
	if err == nil {
		t.Errorf("Error: must be error")
	}
	fmt.Println(string(output))
	checkInfo(t, true, name)
}

func info(name string) ([]byte, error) {
	if name == "" {
		return exec.Command(testbin, "info").Output()
	} else {
		return exec.Command(testbin, "info", name).Output()
	}
}

func checkInfo(t *testing.T, init bool, name string) {
	output, err := info(name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	str := string(output)
	if init {
		if !strings.Contains(str, "init: true") {
			t.Errorf("Error: must be init: true")
		}
	} else {
		if !strings.Contains(str, "init: false") {
			t.Errorf("Error: must be init: false")
		}
	}
	fmt.Println(str)
}

func get(key, name string) ([]byte, error) {
	if name == "" {
		return exec.Command(testbin, "get", key).Output()
	} else {
		return exec.Command(testbin, "get", key, name).Output()
	}
}

func set(key, value, name string) ([]byte, error) {
	if name == "" {
		return exec.Command(testbin, "set", key, value).Output()
	} else {
		return exec.Command(testbin, "set", key, value, name).Output()
	}
}

func testSetAndGet(t *testing.T, name string) {
	output, err := set("key", "TestSetAndGet", name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(output))

	output, err = get("key", name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if !strings.Contains(string(output), "TestSetAndGet") {
		t.Errorf("Error: %v", string(output))
	}

	fmt.Println(string(output))
}
