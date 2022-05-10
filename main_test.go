package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicense(t *testing.T) {
	assert := assert.New(t)
	var out bytes.Buffer
	err := bufferAndRun([]string{"test", "-L"}, &out, &bytes.Buffer{})
	assert.NoError(err)
	output := string(out.Bytes())
	assert.Equal(license, output)
}

func TestKind(t *testing.T) {
	assert := assert.New(t)
	var out bytes.Buffer
	err := bufferAndRun([]string{"test", "-k"}, &out, &bytes.Buffer{})
	assert.NoError(err)
	output := strings.TrimSpace(string(out.Bytes()))
	assert.Equal("Plain", output)
}
