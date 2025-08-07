package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

const (
	dbname      = db.DBName
	projectName = "test-project"
	aliasName   = "tp"
	status      = "testing"
	testDir     = "test_dir"
)

func setup(t *testing.T) {
	err := os.MkdirAll(testDir+"/"+projectName, 0755)
	require.NoError(t, err)
	_, err = os.Create(testDir + "/" + projectName + "/README.md")
	require.NoError(t, err)

	rootCmd.SetArgs([]string{"init", testDir})
	err = rootCmd.Execute()
	require.NoError(t, err)
}

func teardown(t *testing.T) {
	err := db.DeleteDb(dbname)
	require.NoError(t, err)
	err = os.RemoveAll(testDir)
	require.NoError(t, err)
}

func Test_SetCmd(t *testing.T) {
	t.Run("Test set status of a project with no alias", func(t *testing.T) {
		setup(t)
		t.Cleanup(func() {
			teardown(t)
		})

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
		t.Cleanup(func() {
			teardown(t)
		})

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
		t.Cleanup(func() {
			teardown(t)
		})

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
