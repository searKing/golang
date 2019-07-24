package exec

import "os/exec"

func KillProcByName(pname string) {
	params := []string{
		"killall",
		"-9",
		pname,
	}
	exec.Command(params[0], params[1:]...).CombinedOutput()
}
