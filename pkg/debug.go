package pkg

import (
	"context"
	"fmt"
)

func (d *Deployment) Debug(ctx context.Context) error {
	repo, _, err := d.v3client.Repositories.Get(ctx, d.config.GithubOwner, d.config.GithubRepo)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Printf("deployments url: %s\n", *repo.DeploymentsURL)

	//deployment, _, err := v3client.Repositories.GetDeployment(ctx, owner, "gh-deployments", 279652485)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	return err
	//}
	//fmt.Printf("deployment: %+v\n", deployment)

	return nil
}
