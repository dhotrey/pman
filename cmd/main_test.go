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
