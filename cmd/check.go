package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/cobra"
)

var (
	serverListPath string
	userListPath   string
)

type Job struct {
	Server   string
	Username string
	Password string
}

type Result struct {
	Server     string
	Username   string
	SSHAccess  bool
	SudoAccess bool
	Err        error
}

type UserCred struct {
	Username string
	Password string
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check SSH and sudo access",
	Run: func(cmd *cobra.Command, args []string) {
		if serverListPath == "" {
			fmt.Println("--serverList is required")
			return
		}

		servers, err := readLines(serverListPath)
		if err != nil {
			fmt.Printf("Error reading server list: %v\n", err)
			return
		}

		if userListPath == "" {
			fmt.Println("--userList is required")
			return
		}

		userLines, err := readLines(userListPath)
		if err != nil {
			fmt.Printf("Error reading user list: %v\n", err)
			return
		}

		var userCreds []UserCred
		for _, line := range userLines {
			parts := strings.SplitN(line, ",", 2)
			if len(parts) != 2 {
				fmt.Printf("Invalid user line: %s\n", line)
				continue
			}
			userCreds = append(userCreds, UserCred{
				Username: strings.TrimSpace(parts[0]),
				Password: strings.TrimSpace(parts[1]),
			})
		}

		var allResults []Result
		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, user := range userCreds {
			wg.Add(1)
			go func(user UserCred) {
				defer wg.Done()
				userResults := checkUserAccess(user, servers)
				mu.Lock()
				allResults = append(allResults, userResults...)
				mu.Unlock()
			}(user)
		}

		wg.Wait()
		writeCSV("check_results_k.csv", allResults)
		fmt.Printf("Check completed. Results saved to file")
	},
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func checkUserAccess(user UserCred, servers []string) []Result {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			sshOk, sudoOk := checkSSHAndSudo(server, user.Username, user.Password)
			resultsChan <- Result{
				Server:     server,
				Username:   user.Username,
				SSHAccess:  sshOk,
				SudoAccess: sudoOk,
			}
		}(strings.TrimSpace(server))
	}

	wg.Wait()
	close(resultsChan)

	var results []Result
	for r := range resultsChan {
		results = append(results, r)
	}
	return results
}

func checkSSHAndSudo(server, username, password string) (bool, bool) {
	sshSuccess := false
	sudoSuccess := false

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
		conn, err := ssh.Dial("tcp", server+":22", config)
		if err != nil {
			fmt.Printf("SSH dial failed (%s): %v\n", server, err)
			continue
		} else {
			for i := 0; i < 10; i++ {
				session, err := conn.NewSession()
				if err != nil {
					time.Sleep(5 * time.Second)
					fmt.Printf("SSH session failed (%s): %v\n", server, err)
					continue
				} else {
					sshSuccess = true
					var stdoutBuf, stderrBuf bytes.Buffer
					session.Stdout = &stdoutBuf
					session.Stderr = &stderrBuf

					cmd := fmt.Sprintf("echo %s | sudo -S -l", password)
					err = session.Run(cmd)
					if err != nil {
						fmt.Println("Error running sudo -l:", err)
					}

					out := stdoutBuf.String()
					errOut := stderrBuf.String()

					if bytes.Contains([]byte(out), []byte("may run the following commands")) {
						sudoSuccess = true
					}

					session.Close()
					conn.Close()
					if errOut == "" && err == nil {
						break
					}
				}
			}
			if sudoSuccess {
				break
			}
		}
	}

	return sshSuccess, sudoSuccess
}

func writeCSV(filename string, results []Result) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"user", "server", "access", "sudo access"})
	for _, r := range results {
		writer.Write([]string{
			r.Username,
			r.Server,
			fmt.Sprintf("%v", r.SSHAccess),
			fmt.Sprintf("%v", r.SudoAccess),
		})
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVar(&serverListPath, "serverList", "", "Path to file with list of servers")
	checkCmd.Flags().StringVar(&userListPath, "userList", "", "Path to file with list of users (username,password)")
}
