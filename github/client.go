package gh

import (
	"context"
	"fmt"
	//"log"
	"os"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func GithubClient(ctx context.Context) *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

func Repos() error {

	ctx := context.Background()

	client := GithubClient(ctx)

	opt := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 20}}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, "leapfrogtechnology", opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	//if err != nil {
	//return err
	//}

	for _, repo := range allRepos {
		op := &github.PullRequestListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
		name := repo.GetName()
		var allPrs []*github.PullRequest
		for {
			prs, resp, err := client.PullRequests.List(ctx, "leapfrogtechnology", name, op)

			if err != nil {
				return err
			}

			var prIn2020 []*github.PullRequest

			for _, pr := range prs {

				year2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				createdAt := pr.CreatedAt
				isIn2020 := createdAt.After(year2020)
				if !isIn2020 {
					break
				}

				prIn2020 = append(prIn2020, pr)
			}

			//allPrs = append(allPrs, prs...)
			allPrs = append(allPrs, prIn2020...)

			if resp.NextPage == 0 {
				break
			}
			op.Page = resp.NextPage
		}

		fmt.Println(name, ",", len(allPrs))
	}
	return nil
}

func Prs() error {

	ctx := context.Background()

	client := GithubClient(ctx)
	op := &github.PullRequestListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
	name := "jump"
	var allPrs []*github.PullRequest
	for {
		prs, resp, err := client.PullRequests.List(ctx, "leapfrogtechnology", name, op)

		if err != nil {
			return err
		}
		fmt.Println("NEt", resp.NextPage)
		var prIn2020 []*github.PullRequest

		for _, pr := range prs {

			year2020 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			createdAt := pr.CreatedAt
			isIn2020 := createdAt.After(year2020)
			if !isIn2020 {
				fmt.Println(createdAt, pr)
				break
			}

			prIn2020 = append(prIn2020, pr)
		}
		allPrs = append(allPrs, prIn2020...)
		if resp.NextPage == 0 {
			break
		}
		op.Page = resp.NextPage
	}

	fmt.Println(name, len(allPrs))

	return nil
}
