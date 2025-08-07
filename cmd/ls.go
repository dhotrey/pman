package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/ui"
	"github.com/theredditbandit/pman/pkg/utils"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all indexed projects along with their status",
	Long: `List all indexed projects along with their status
    Usage : pman ls
    `,
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, _ []string) error {
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.Suffix = " Compiling..."
		s.Start()
		defer s.Stop()
		filterFlag, _ := cmd.Flags().GetString("f")
		refreshLastEditTime, _ := cmd.Flags().GetBool("r")
		data, err := db.GetAllRecords(db.DBName, c.StatusBucket)
		if err != nil {
			return err
		}
		if filterFlag != "" {
			fmt.Println("Filtering by status : ", filterFlag)
			data = utils.FilterByStatuses(data, strings.Split(filterFlag, ","))
		}
		return ui.RenderTable(db.DBName, data, refreshLastEditTime)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls --f <status1[,status2]>")
	lsCmd.Flags().Bool("r", false, "Refresh Last Edited time: pman ls --r")
}
