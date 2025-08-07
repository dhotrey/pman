package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

func Test_StatusCmd(t *testing.T) {
	t.Run("Test get status of a project with no alias", func(t *testing.T) {
		setup(t)
		t.Cleanup(func() {
			teardown(t)
		})

		// get current status
		currentStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)

		// execute
		r, w, _ := os.Pipe()
		os.Stdout = w
		rootCmd.SetArgs([]string{"status", projectName})
		err = rootCmd.Execute()
		require.NoError(t, err)
		w.Close()
		out, _ := io.ReadAll(r)

		// verify
		expected := fmt.Sprintf("Status of %s  : %s\n", projectName, utils.TitleCase(currentStatus))
		assert.Equal(t, expected, string(out))
	})

	t.Run("Test get status of a project using its alias", func(t *testing.T) {
		setup(t)
		t.Cleanup(func() {
			teardown(t)
		})

		// setup alias
		rootCmd.SetArgs([]string{"alias", projectName, aliasName})
		err := rootCmd.Execute()
		require.NoError(t, err)

		// get current status
		currentStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)

		// execute
		r, w, _ := os.Pipe()
		os.Stdout = w
		rootCmd.SetArgs([]string{"status", aliasName})
		err = rootCmd.Execute()
		require.NoError(t, err)
		w.Close()
		out, _ := io.ReadAll(r)

		// verify
		expected := fmt.Sprintf("Status of %s (%s) : %s\n", projectName, aliasName, utils.TitleCase(currentStatus))
		assert.Equal(t, expected, string(out))
	})

	t.Run("Test get status of a project that has an alias, using its project name", func(t *testing.T) {
		setup(t)
		t.Cleanup(func() {
			teardown(t)
		})

		// setup alias
		rootCmd.SetArgs([]string{"alias", projectName, aliasName})
		err := rootCmd.Execute()
		require.NoError(t, err)

		// get current status
		currentStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)

		// execute
		r, w, _ := os.Pipe()
		os.Stdout = w
		rootCmd.SetArgs([]string{"status", projectName})
		err = rootCmd.Execute()
		require.NoError(t, err)
		w.Close()
		out, _ := io.ReadAll(r)

		// verify
		expected := fmt.Sprintf("Status of %s  : %s\n", projectName, utils.TitleCase(currentStatus))
		assert.Equal(t, expected, string(out))
	})
}
