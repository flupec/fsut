package files

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExploreFiles(t *testing.T) {
	tempFilePath := os.TempDir()

	fileProps, err := ExploreFiles(tempFilePath)
	assert.Nil(t, err)
	assert.NotNil(t, fileProps)
	t.Logf("%+v", fileProps)

	notDirFileProp, ok := GetFirstNonDirFile(fileProps)
	if ok {
		t.Logf("Non dir fileProp is %+v", notDirFileProp)
		fileProperty, err := ExploreFiles(tempFilePath + "/" + notDirFileProp.Name)
		assert.Nil(t, fileProperty)
		assert.NotNil(t, err)
		assert.Equal(t, ErrCannotExploreFilesFromFile, err)
	}
}

func GetFirstNonDirFile(fileProps []FileProperties) (FileProperties, bool) {
	for _, fileProp := range fileProps {
		if fileProp.Type == File {
			return fileProp, true
		}
	}
	return FileProperties{}, false
}