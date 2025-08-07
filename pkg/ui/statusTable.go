package ui

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

var RenderTable = func(dbname string, data map[string]string, refreshLastEditedTime bool) error {
	var tableData [][]string
	var lastEdited string
	var timestamp int64

	if refreshLastEditedTime {
		err := utils.UpdateLastEditedTime()
		if err != nil {
			return err
		}
	} else {
		rec, err := db.GetRecord(dbname, "lastRefreshTime", c.ConfigBucket)
		if err != nil { // lastRefreshTime key does not exist in db
			refreshLastEditedTime = true
			err := utils.UpdateLastEditedTime()
			if err != nil {
				return err
			}
		}
		if utils.DayPassed(rec) { // lastEdited values are more than a day old. need to refresh them
			refreshLastEditedTime = true
			err := utils.UpdateLastEditedTime()
			if err != nil {
				return err
			}
		}
	}
	for p, status := range data {
		// log.Println("===============================================")
		alias, err := db.GetRecord(dbname, p, c.ProjectAliasBucket)
		// log.Println("after get record func", "  err=", err)
		// log.Println("p=", p)
		// log.Println("status=", status)
		if refreshLastEditedTime {
			t := utils.GetLastModifiedTime(dbname, p)
			lastEdited, timestamp = utils.ParseTime(t)
			rec := map[string]string{p: fmt.Sprintf("%s-%d", lastEdited, timestamp)}
			err := db.WriteToDB(dbname, rec, c.LastUpdatedBucket)
			if err != nil {
				return err
			}
		} else {
			// log.Println("b4 key nf error", "p=", p, "bucket=", c.LastUpdatedBucket)
			lE, err := db.GetRecord(dbname, p, c.LastUpdatedBucket)
			if err != nil {
				log.Println("aftert key nf error")
				return err
			}
			out := strings.Split(lE, "-")
			lastEdited = out[0]
			tstmp, err := strconv.ParseInt(out[1], 10, 64)
			if err != nil {
				return err
			}
			timestamp = tstmp
		}
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", p, alias)
			row := []string{utils.TitleCase(status), pname, lastEdited, fmt.Sprint(timestamp)} // Status | projectName (alias) | lastEdited | timestamp
			tableData = append(tableData, row)
		} else {
			row := []string{utils.TitleCase(status), p, lastEdited, fmt.Sprint(timestamp)} // Status | projectName | lastEdited | timestamp
			tableData = append(tableData, row)
		}
	}
	if len(tableData) == 0 {
		fmt.Printf("No projects found in the database\n\n")
		fmt.Printf("Add projects to the database using \n\n")
		fmt.Println("pman init .")
		fmt.Println("or")
		fmt.Println("pman add <projectDir>")
		return fmt.Errorf("no database initialized")
	}
	sort.Slice(tableData, func(i, j int) bool {
		valI, _ := strconv.ParseInt(tableData[i][3], 10, 64)
		valJ, _ := strconv.ParseInt(tableData[j][3], 10, 64)
		return valI > valJ
	})

	cleanUp := func(tbl [][]string) [][]string { // cleanUp func removes the unix timestamp col from the tabledata
		result := make([][]string, len(tbl))
		for i, inner := range tbl {
			n := len(inner)
			result[i] = inner[:n-1]
		}
		return result
	}

	tableData = cleanUp(tableData)

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	// selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	statusColors := map[string]lipgloss.Color{
		"Idea":        lipgloss.Color("#FF87D7"),
		"Indexed":     lipgloss.Color("#727272"),
		"Not Started": lipgloss.Color("#D7FF87"),
		"Ongoing":     lipgloss.Color("#00E2C7"),
		"Started":     lipgloss.Color("#00E2C7"),
		"Paused":      lipgloss.Color("#7D5AFC"),
		"Completed":   lipgloss.Color("#75FBAB"),
		"Aborted":     lipgloss.Color("#FF875F"),
		"Default":     lipgloss.Color("#929292"),
	}
	headers := []string{"Status", "Project Name", "Last Edited"}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(90).
		Rows(tableData...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			color, ok := statusColors[fmt.Sprint(tableData[row-1][0])]
			if ok {
				return baseStyle.Foreground(color)
			}
			color = statusColors["Default"]
			return baseStyle.Foreground(color)
		})
	fmt.Println(t)
	return nil
}
