package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

var (
	ErrBadUsageStatusCmd = errors.New("bad usage of status command")
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a project",
	Long:  `Query the database for the status of a project.`,
	RunE: func(_ *cobra.Command, args []string) error {
		var alias string
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return ErrBadUsageStatusCmd
		}
		projName := args[0]
		// check if user has supplied an alias instead of actual project name
		proj, err := db.GetRecord(db.DBName, projName, c.ProjectAliasBucket)
		if err == nil {
			_, err2 := db.GetRecord(db.DBName, proj, c.StatusBucket)
			if err2 == nil {
				alias = projName
				projName = proj
			}
		}
		status, err := db.GetRecord(db.DBName, projName, c.StatusBucket)
		if err != nil {
			fmt.Printf("%s is not a valid project name : Err -> %s", projName, err)
			return err
		}
		if alias != "" {
			fmt.Printf("Status of %s (%s) : %s\n", projName, alias, utils.TitleCase(status))
		} else {
			fmt.Printf("Status of %s  : %s\n", projName, utils.TitleCase(status))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
