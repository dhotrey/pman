package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

func TestIndexDir(t *testing.T) {
	// Create a temporary directory for our test setup
	tmpDir, err := os.MkdirTemp("", "pman-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a directory structure for testing
	// Project 1: valid project with a README.md
	project1Dir := filepath.Join(tmpDir, "project1")
	err = os.Mkdir(project1Dir, 0755)
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(project1Dir, "README.md"))
	assert.NoError(t, err)

	// Project 2: valid project with a .git directory
	project2Dir := filepath.Join(tmpDir, "project2")
	err = os.Mkdir(project2Dir, 0755)
	assert.NoError(t, err)
	err = os.Mkdir(filepath.Join(project2Dir, ".git"), 0755)
	assert.NoError(t, err)

	// Project 3: not a project directory, no README.md or .git
	project3Dir := filepath.Join(tmpDir, "project3")
	err = os.Mkdir(project3Dir, 0755)
	assert.NoError(t, err)

	// Nested project, should not be indexed
	nestedProjectDir := filepath.Join(project1Dir, "nested-project")
	err = os.Mkdir(nestedProjectDir, 0755)
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(nestedProjectDir, "README.md"))
	assert.NoError(t, err)

	// Call the function we are testing
	projDirs, err := indexDir(tmpDir)
	assert.NoError(t, err)

	// Assertions
	assert.Len(t, projDirs, 2, "should have indexed 2 projects")

	absProject1, _ := filepath.Abs(project1Dir)
	absProject2, _ := filepath.Abs(project2Dir)

	_, p1Indexed := projDirs[absProject1]
	_, p2Indexed := projDirs[absProject2]

	assert.True(t, p1Indexed, "project1 should be indexed")
	assert.True(t, p2Indexed, "project2 should be indexed")

	absProject3, _ := filepath.Abs(project3Dir)
	_, p3Indexed := projDirs[absProject3]
	assert.False(t, p3Indexed, "project3 should not be indexed")

	absNested, _ := filepath.Abs(nestedProjectDir)
	_, nestedIndexed := projDirs[absNested]
	assert.False(t, nestedIndexed, "nested project should not be indexed")
}

func TestInitDirs_PreservesStatus(t *testing.T) {
	db.DeleteDb(db.DBName)
	defer db.DeleteDb(db.DBName)
	// Create a temporary directory for our test setup
	tmpDir, err := os.MkdirTemp("", "pman-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a directory structure for testing
	project1Dir := filepath.Join(tmpDir, "project1")
	err = os.Mkdir(project1Dir, 0755)
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(project1Dir, "README.md"))
	assert.NoError(t, err)

	// --- First run of InitDirs ---
	err = InitDirs([]string{tmpDir})
	assert.NoError(t, err)

	// --- Check status is "indexed" ---
	status, err := db.GetRecord(db.DBName, "project1", c.StatusBucket)
	assert.NoError(t, err)
	assert.Equal(t, "indexed", status)

	// --- Manually change status ---
	err = db.UpdateRec(db.DBName, "project1", "active", c.StatusBucket)
	assert.NoError(t, err)

	// --- Second run of InitDirs ---
	err = InitDirs([]string{tmpDir})
	assert.NoError(t, err)

	// --- Check status is still "active" ---
	status, err = db.GetRecord(db.DBName, "project1", c.StatusBucket)
	assert.NoError(t, err)
	assert.Equal(t, "active", status, "InitDirs should not overwrite existing status")
}
