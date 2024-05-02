package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/theredditbandit/pman/pkg/db"
)

const StatusBucket = "projects"
const ProjectPaths = "projectPaths"
const ProjectAliasBucket = "projectAliases"

// InitDirs indexes a directory for project directories and writes the data to the db
func InitDirs(args []string) {
	// the file which identifies a project directory
	projIdentifier := "README.md"
	if len(args) != 1 {
		fmt.Println("Please provide a directory name")
		return
	}
	dirname := args[0]
	if stat, err := os.Stat(dirname); os.IsNotExist(err) {
		fmt.Printf("%s is not a directory \n", dirname)
		return
	} else if !stat.IsDir() {
		fmt.Printf("%s is a file and not a directory \n", dirname)
		return
	}
	projDirs, err := indexDir(dirname, projIdentifier)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Indexed %d project directories . . .\n", len(projDirs))
	projectStatusMap := make(map[string]string)
	projectPathMap := make(map[string]string)
	for k, v := range projDirs { // k : full project path, v : project status ,
		projectStatusMap[filepath.Base(k)] = v // filepath.Base(k) : project name
		projectPathMap[filepath.Base(k)] = k
	}
	err = db.WriteToDB(projectStatusMap, StatusBucket)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = db.WriteToDB(projectPathMap, ProjectPaths)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// indexDir indexes a directory for project directories
func indexDir(path, identifier string) (map[string]string, error) {
	projDirs := make(map[string]string)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			pname := filepath.Dir(path)
			absPath, _ := filepath.Abs(pname)
			projDirs[absPath] = "indexed"
			return filepath.SkipDir
		}
		if !info.IsDir() && info.Name() == identifier {
			pname := filepath.Dir(path)
			absPath, _ := filepath.Abs(pname)
			projDirs[absPath] = "indexed"
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return projDirs, nil
}
