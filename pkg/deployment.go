package pkg

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type Deployment struct {
	config   Config
	v3client *github.Client
}

func (d *Deployment) CreateDeployment(ctx context.Context) error {
	createDeployment := &github.DeploymentRequest{
		Ref:                  github.String(d.config.Ref),
		Task:                 github.String("deploy"),
		Environment:          github.String(d.config.Environment),
		TransientEnvironment: github.Bool(false),
	}
	createdDeployment, _, err := d.v3client.Repositories.CreateDeployment(ctx, d.config.GithubOwner, d.config.GithubRepo, createDeployment)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Printf("env=%s\n", d.config.Environment)
	fmt.Printf("deployment_id=%d\n", *createdDeployment.ID)

	return nil
}

func (d *Deployment) UpdateDeploymentStatus(ctx context.Context) error {
	createDeploymentStatus := &github.DeploymentStatusRequest{
		State:       github.String(d.config.State),
		Description: github.String(d.config.Description),
		Environment: github.String(d.config.Environment),
	}
	deploymentStatus, _, err := d.v3client.Repositories.CreateDeploymentStatus(ctx, d.config.GithubOwner, d.config.GithubRepo, d.config.DeploymentId, createDeploymentStatus)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Printf("env=%s\n", d.config.Environment)
	fmt.Printf("deployment_status_id=%d\n", *deploymentStatus.ID)

	return nil
}

func NewDeployment(ctx context.Context, config Config) *Deployment {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Deployment{
		config:   config,
		v3client: github.NewClient(tc),
	}
}
