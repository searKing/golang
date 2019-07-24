package exec

import "os/exec"

func KillProcByName(pname string) {
	params := []string{
		"taskkill",
		"/F",
		"/IM",
		pname,
		"/T",
	}
	exec.Command(params[0], params[1:]...).CombinedOutput()
}
