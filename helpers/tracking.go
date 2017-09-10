package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

// Get the docker container id.
// The following function returns the docker container id as string which is running the app right now.
func GetContainerId() string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "/bin/bash"
	cmdArgs := []string{"./scripts/get_container_id.sh"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
	}
	dockerContainerContainerID := string(cmdOut)
	return dockerContainerContainerID
}
