package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Measurement struct {
	config   Config
	v4client *githubv4.Client
}

type PageInfo struct {
	HasNextPage bool
	EndCursor   githubv4.String
}

type Repository struct {
	Name          string
	NameWithOwner string
	Url           string
}

type RepositoriesQuery struct {
	Organization struct {
		Repositories struct {
			PageInfo PageInfo
			Nodes    []Repository
		} `graphql:"repositories(first:100, after: $repositoriesCursor)"`
	} `graphql:"organization(login: $owner)"`
}

var queryAllRepositories RepositoriesQuery

type DeploymentStatus struct {
	Id        string
	CreatedAt string
	Creator   struct {
		Login string
	}
	State       string
	Description string
	Environment string
	LogUrl      string
	Deployment  struct {
		Id         string
		DatabaseId int64
		CreatedAt  string
	}
}

type DeploymentType struct {
	State       string
	Environment string
	Description string
	Payload     string
	CreatedAt   string
	Creator     struct {
		Login string
	}
	LatestStatus DeploymentStatus
	Commit       struct {
		AuthoredDate           string
		Oid                    string
		Url                    string
		AssociatedPullRequests struct {
			PageInfo PageInfo
			Nodes    []struct {
				CreatedAt string
				Url       string
			}
		} `graphql:"associatedPullRequests(first:5)"`
	}
	Statuses struct {
		PageInfo PageInfo
		Nodes    []DeploymentStatus
	} `graphql:"statuses(last: 100)"`
}

var queryAllDeployments struct {
	Repository struct {
		Name        string
		Description string
		Url         string
		Deployments struct {
			PageInfo PageInfo
			Nodes    []DeploymentType
		} `graphql:"deployments(first: 100, orderBy:{field:CREATED_AT, direction:DESC})"`
	} `graphql:"repository(owner:$owner, name:$name)"`
}

func (m *Measurement) getRepositoriesPage(ctx context.Context, variables map[string]interface{}) (RepositoriesQuery, error) {
	err := m.v4client.Query(ctx, &queryAllRepositories, variables)
	if err != nil {
		fmt.Printf("%v\n", err)
		return RepositoriesQuery{}, err
	}
	return queryAllRepositories, err
}

func (m *Measurement) GetAllRepositories(ctx context.Context) ([]Repository, error) {
	var allRepos []Repository

	variables := map[string]interface{}{
		"owner":              githubv4.String(m.config.GithubOwner),
		"repositoriesCursor": (*githubv4.String)(nil), // nil for the first page
	}

	for {
		page, err := m.getRepositoriesPage(ctx, variables)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, page.Organization.Repositories.Nodes...)
		if !page.Organization.Repositories.PageInfo.HasNextPage {
			return allRepos, nil
		}

		variables["repositoriesCursor"] = githubv4.NewString(page.Organization.Repositories.PageInfo.EndCursor)
	}
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

	return queryAllDeployments, err
}

type withPreviewHeader struct {
	http.Header
	rt             http.RoundTripper
	previewHeaders []string
}

func (h withPreviewHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	currentAccept := req.Header.Get("Accept")
	// shadow-cat-preview: Draft pull requests
	req.Header.Set("Accept", fmt.Sprintf("%s,%s", currentAccept, strings.Join(h.previewHeaders, ",")))
	return h.rt.RoundTrip(req)
}

func roundTripperWithPreviewHeader(rt http.RoundTripper, previewHeaders []string) withPreviewHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	headers := make(http.Header)
	return withPreviewHeader{Header: headers, rt: rt, previewHeaders: previewHeaders}
}

func NewMeasurement(ctx context.Context, config Config) *Measurement {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	previewV4Headers := []string{
		"application/vnd.github.antiope-preview+json",
		"application/vnd.github.bane-preview+json",
		"application/vnd.github.flash-preview+json",
		"application/vnd.github.shadow-cat-preview+json",
	}
	tc.Transport = roundTripperWithPreviewHeader(tc.Transport, previewV4Headers)

	return &Measurement{
		config:   config,
		v4client: githubv4.NewClient(tc),
	}
}
