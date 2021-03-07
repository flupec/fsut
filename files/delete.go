package files

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)


func DeleteRecursively(root string) error {
	// Single file deletion case
	if ok, err := IsDir(root); err == nil {
		if !ok {
			return os.Remove(root)
		}
	} else {
		return err
	}

	delWalker := DeleteWalker{visited: make(map[string]struct{}, 0)} // TODO: Map size estimation?
	return delWalker.Walk(root)
}

type DeleteWalker struct {
	visited map[string]struct{}
}

// Walk performs dfs traverse
func (d DeleteWalker) Walk(root string) error {
	if _, ok := d.visited[root]; !ok {
		if ok, err := IsDir(root); err == nil && ok {
			if filesInDir, err := ioutil.ReadDir(root); err == nil {
				for _, file := range filesInDir {
					if err := d.Walk(root + "/" + file.Name()); err != nil {
						return err
					}
				}
			} else {
				return err
			}
		} else if err != nil {
			return err
		}

		d.visited[root] = struct{}{}
		log.Debugf("%s visited for deletion", root)
		if ok, err := IsDir(root); err != nil {
			return err
		} else if !ok {
			return os.Remove(root)
		} else {
			return os.RemoveAll(root)
		}
	}
	return nil
}