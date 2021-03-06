package files

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyWalker(t *testing.T) {
	rootDir := os.TempDir()
	fileUnderRootDir := rootDir + "/" + "file1.txt"
	err := ioutil.WriteFile(fileUnderRootDir, []byte("file1"), 0777)
	assert.Nil(t, err)

	dirUnderRoot := rootDir + "/" + "dir"
	err = os.Mkdir(dirUnderRoot, 0777)
	assert.Nil(t, err)
	t.Log(err)

	fileUnderDirUnderRoot := dirUnderRoot + "/" + "file2.txt"
	err = ioutil.WriteFile(fileUnderDirUnderRoot, []byte("file2"), 0777)
	assert.Nil(t, err)
	t.Log(err)

	copyTo := rootDir + "/" + "copyDest"
	err = os.Mkdir(copyTo, 0777)
	assert.Nil(t, err)
	t.Log(err)

	walker := NewCopyWalker(copyTo, nil)
	err = walker.Walk(dirUnderRoot)
	assert.Nil(t, err)
	t.Log(err)

	_ = os.Remove(fileUnderRootDir)
	_ = os.Remove(fileUnderDirUnderRoot)
	_ = os.RemoveAll(copyTo)
	_ = os.RemoveAll(dirUnderRoot)
}

func TestCopy_SingleDirectory(t *testing.T) {
	copySrcDir := os.TempDir() + "/dir"
	err := os.MkdirAll(copySrcDir, 0777)
	assert.Nil(t, err)

	copyFilesSrc := []string {
		copySrcDir + "/file16.txt",
		copySrcDir + "/file26.txt",
	}
	for _, cpy := range copyFilesSrc {
		err = ioutil.WriteFile(cpy, []byte(cpy), 0777)
		assert.Nil(t, err)
	}

	anotherSrcDir := copySrcDir + "/" + "dir2"
	err = os.MkdirAll(anotherSrcDir, 0777)
	assert.Nil(t, err)

	anotherCopyFilesSrc := []string {
		anotherSrcDir + "/file36.txt",
		anotherSrcDir + "/file46.txt",
	}
	for _, cpy := range anotherCopyFilesSrc {
		err = ioutil.WriteFile(cpy, []byte(cpy), 0777)
		assert.Nil(t, err)
	}

	copyDest := os.TempDir() + "/copyDest"
	err = os.MkdirAll(copyDest, 0777)
	assert.Nil(t, err)

	err = Copy(copyDest, nil, copySrcDir)
	assert.Nil(t, err)

	expectedToBeCopied := []string {
		copyDest + "/file16.txt",
		copyDest + "/file26.txt",
		copyDest + "/dir2/file36.txt",
		copyDest + "/dir2/file46.txt",
	}
	for _, expCopy := range expectedToBeCopied {
		f, err := os.Open(expCopy)
		assert.Nil(t, err)
		assert.Equal(t, expCopy, f.Name())
	}

	// Delete
	err = os.RemoveAll(copyDest)
	assert.Nil(t, err)
	err = os.RemoveAll(copySrcDir)
	assert.Nil(t, err)
}