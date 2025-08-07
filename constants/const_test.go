package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, "projects", StatusBucket)
	assert.Equal(t, "projectPaths", ProjectPaths)
	assert.Equal(t, "projectAliases", ProjectAliasBucket)
	assert.Equal(t, "config", ConfigBucket)
	assert.Equal(t, "1.1.1", Version)
	assert.Equal(t, "lastUpdated", LastUpdatedBucket)
}
