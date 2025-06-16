package utils

import (
	"fmt"
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
		conn, errCon = ssh.Dial("tcp", server+":22", config)
		if errCon != nil {
			err = errCon
			continue
		} else {
			for i := 0; i < 10; i++ {
				if conn != nil {
					session, errSesh = conn.NewSession()
					if errSesh != nil {
						err = errSesh
						continue
					} else {
						return conn, session, nil
					}
				} else {
					break
				}
			}
		}
	}
	return nil, nil, err
}

func Close(session *ssh.Session, conn *ssh.Client) {
	session.Close()
	conn.Close()
}
