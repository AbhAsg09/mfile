package utils

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	errCon  error
	errSesh error
)

func CreateConnection(server string, username string, password string) (conn *ssh.Client, session *ssh.Session, err error) {
	err = nil
	session = nil
	conn = nil
	fmt.Printf("Connecting to %s...\n", server)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	for i := 0; i < 10; i++ {
		conn, errCon := ssh.Dial("tcp", server+":22", config)
		if errCon != nil {
			fmt.Printf("SSH dial failed (%s): %v\n", server, errCon)
			continue
		} else {
			for i := 0; i < 10; i++ {
				session, errSesh := conn.NewSession()
				if errSesh != nil {
					time.Sleep(5 * time.Second)
					fmt.Printf("SSH session failed (%s): %v\n", server, errSesh)
					continue
				} else {
					return conn, session, nil
				}
			}
			return nil, nil, errSesh
		}
	}
	return nil, nil, errCon
}

func Close(session *ssh.Session, conn *ssh.Client) {
	session.Close()
	conn.Close()
}

func SudoAccess(session *ssh.Session, password string) bool {
	for i := 0; i < 10; i++ {
		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		cmd := "sudo -S -l"
		stdin, _ := session.StdinPipe()
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, password+"\n")
		}()

		err := session.Run(cmd)
		if err == nil {
			return true
		} else {
			continue
		}
	}
	return false
}
