package filehelp

import (
	"fmt"
	"os"
	"path/filepath"
)

func RelativePath(filename string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s", dir, filename)
}
