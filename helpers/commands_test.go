package helpers_test

import (
	"github.com/cpapidas/pegasus/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContainerID(t *testing.T) {
	var id string

	// Should return commandSuccess
	helpers.GetContainerIDScriptPath = "../tests/fixtures/commands/get_container_id.sh"
	id = helpers.GetContainerID()
	assert.Equal(t, "commandSuccess", id)

	// Should return the default error message
	helpers.GetContainerIDScriptPath = "invalid command"
	id = helpers.GetContainerID()
	assert.Equal(t, "Container ID not found", id)

	// Should return the default error message for undefined container id
	helpers.GetContainerIDScriptPath = `echo ""`
	id = helpers.GetContainerID()
	assert.Equal(t, "Container ID not found", id)
}