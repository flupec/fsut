package files

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Copy copies files at filePaths to 'to' directory. 'to' directory must already exist
// strategy param can be nil (will be replaced with always rewrite strategy)
func Copy(to string, strategy CopyStrategy, filePaths ...string) error {
	if strategy == nil {
		strategy = AlwaysRewriteCopyStrategy{}
	}
	var aggregatedError string

	for _, filePath := range filePaths {
		if err := copyFile(to, filePath, strategy); err != nil {
			if len(aggregatedError) == 0 {
				aggregatedError = err.Error()
			} else {
				aggregatedError = fmt.Sprintf("%s\n%s", aggregatedError, err.Error())
			}
		}
	}

	if len(aggregatedError) > 0 {
		return fmt.Errorf("cannot copy file. Errors: %s", aggregatedError)
	}

	return nil
}

func copyFile(to string, filePath string, strategy CopyStrategy) error {
	var copySrc *os.File
	var err error

	if copySrc, err = os.Open(filePath); err != nil {
		return err
	}

	if ok, _ := IsDir(copySrc.Name()); ok {
		walker := NewCopyWalker(to, strategy)
		if err := walker.Walk(copySrc.Name()); err != nil {
			return err
		}
	} else {
		if err := copySingleFile(to, copySrc, strategy); err != nil {
			return err
		}
	}
	return nil
}

type Walker interface {
	Walk(root string) error
}

type CopyWalker struct {
	copyTo        string
	alreadyCopied map[string]struct{}
	root          string
	copyStrategy  CopyStrategy
}

func (cw *CopyWalker) Walk(root string) error {
	cw.root = root
	return filepath.Walk(root, cw.walk)
}

func (cw *CopyWalker) walk(path string, info os.FileInfo, err error) error {
	log.Infof("Walk at %s", path)
	if copySrc, err := os.Open(path); err == nil {
		isDir, _ := IsDir(path)
		if !cw.isAlreadyCopied(path) && !isDir {
			if err := copySingleFile(cw.diffRelativePath(copySrc.Name()), copySrc, cw.copyStrategy); err != nil {
				return err
			}
			cw.addToCopied(path)
		}
	} else {
		return filepath.SkipDir
	}
	return nil
}

func (cw CopyWalker) diffRelativePath(copySrcFilePath string) string {
	copySrcDir, _ := filepath.Split(copySrcFilePath)
	copyDestDir := strings.Split(copySrcDir, cw.root)[1]
	return cw.copyTo + copyDestDir
}

func (cw *CopyWalker) addToCopied(filepath string) {
	cw.alreadyCopied[filepath] = struct{}{}
}

func (cw CopyWalker) isAlreadyCopied(filepath string) bool {
	_, ok := cw.alreadyCopied[filepath]
	return ok
}

func NewCopyWalker(copyTo string, strategy CopyStrategy) Walker {
	if strategy == nil {
		strategy = AlwaysRewriteCopyStrategy{}
	}
	return &CopyWalker{
		copyTo:        copyTo,
		alreadyCopied: make(map[string]struct{}, 0),
		copyStrategy:  strategy,
	}
}

func copySingleFile(to string, copySrc *os.File, strategy CopyStrategy) error {
	copyFilepath := to + "/" + filepath.Base(copySrc.Name())
	fileBytes, err := ioutil.ReadAll(copySrc)
	if err != nil {
		return err
	}
	cleanCopyFilePath := filepath.Clean(copyFilepath)
	if err := os.MkdirAll(filepath.Dir(cleanCopyFilePath), 0777); err != nil {
		return err
	}

	_, err = os.Open(cleanCopyFilePath)
	fileExists := os.IsExist(err)

	if fileExists {
		if strategy.CopyWithRewrite() {
			if err := ioutil.WriteFile(cleanCopyFilePath, fileBytes, 0666); err != nil {
				return err
			}
		}
	} else {
		if err := ioutil.WriteFile(cleanCopyFilePath, fileBytes, 0666); err != nil {
			return err
		}
	}
	return nil
}

func IsDir(filename string) (bool, error) {
	if fInfo, err := os.Stat(filename); err == nil {
		return fInfo.IsDir(), nil
	} else {
		return false, err
	}
}

// CopyStrategy can decide if particular file must be rewritten (for example, based on user input)
type CopyStrategy interface {
	CopyWithRewrite() bool
}

type AlwaysRewriteCopyStrategy struct {}

func (a AlwaysRewriteCopyStrategy) CopyWithRewrite() bool {
	return true
}

type NeverRewriteCopyStrategy struct {}

func (n NeverRewriteCopyStrategy) CopyWithRewrite() bool {
	return false
}

