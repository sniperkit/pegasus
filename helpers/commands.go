package helpers

import (
	"bitbucket.org/code_horse/pegasus/blunder"
	"errors"
	"net/url"
	"os/exec"
	"strings"
)

// GetContainerIDScriptPath is the shell script file path which returns the docker container id
var GetContainerIDScriptPath = "./scripts/get_container_id.sh"

// GetContainerID returns the docker running container id as string. It needs the ./scripts/get_container_id.sh
// in order to get get container id. If something goes wrong then it will return "Container ID not found".
func GetContainerID() string {
	id, err := commandRunner(GetContainerIDScriptPath, "Container ID not found")
	if err != nil {
		return err.Error()
	}
	return id
}

// commandRunner will run a command. If the script runs successfully then it will return the exported data as string
// if something doesn't go well then it will return an error with the given error message.
func commandRunner(command string, errMsg string) (string, error) {

	var (
		cmdOut []byte
		err    error
	)

	cmdName := "/bin/bash"
	cmdArgs := []string{command}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		blunder.Set("Error on helpers/commands.go function commandRunner", err).Handle()
		return "", errors.New(errMsg)
	}

	// Remove the last %0A from string
	results := strings.Replace(url.QueryEscape(string(cmdOut)), "%0A", "", -1)
	return results, nil
}
