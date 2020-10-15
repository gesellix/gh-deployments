package pkg

type Config struct {
	GithubToken  string
	GithubOwner  string
	GithubRepo   string
	Ref          string
	Payload      string
	Description  string
	Environment  string
	DeploymentId int64
	State        string
	LogUrl       string
}
