package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func execc(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	return c.Run()
}

func create(initPath string) {
	t := time.Now()
	d, err := os.MkdirTemp("", "container-")
	if err != nil {
		log.Fatal("tmpdir ", err)
	}
	defer os.RemoveAll(d)
	log.Print(d)

	if err := execc("cp", "-r", "./fs", d); err != nil {
		log.Fatal("cp: ", err)
	}

	if err := execc("cp", initPath, d+"/fs/"); err != nil {
		log.Fatal("cp wrapper: ", err)
	}

	cmd := exec.CommandContext(context.Background(), "/proc/self/exe", "container", d+"/fs")
	cmd.Env = []string{"FUNC_INIT=" + initPath}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWNET,
		Unshareflags: syscall.CLONE_NEWNS,
		// Foreground:   true,
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	must("cmd: ", cmd.Run())

	log.Print(time.Since(t))
}
