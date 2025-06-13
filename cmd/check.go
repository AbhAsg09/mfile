package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mfile/utils"
	"os"
	"strings"
	"sync"

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

	conn, session, err := utils.CreateConnection(server, username, password)
	if err != nil {
		fmt.Printf("Error connecting to server: %s \nError: %v", server, err)
	} else {
		sshSuccess = true
		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		cmd := "sudo -S -l"
		stdin, _ := session.StdinPipe()
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, password+"\n")
		}()

		err = session.Run(cmd)
		if err == nil {
			sudoSuccess = true
		}
	}
	utils.Close(session, conn)
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
