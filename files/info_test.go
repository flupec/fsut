package files

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetFileProperties_RegularFile(t *testing.T) {
	testFilePath := os.TempDir() + "/TestGetFileProperties_RegularFile.txt"
	err := ioutil.WriteFile(testFilePath, []byte("AAAAAAAAAAAAAAAAAAAAAAA"), 0666)
	assert.Nil(t, err)

	fileProps, err := GetFileProperties(testFilePath)
	assert.Nil(t, err)
	t.Logf("%+v", fileProps)

	err = os.Remove(testFilePath)
	assert.Nil(t, err)
}

func TestGetFileProperties_Directory(t *testing.T) {
	testFilePath := os.TempDir() + "/TestGetFileProperties_Directory"
	err := os.Mkdir(testFilePath, 0777)
	assert.Nil(t, err)

	fileProps, err := GetFileProperties(testFilePath)
	assert.Nil(t, err)
	t.Logf("%+v", fileProps)

	assert.Equal(t, Directory, fileProps.Type)

	err = os.Remove(testFilePath)
	assert.Nil(t, err)
}