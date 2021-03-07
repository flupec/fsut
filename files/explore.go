package files

import (
	"errors"
	"io/ioutil"
)

var ErrCannotExploreFilesFromFile = errors.New("cannot list files in regular file")

func ExploreFiles(root string) ([]FileProperties, error) {
	if ok, err := IsDir(root); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrCannotExploreFilesFromFile
	}

	if filesInDir, err := ioutil.ReadDir(root); err != nil {
		return nil, err
	} else {
		fileProps := make([]FileProperties, len(filesInDir))
		for i, fileInDir := range filesInDir {
			if fileProps[i], err = GetFileProperties(root + "/" + fileInDir.Name()); err != nil {
				return nil, err
			}
		}
		return fileProps, nil
	}
}