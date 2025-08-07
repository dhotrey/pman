package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

func Test_SetCmd(t *testing.T) {
	t.Run("Test set status of a project with no alias", func(t *testing.T) {
		setup(t)

		// execute
		rootCmd.SetArgs([]string{"set", projectName, status})
		err := rootCmd.Execute()
		require.NoError(t, err)

		// verify
		actualStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)
		assert.Equal(t, status, actualStatus)
	})

	t.Run("Test set status of a project using its alias", func(t *testing.T) {
		setup(t)

		// setup alias
		rootCmd.SetArgs([]string{"alias", projectName, aliasName})
		err := rootCmd.Execute()
		require.NoError(t, err)

		// execute
		rootCmd.SetArgs([]string{"set", aliasName, status})
		err = rootCmd.Execute()
		require.NoError(t, err)

		// verify
		actualStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)
		assert.Equal(t, status, actualStatus)
	})

	t.Run("Test set status of a project that has an alias, using its project name", func(t *testing.T) {
		setup(t)

		// setup alias
		rootCmd.SetArgs([]string{"alias", projectName, aliasName})
		err := rootCmd.Execute()
		require.NoError(t, err)

		// execute
		rootCmd.SetArgs([]string{"set", projectName, status})
		err = rootCmd.Execute()
		require.NoError(t, err)

		// verify
		actualStatus, err := db.GetRecord(dbname, projectName, c.StatusBucket)
		require.NoError(t, err)
		assert.Equal(t, status, actualStatus)
	})
}
