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
	testbin     = "./tmp/vault-cli-test"
)

func rmTestData() {
	exec.Command("rm", "-rf", testDataDir).Run()
	exec.Command("rm", "-rf", testLogDir).Run()
}

func TestMain(m *testing.M) {
	// remove test data
	rmTestData()

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
	// remove test data
	rmTestData()

	name := ""
	testInfo(t, name)
	testInit(t, name)
	testSetAndGet(t, name)
	testDelete(t, name)
}

func TestNamedVault(t *testing.T) {
	// remove test data
	rmTestData()

	name := "test"
	testInfo(t, name)
	testInit(t, name)
	testSetAndGet(t, name)
	testDelete(t, name)
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
	output, err := initCommand(name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(output))
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

func delete(key, name string) ([]byte, error) {
	if name == "" {
		fmt.Println("delete key")
		return exec.Command(testbin, "delete", key).Output()
	} else {
		fmt.Println("delete key name")
		return exec.Command(testbin, "delete", key, name).Output()
	}
}

func testDelete(t *testing.T, name string) {
	output, err := delete("key", name)
	if err != nil {
		t.Errorf("Error: %v %s", err, output)
	}
	fmt.Println(string(output))
}
