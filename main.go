package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/go-github/v60/github"
	"github.com/jessevdk/go-flags"
	"github.com/rodaine/table"
	"golang.org/x/oauth2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	fmt.Printf("github-projects %s, commit %s, built at %s by %s\n\n", version, commit, date, builtBy)

	var opts struct {
		GithubUser  string `short:"u" long:"user" default:"attachmentgenie" required:"true" name:"github user"`
		GithubToken string `short:"t" long:"token" env:"GITHUB_TOKEN" required:"true" name:"github auth token"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(2)
	}

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

	// repository counters
	forks :=0
	archives :=0
	private :=0
	//template_repos := 0
	// issue counter
	openIssues := 0
	for _, repository := range allRepos {
		if !*repository.Archived {
			if *repository.Fork {
				forks++
			}
			if *repository.Private {
				private++
			}
			// @todo doesnt seemed to be returned as a property
			//if *repository.IsTemplate {
			//	template_repos++
			//}

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
		} else {
			archives++
		}
	}

	tbl.Print()
	fmt.Println("Repos:", len(allRepos))
	fmt.Println("Forked repos:", forks)
	fmt.Println("Archived repos:", archives)
	fmt.Println("Private repos:", private)
	//fmt.Println("Templates:", template_repos)
	fmt.Println("Open Issues:", openIssues)
}
