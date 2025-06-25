package cmd

var (
	serverListPath string
	userListPath   string
	username       string
	password       string
	addPath        string
	dropPath       string
	ans            string
	cfgFile        string
	fileName       string
	filePath       string
	server         string
	search         string
	date           string
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
