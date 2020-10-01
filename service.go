package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type ServiceInfo interface{}

type StandaloneProcessServiceInfo struct {
	Path       string
	LaunchArgs []string
}

const (
	ServiceStopped int32 = iota
	ServiceRunning
	ServiceStarting
	ServiceExited
)

type ClientService interface {
	Start(onExit chan interface{}, id interface{})
	Stop()
	Exit()
	Status() int32
	ID() string
}

type StandaloneProcessService struct {
	info         StandaloneProcessServiceInfo
	id           string
	cmd          *exec.Cmd
	status       int32
	additionArgs []string
}

func NewService(info ServiceInfo, id string, additionArgs []string) ClientService {
	if i, ok := info.(StandaloneProcessServiceInfo); ok {
		return &StandaloneProcessService{
			info:         i,
			status:       ServiceStopped,
			id:           id,
			additionArgs: additionArgs,
		}
	}
	return nil
}
func (s *StandaloneProcessService) Start(onExit chan interface{}, id interface{}) {
	s.status = ServiceStarting
	params := url.Values{}
	params.Add("x-service-id", s.ID())
	params.Add("x-mother-id", MotherID)
	params.Add("token", ManagerToken)
	actualArgs := make([]string, 0, len(s.info.LaunchArgs)+2+len(s.additionArgs))
	actualArgs = append(actualArgs, s.info.LaunchArgs...)
	actualArgs = append(actualArgs, "applyto", GetUBotAddr("ws", "/api/manager")+"?"+params.Encode())
	actualArgs = append(actualArgs, s.additionArgs...)
	cmd := exec.Command(s.info.Path, actualArgs...)
	s.cmd = cmd
	logFile, _ := os.OpenFile(LogFilePath(s.id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	setDeathsig(cmd.SysProcAttr)
	go func() {
		err := cmd.Start()
		if err != nil {
			log.Printf("failed to start %s: %v", s.id, err)
			_, _ = logFile.WriteString(fmt.Sprintf(
				"%s failed to start %s: %v\n",
				time.Now().UTC().Format(time.RFC3339),
				s.id,
				err))
			_ = logFile.Close()
			s.cmd = nil
			onExit <- id
			return
		}
		s.status = ServiceRunning
		_ = cmd.Wait()
		_ = logFile.Close()
		if s.status == ServiceRunning {
			s.status = ServiceExited
		}
		s.cmd = nil
		onExit <- id
	}()
}
func (s *StandaloneProcessService) Stop() {
	s.status = ServiceStopped
	ExitCmd(s.cmd)
}
func (s *StandaloneProcessService) Exit() {
	ExitCmd(s.cmd)
}
func (s *StandaloneProcessService) Status() int32 {
	return s.status
}
func (s *StandaloneProcessService) ID() string {
	return s.id
}
