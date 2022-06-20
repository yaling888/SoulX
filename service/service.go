//go:build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows"
)

type Config struct {
	Name, DisplayName, Description string

	Exec string

	Stderr, Stdout string
}

var logger service.Logger

type program struct {
	exit    chan struct{}
	service service.Service

	*Config

	cmd *exec.Cmd
}

func (p *program) Start(s service.Service) error {
	fullExec, err := exec.LookPath(p.Exec)
	if err != nil {
		return fmt.Errorf("failed to find executable %q: %v", p.Exec, err)
	}

	p.cmd = exec.Command(fullExec, "-d", homeDir)
	p.cmd.Dir = ExecutableDir()

	go p.run()
	return nil
}

func (p *program) run() {
	_ = logger.Info("Starting ", p.DisplayName)
	defer func() {
		if service.Interactive() {
			_ = p.Stop(p.service)
		} else {
			_ = p.service.Stop()
		}
	}()

	if p.Stderr != "" {
		f, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			_ = logger.Warningf("Failed to open std err %q: %v", p.Stderr, err)
			return
		}
		defer f.Close()
		p.cmd.Stderr = f
	}
	if p.Stdout != "" {
		f, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			_ = logger.Warningf("Failed to open std out %q: %v", p.Stdout, err)
			return
		}
		defer f.Close()
		p.cmd.Stdout = f
	}

	err := p.cmd.Run()
	if err != nil {
		_ = logger.Warningf("Error running: %v", err)
	}

	return
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	_ = logger.Info("Stopping ", p.DisplayName)
	if p.cmd.Process != nil {
		err := terminateProc(p.cmd.Process)
		if err != nil {
			_ = p.cmd.Process.Kill()
		}
	}
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the Soul system service.")
	flag.Parse()

	config := &Config{
		Name:        "Soul",
		DisplayName: "Soul Service",
		Description: "Control the Soul system service.",
		Exec:        CoreExecPath(),
		Stdout:      LogOutFile(),
		Stderr:      LogErrFile(),
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{
		exit: make(chan struct{}),

		Config: config,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	defer close(errs)

	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err, ok := <-errs
			if !ok {
				break
			}
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}

func terminateProc(process *os.Process) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	pid := process.Pid

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}
	return nil
}
