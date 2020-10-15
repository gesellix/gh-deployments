package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
	owner := "gesellix"
	repository := "gh-deployments"
	ref := "e63c1a01c8093c28f8ba4886dcab128464815889"

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repo, _, err := client.Repositories.Get(ctx, owner, repository)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("deployments url: %s\n", *repo.DeploymentsURL)

	//deployment, _, err := client.Repositories.GetDeployment(ctx, owner, "gh-deployments", 279652485)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	//os.Exit(1)
	//}
	//fmt.Printf("deployment: %+v\n", deployment)

	createDeployment := &github.DeploymentRequest{
		Ref:                  github.String(ref),
		Task:                 github.String("deploy"),
		TransientEnvironment: github.Bool(true),
	}
	createdDeployment, _, err := client.Repositories.CreateDeployment(ctx, owner, repository, createDeployment)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("deployment created: %d\n", *createdDeployment.ID)

	createDeploymentStatus := &github.DeploymentStatusRequest{
		State:       github.String("pending"),
		Description: github.String("foo"),
	}
	deploymentStatus, _, err := client.Repositories.CreateDeploymentStatus(ctx, owner, repository, *createdDeployment.ID, createDeploymentStatus)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("deployment status: %d\n", *deploymentStatus.ID)
}
