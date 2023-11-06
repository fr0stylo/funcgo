package runtime

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type Worker struct {
	api      *WorkerApi
	name     string
	busy     bool
	m        sync.Mutex
	initPath string
	tmpDir   string
	lastExec time.Time
	pid      int
	p        *os.Process
	hostname string
	cancel   context.CancelFunc
	ip       *netlink.Addr
}

type WorkerOpts struct {
	InitPath    string
	FilesToCopy []Files
}

type Files struct {
	From string
	To   string
}

func FileList(f ...Files) []Files {
	return f
}

func NewWorker(name string, ip *netlink.Addr, opts *WorkerOpts) *Worker {
	d := prepareFilesystem(opts.FilesToCopy)
	return &Worker{
		api:      NewWorkerApi(ip.IP.String()),
		name:     name,
		busy:     false,
		ip:       ip,
		initPath: opts.InitPath,
		m:        sync.Mutex{},
		tmpDir:   d,
		hostname: name,
		lastExec: time.Now(),
	}
}

var (
	defaultEnv = []string{
		"HOME=/root",
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"TERM=xterm",
	}
)

func (r *Worker) Start(c context.Context) error {
	ctx, cancel := context.WithCancel(c)
	r.cancel = cancel

	cmd := exec.CommandContext(ctx, "/proc/self/exe", "container", r.tmpDir+"/fs")
	cmd.Env = append(defaultEnv, []string{
		"FUNC_INIT=" + r.initPath,
		"HOSTNAME=" + r.hostname,
		"IP=" + r.ip.String(),
	}...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWUSER |
			unix.CLONE_NEWNS,
		Unshareflags: unix.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Geteuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getegid(),
				Size:        1,
			},
		},
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	putIface(cmd.Process.Pid)

	go cmd.Wait()

	raw_connect(r.ip.IP.String(), "9999")
	r.pid = cmd.Process.Pid
	r.p = cmd.Process
	return nil
}

func (w *Worker) Stop() {
	w.cancel()
	defaultIPManager.Release(w.ip)
}

func prepareFilesystem(fs []Files) string {
	d, err := os.MkdirTemp("", "container-")
	if err != nil {
		log.Fatal("tmpdir ", err)
	}
	log.Print(d)

	if err := execc("cp", "-r", "./fs", d); err != nil {
		log.Fatal("cp: ", err)
	}

	for _, v := range fs {
		if err := execc("cp", v.From, fmt.Sprintf("%s/fs%s", d, v.To)); err != nil {
			log.Fatal("cp wrapper: ", err)
		}
	}

	return d
}

func (r *Worker) SinceLastExecution() time.Duration {
	log.Printf("[%s]: %s", r.name, time.Since(r.lastExec))
	return time.Since(r.lastExec)
}

func (r *Worker) Cleanup() {
	defer os.RemoveAll(r.tmpDir)
}

func (r *Worker) setNotBusy() {
	r.busy = false
	r.m.Unlock()
}

func (r *Worker) setBusy() {
	r.busy = true
	r.m.Lock()
}

func (r *Worker) Execute(obj any) ([]byte, error) {
	r.setBusy()
	defer r.setNotBusy()
	log.Printf("[%s] exec: started", r.name)
	defer log.Printf("[%s] exec: end", r.name)
	log.Print()

	res, err := r.api.Execute(obj)
	if err != nil {
		return nil, err
	}
	r.lastExec = time.Now()

	return res, nil
}

func (r *Worker) IsBusy() bool {
	return r.busy
}

func execc(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	return c.Run()
}

func raw_connect(host string, port string) {
	for {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
		}
		if conn != nil {
			defer conn.Close()
			fmt.Println("ready", net.JoinHostPort(host, port))
			return
		}
	}
}
