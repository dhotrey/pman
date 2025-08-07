package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestInitAndAdd_PreservesStatus(t *testing.T) {
	setup(t)
	defer teardown(t)

	// Set status to active
	rootCmd.SetArgs([]string{"set", projectName, "active"})
	err := rootCmd.Execute()
	require.NoError(t, err)

	// Run init again
	rootCmd.SetArgs([]string{"init", testDir})
	err = rootCmd.Execute()
	require.NoError(t, err)

	// Check status is still active
	status, err := db.GetRecord(dbname, projectName, "projects")
	require.NoError(t, err)
	require.Equal(t, "active", status, "pman init should not overwrite existing status")

	// Set status to inactive
	rootCmd.SetArgs([]string{"set", projectName, "inactive"})
	err = rootCmd.Execute()
	require.NoError(t, err)

	// Run add
	rootCmd.SetArgs([]string{"add", testDir})
	err = rootCmd.Execute()
	require.NoError(t, err)

	// Check status is still inactive
	status, err = db.GetRecord(dbname, projectName, "projects")
	require.NoError(t, err)
	require.Equal(t, "inactive", status, "pman add should not overwrite existing status")
}
