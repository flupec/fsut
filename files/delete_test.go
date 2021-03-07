package files

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestDeleteRecursively_AppliedToFile(t *testing.T) {
	asrt := assert.New(t)
	CreateTestCase(asrt)
	defer AssertDeletedTestCase(asrt)
	defer PurgeTestCase(asrt)

	deleteTargetFilePath := targetFilePath + "/" + filesUnderTarget[0]

	err := DeleteRecursively(deleteTargetFilePath)
	asrt.Nil(err)

	_, err = os.Open(deleteTargetFilePath)
	asrt.True(os.IsNotExist(err))

}

func TestDeleteRecursively_AppliedToDir(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	asrt := assert.New(t)
	CreateTestCase(asrt)
	defer AssertDeletedTestCase(asrt)

	err := DeleteRecursively(targetFilePath)
	asrt.Nil(err)
}

var(
	targetFilePath = os.TempDir() + "/" + "deleteTest"
	dirUnderTarget = targetFilePath + "/" + "dirUnderRoot"

	filesUnderTarget = [2]string{
		"fileUnderRoot1.txt",
		"fileUnderRoot2.txt",
	}
	filesUnderSecondDir = [2]string{
		"fileUnderSecondDir1.txt",
		"fileUnderSecondDir2.txt",
	}
)

func CreateTestCase(asrt *assert.Assertions) {
	_, err := os.Open(targetFilePath)
	asrt.True(os.IsNotExist(err))
	err = os.Mkdir(targetFilePath, 0777)
	asrt.Nil(err)
	for _, f := range filesUnderTarget {
		filename := targetFilePath + "/" + f
		_, err = os.Open(filename)
		asrt.True(os.IsNotExist(err))
		err = ioutil.WriteFile(filename, []byte(f), 0777)
		asrt.Nil(err)
	}

	_, err = os.Open(dirUnderTarget)
	asrt.True(os.IsNotExist(err))
	err = os.Mkdir(dirUnderTarget, 0777)
	asrt.Nil(err)

	for _, f := range filesUnderSecondDir {
		filename := dirUnderTarget + "/" + f
		_, err = os.Open(filename)
		err = ioutil.WriteFile(dirUnderTarget + "/" + f, []byte(f), 0777)
		asrt.Nil(err)
	}
}

func PurgeTestCase(asrt *assert.Assertions) {
	err := os.RemoveAll(dirUnderTarget)
	asrt.Nil(err)

	err = os.RemoveAll(targetFilePath)
	asrt.Nil(err)

	AssertDeletedTestCase(asrt)
}

func AssertDeletedTestCase(asrt *assert.Assertions) {
	for _, f := range filesUnderSecondDir {
		filename := dirUnderTarget + "/" + f
		_, err := os.Open(filename)
		asrt.True(os.IsNotExist(err))
	}
	_, err := os.Open(dirUnderTarget)
	asrt.True(os.IsNotExist(err))

	for _, f := range filesUnderTarget {
		filename := targetFilePath + "/" + f
		_, err := os.Open(filename)
		asrt.True(os.IsNotExist(err))
	}
	_, err = os.Open(targetFilePath)
	asrt.True(os.IsNotExist(err))
}