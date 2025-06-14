package cmd

var (
	serverListPath string
	userListPath   string
	username       string
	addPath        string
	dropPath       string
	ans            string
	password       string
	server         string
	cfgFile        string
	checkSudo      bool
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
