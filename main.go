package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"gh-deployments/pkg"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var config pkg.Config

var appFlags []cli.Flag

func init() {
	appFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "github.token",
			Usage:   "Access Token for the GitHub API",
			EnvVars: []string{"GITHUB_TOKEN"},
			Hidden:  false,
			//Required:    true,
			Destination: &config.GithubToken,
		},
		&cli.StringFlag{
			Name:    "owner",
			Usage:   "Owner of the GitHub repository",
			EnvVars: []string{"GITHUB_OWNER"},
			Hidden:  false,
			//Required:    true,
			Destination: &config.GithubOwner,
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   "GitHub repository name",
			EnvVars: []string{"GITHUB_REPO"},
			Hidden:  false,
			//Required:    true,
			Destination: &config.GithubRepo,
		},
	}
}

func main() {
	ctx := context.Background()

	var defaultAction = func(c *cli.Context) error {
		deployment := pkg.NewDeployment(ctx, config)
		err := deployment.Debug(ctx)
		return err
	}

	app := cli.NewApp()
	app.Name = "GitHub Deployments Helper"
	app.Usage = ""
	app.Description = "Adopts common use cases for GitHub's Deployments API"
	app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	app.Flags = appFlags
	app.Before = beforeApp()
	app.Action = defaultAction
	app.Commands = []*cli.Command{
		{
			Name: "create",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ref",
					Usage:       "The ref to deploy. This can be a branch, tag, or SHA.",
					Hidden:      false,
					Required:    true,
					Destination: &config.Ref,
				},
				&cli.StringFlag{
					Name:        "payload",
					Usage:       "JSON payload with extra information about the deployment.",
					Hidden:      false,
					Destination: &config.Payload,
				},
				&cli.StringFlag{
					Name:        "description",
					Usage:       "Short description of the deployment.",
					Hidden:      false,
					Destination: &config.Description,
				},
				&cli.StringFlag{
					Name:        "environment",
					Usage:       "Name for the target deployment environment (e.g., production, staging, qa).",
					Hidden:      false,
					Value:       "production",
					Destination: &config.Environment,
				},
			},
			Action: func(c *cli.Context) error {
				deployment := pkg.NewDeployment(ctx, config)
				err := deployment.CreateDeployment(ctx)
				return err
			},
		},
		{
			Name: "status",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "description",
					Usage:       "Short description of the status.",
					Hidden:      false,
					Destination: &config.Description,
				},
				&cli.Int64Flag{
					Name:        "deployment-id",
					Usage:       "deployment_id parameter",
					Hidden:      false,
					Destination: &config.DeploymentId,
				},
				&cli.StringFlag{
					Name:        "state",
					Usage:       "The state of the status.",
					Hidden:      false,
					Required:    true,
					Destination: &config.State,
				},
				&cli.StringFlag{
					Name:        "log_url",
					Usage:       "The full URL of the deployment's output.",
					Hidden:      false,
					Destination: &config.LogUrl,
				},
				&cli.StringFlag{
					Name:        "environment",
					Usage:       "Name for the target deployment environment (e.g., production, staging, qa).",
					Hidden:      false,
					Value:       "production",
					Destination: &config.Environment,
				},
			},
			Action: func(c *cli.Context) error {
				deployment := pkg.NewDeployment(ctx, config)
				err := deployment.UpdateDeploymentStatus(ctx)
				return err
			},
		},
		{
			Name:  "measurements",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) error {
				deployment := pkg.NewMeasurement(ctx, config)
				_, err := deployment.GetAllDeployments(ctx)
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func beforeApp() cli.BeforeFunc {
	return func(context *cli.Context) error {
		//if context.String("github.token") == "" {
		//	return errors.New("GITHUB_TOKEN required")
		//}
		//if context.String("owner") == "" {
		//	return errors.New("GITHUB_OWNER required")
		//}
		//if context.String("repo") == "" {
		//	return errors.New("GITHUB_REPO required")
		//}
		return nil
	}
}
