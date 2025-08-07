package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

var (
	ErrFlagNotImplemented = errors.New("flag not implemented yet")
	ErrBadUsageSetCmd     = errors.New("bad usage of set command")
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the status of a project",
	Long: `Set the status of a project to a specified status
    Usage:
    pman set <project_name> <status>

    Common statuses: Indexed (default), Idea, Started, Paused, Completed, Aborted, Ongoing, Not Started
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		interactiveFlag, _ := cmd.Flags().GetBool("i") // TODO: Implement this
		if interactiveFlag {
			cmd.SilenceUsage = true
			return ErrFlagNotImplemented
		}
		if len(args) != 2 {
			fmt.Println("Please provide a directory name")
			return ErrBadUsageSetCmd
		}
		pname := args[0]
		status := args[1]

		// check if user has supplied an alias instead of actual project name
		proj, err := db.GetRecord(db.DBName, pname, c.ProjectAliasBucket)
		if err == nil {
			_, err2 := db.GetRecord(db.DBName, proj, c.StatusBucket)
			if err2 == nil {
				pname = proj
			}
		}
		err = db.UpdateRec(db.DBName, pname, status, c.StatusBucket)
		if err != nil {
			fmt.Println("Error updating record : ", err)
			return err
		}

		lastEdit := make(map[string]string)
		lastEdit["lastWrite"] = fmt.Sprint(time.Now().Format("02 Jan 06 15:04"))
		err = db.WriteToDB(db.DBName, lastEdit, c.ConfigBucket)
		if err != nil {
			log.Print(err)
			return err
		}

		fmt.Printf("Project %s set to status %s\n", pname, status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("i", false, "Set the status of projects interactively")
}
