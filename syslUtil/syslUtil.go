package syslutil

import (
	"os"
	"os/exec"
)

func Execute(args []string) {
	app := "sysl"
	currentDir, _ := os.Getwd()
	// err := os.Chdir(filepath.Join(currentDir, "tmp"))
	// check(err)
	cmd := exec.Command(app, args[0:]...)
	cmd.Dir = currentDir
	_, e := cmd.Output()
	check(e)
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
