package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/glamour"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

var (
	ErrBeautifyMD = errors.New("error beautifying markdown")
	ErrGetAlias   = errors.New("error getting alias")
	ErrGetProject = errors.New("error getting project")
	ErrReadREADME = errors.New("error reading README")
)

func TitleCase(s string) string {
	t := cases.Title(language.English)
	return t.String(s)
}

func FilterByStatuses(data map[string]string, status []string) map[string]string {
	filteredData := make(map[string]string)
	for k, v := range data {
		for _, s := range status {
			if v == s {
				filteredData[k] = v
			}
		}
	}
	return filteredData
}

// GetLastModifiedTime returns the last modified time
func GetLastModifiedTime(dbname, pname string) string {
	var lastModTime time.Time
	pPath, err := db.GetRecord(dbname, pname, c.ProjectPaths)
	if err != nil {
		return "Something went wrong"
	}
	_ = filepath.Walk(pPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.ModTime().After(lastModTime) {
			lastModTime = info.ModTime()
		}
		return nil
	})
	return fmt.Sprint(lastModTime.Format("02 Jan 06 15:04"))
}

// BeautifyMD: returns styled markdown
func BeautifyMD(data []byte) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
		glamour.WithAutoStyle(),
	)
	if err != nil {
		log.Print("something went wrong while creating renderer: ", err)
		return "", errors.Join(ErrBeautifyMD, err)
	}
	out, _ := r.Render(string(data))
	return out, nil
}

// GetProjectPath: returns the path to the README inside the project
func GetProjectPath(dbname, projectName string) (string, error) {
	path, err := db.GetRecord(dbname, projectName, c.ProjectPaths)
	if err != nil {
		actualName, err := db.GetRecord(dbname, projectName, c.ProjectAliasBucket)
		if err != nil {
			log.Printf("project: %v not a valid project\n", projectName)
			return "", errors.Join(ErrGetAlias, err)
		}
		projectName = actualName
		path, err = db.GetRecord(dbname, projectName, c.ProjectPaths)
		if err != nil {
			log.Printf("project: %v not a valid project\n", projectName)
			return "", errors.Join(ErrGetProject, err)
		}
	}
	pPath := filepath.Join(path, "README.md")
	_, err = os.Stat(pPath)
	return pPath, err
}

// ReadREADME: returns the byte array of README.md of a project
func ReadREADME(dbname, projectName string) ([]byte, error) {
	pPath, err := GetProjectPath(dbname, projectName)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("# README does not exist for %s", projectName)
			return []byte(msg), nil
		}
		return nil, err
	}
	data, err := os.ReadFile(pPath)
	if err != nil {
		return nil, errors.Join(ErrReadREADME, fmt.Errorf("something went wrong while reading README for %s: %w", projectName, err))
	}
	return data, nil
}

func UpdateLastEditedTime() error {
	r := fmt.Sprint(time.Now().Unix())
	rec := map[string]string{"lastRefreshTime": r}
	err := db.WriteToDB(db.DBName, rec, c.ConfigBucket)
	if err != nil {
		return err
	}
	return nil

}

func DayPassed(t string) bool {
	oneDay := 86400
	now := time.Now().Unix()
	recTime, _ := strconv.ParseInt(t, 10, 64)
	return now-recTime >= int64(oneDay)
}

func ParseTime(tstr string) (string, int64) {
	layout := "02 Jan 06 15:04"
	p, err := time.Parse(layout, tstr)
	timeStamp := p.Unix()
	if err != nil {
		return "unnkown", 0
	}
	today := time.Now()
	switch fmt.Sprint(p.Date()) {
	case fmt.Sprint(today.Date()):
		return fmt.Sprintf("Today %s", p.Format("15:04")), timeStamp
	case fmt.Sprint(today.AddDate(0, 0, -1).Date()):
		return fmt.Sprintf("Yesterday %s", p.Format("17:00")), timeStamp
	}
	return tstr, timeStamp
}
