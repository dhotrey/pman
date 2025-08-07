package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/ui"
)

// Mocking the db.GetAllRecords function
func mockGetAllRecords(dbname, bucketName string) (map[string]string, error) {
	if dbname == "error" {
		return nil, errors.New("database error")
	}
	return map[string]string{"project1": "ongoing", "project2": "completed"}, nil
}

// Mocking the ui.RenderTable function
func mockRenderTable(data map[string]string, refreshLastEditedTime bool) error {
	if data == nil {
		return errors.New("render error")
	}
	return nil
}

func TestLsCmd(t *testing.T) {
	// Replace the original functions with mocks
	originalGetAllRecords := db.GetAllRecords
	originalRenderTable := ui.RenderTable
	db.GetAllRecords = mockGetAllRecords
	ui.RenderTable = mockRenderTable
	defer func() {
		db.GetAllRecords = originalGetAllRecords
		ui.RenderTable = originalRenderTable
	}()

	// Create a new lsCmd instance
	cmd := &cobra.Command{Use: "ls"}
	cmd.RunE = lsCmd.RunE

	// Execute the command
	err := cmd.RunE(cmd, []string{})

	// Assertions
	assert.NoError(t, err, "ls command should not return an error")

	// Test with a database error
	db.GetAllRecords = func(dbname, bucketName string) (map[string]string, error) {
		return nil, errors.New("database error")
	}
	err = cmd.RunE(cmd, []string{})
	assert.Error(t, err, "ls command should return an error when the database fails")
	db.GetAllRecords = mockGetAllRecords // Reset for other tests

	// Test with a filter
	cmd.Flags().String("f", "", "")
	cmd.Flags().Set("f", "ongoing")
	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err, "ls command with filter should not return an error")

	// Test with refresh flag
	cmd.Flags().Bool("r", false, "")
	cmd.Flags().Set("r", "true")
	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err, "ls command with refresh flag should not return an error")
}
