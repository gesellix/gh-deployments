package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Measurement struct {
	config   Config
	v4client *githubv4.Client
}

var queryAllDeployments struct {
	Repository struct {
		Description string
		Deployments struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []struct {
				State        string
				Environment  string
				CreatedAt    string
				LatestStatus struct {
					CreatedAt string
				}
				Commit struct {
					AuthoredDate           string
					AssociatedPullRequests struct {
						Nodes []struct {
							CreatedAt string
						}
					} `graphql:"associatedPullRequests(first:1)"`
				}
			}
		} `graphql:"deployments(first: 100, orderBy:{field:CREATED_AT, direction:DESC})"`
	} `graphql:"repository(owner:$owner,name:$name)"`
}

func (m *Measurement) GetAllDeployments(ctx context.Context) (interface{}, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(m.config.GithubOwner),
		"name":  githubv4.String(m.config.GithubRepo),
	}
	err := m.v4client.Query(ctx, &queryAllDeployments, variables)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	log.Printf("%+v\n", &queryAllDeployments)
	return queryAllDeployments, err
}

func NewMeasurement(ctx context.Context, config Config) *Measurement {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Measurement{
		config:   config,
		v4client: githubv4.NewClient(tc),
	}
}
