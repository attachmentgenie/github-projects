package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/jessevdk/go-flags"
	"github.com/rodaine/table"
	"golang.org/x/oauth2"
)

func main() {

	var opts struct {
		GithubUser  string `short:"u" long:"user" default:"attachmentgenie" required:"true" name:"github user"`
		GithubToken string `short:"t" long:"token" required:"true" name:"github auth token"`
	}
	flags.Parse(&opts)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: opts.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		Affiliation: "owner",
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				log.Println("hit rate limit")
			} else {
				fmt.Println("error : ", err)
			}
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Repo", "Type", "Number", "Description", "Created")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	openIssues := 0
	for _, repository := range allRepos {

		issues, _, _ := client.Issues.ListByRepo(ctx, opts.GithubUser, *repository.Name, nil)

		if len(issues) > 0 {
			for _, issue := range issues {

				var issueType = "Issue"
				if issue.IsPullRequest() {
					issueType = "PR"
				}
				tbl.AddRow(*repository.Name, issueType, *issue.Number, *issue.Title, *issue.CreatedAt)
				openIssues++
			}
		}
	}

	tbl.Print()
	fmt.Println("Public Repos:", len(allRepos))
	fmt.Println("Open Issues:", openIssues)
}
