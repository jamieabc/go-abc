package main_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	program     = "go-abc"
	fixtureFile = "fixtures.go"
)

func TestMainWhenBadArguments(t *testing.T) {
	var err error
	cmd := exec.Command(program)
	out, err := cmd.CombinedOutput()
	assert.Equal(t, nil, err, "wrong error")
	assert.Contains(t, string(out), "Usage", "wrong output")
}

func TestMainWhenFileNotExist(t *testing.T) {
	var err error
	cmd := exec.Command(program, "not-exist-file")
	out, err := cmd.CombinedOutput()
	assert.Equal(t, nil, err, "wrong error")
	assert.Contains(t, string(out), "open", "wrong output")
}

func TestMainWhenDirNotExist(t *testing.T) {
	var err error
	cmd := exec.Command(program, "not-exist-dir")
	out, err := cmd.CombinedOutput()
	assert.Equal(t, nil, err, "wrong error")
	assert.Contains(t, string(out), "open", "wrong output")
}

func TestMainWhenDirNotEmpty(t *testing.T) {
	var err error
	cmd := exec.Command(program, "fixtures")
	out, err := cmd.CombinedOutput()
	sout := string(out)
	assert.Equal(t, nil, err, "wrong error")
	assert.Contains(t, sout, fixtureFile, "wrong list")
}

func TestMainFixture(t *testing.T) {
	var err error
	cmd := exec.Command(program, "fixtures")
	out, err := cmd.CombinedOutput()
	sout := string(out)
	assert.Equal(t, nil, err, "wrong error")
	assert.Contains(t, sout, fixtureFile, "wrong list")
	assert.Contains(t, sout, "assignment: 1", "wrong assignment")
	assert.Contains(t, sout, "branch: 2", "wrong branch")
	assert.Contains(t, sout, "condition: 2", "wrong condition")
	assert.Contains(t, sout, "score: 3", "wrong score")
}
