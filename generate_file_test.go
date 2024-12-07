package filetools

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFile(t *testing.T) {
	output, err := GenerateFile(context.Background(), FileRequest{
		OutputPath:   "/gentest/test.txt",
		TemplateName: "test",
		Data: struct {
			Name string
		}{Name: "nuzur"},
		Funcs:           nil,
		DisableGoFormat: true,
	})

	assert.NotNil(t, output)
	assert.NoError(t, err)
	rootDir := CurrentPath()
	finalOutput := path.Join(rootDir, "gentest", "test.txt")
	if _, err := os.Stat(finalOutput); errors.Is(err, os.ErrNotExist) {
		assert.NoError(t, err)
	}

	os.RemoveAll(path.Join(rootDir, "gentest"))
}
