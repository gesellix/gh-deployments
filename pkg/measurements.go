package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Measurement struct {
	config   Config
	v4client *githubv4.Client
}

var queryAllRepositories struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []struct {
				Name          string
				NameWithOwner string
				Url           string
			}
		} `graphql:"repositories(first:100)"`
	} `graphql:"organization(login: $owner)"`
}

var queryAllDeployments struct {
	Repository struct {
		Name        string
		Description string
		Url         string
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
					State       string
					CreatedAt   string
					Description string
					LogUrl      string
				}
				Payload     string
				Description string
				Commit      struct {
					AuthoredDate           string
					AssociatedPullRequests struct {
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
						Nodes []struct {
							CreatedAt string
						}
					} `graphql:"associatedPullRequests(first:1)"`
				}
				Statuses struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
					Nodes []struct {
						State       string
						CreatedAt   string
						Description string
						LogUrl      string
					}
				} `graphql:"statuses(last: 100)"`
			}
		} `graphql:"deployments(first: 100, orderBy:{field:CREATED_AT, direction:DESC})"`
	} `graphql:"repository(owner:$owner, name:$name)"`
}

func (m *Measurement) GetAllRepositories(ctx context.Context) (interface{}, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(m.config.GithubOwner),
	}
	err := m.v4client.Query(ctx, &queryAllRepositories, variables)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	// TODO Pagination

	reposJSON, err := json.Marshal(queryAllRepositories)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", reposJSON)

	return queryAllRepositories, err
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

	deploymentsJSON, err := json.Marshal(queryAllDeployments)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", deploymentsJSON)

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
