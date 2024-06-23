package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/fr0stylo/funcgo/pkg/runtime"
)

func containerInit() {
	initPath, ok := os.LookupEnv("FUNC_INIT")
	if !ok {
		log.Fatal("FUNC_INIT is not defined")
	}
	hostname, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		log.Fatal("HOSTNAME is not defined")
	}
	ip, ok := os.LookupEnv("IP")
	if !ok {
		log.Fatal("HOSTNAME is not defined")
	}
	fmt.Fprintf(os.Stdout, "Container inside\n")
	fmt.Fprintf(os.Stdout, "%s\n", os.Args[2])

	must("hostname: ", syscall.Sethostname([]byte(hostname)))

	cg()

	must("chdir: ", syscall.Chroot(os.Args[2]))
	must("chdir: ", syscall.Chdir("/"))
	must("proc: ", syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("/proc", 0)
	// ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)

	must("Setup net: ", runtime.SetupNet(ip))

	ctx := context.Background()

	cmd := exec.CommandContext(ctx, initPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must("cmd: ", cmd.Run())
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "ourContainer"), 0755)
	os.WriteFile(filepath.Join(pids, "ourContainer/pids.max"), []byte("10"), 0700)
	//up here we limit the number of child processes to 10

	os.WriteFile(filepath.Join(pids, "ourContainer/notify_on_release"), []byte("1"), 0700)

	os.WriteFile(filepath.Join(pids, "ourContainer/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	// up here we write container PIDs to cgroup.procs
}
