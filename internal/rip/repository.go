package rip

import (
	"path"

	"github.com/leslieleung/reaper/internal/config"
	"github.com/leslieleung/reaper/internal/scm/github"
	"github.com/leslieleung/reaper/internal/typedef"
	"github.com/leslieleung/reaper/internal/ui"
)

func GetRepositories(name string) []typedef.Repository {
	repositories := make([]typedef.Repository, 0)
	if name != "" {
		// find repo in config
		for _, repository := range config.GetIns().Repository {
			if repository.Name == name {
				repositories = addRepo(repository, repositories)
			}
		}
		return repositories
	}
	for _, repository := range config.GetIns().Repository {
		repositories = addRepo(repository, repositories)
	}
	return repositories
}

func addRepo(repo typedef.Repository, ret []typedef.Repository) []typedef.Repository {
	switch repo.GetType() {
	case typedef.TypeRepo:
		ret = append(ret, repo)
	case typedef.TypeUser, typedef.TypeOrg:
		// get repos
		client, err := github.New()
		if err != nil {
			ui.Errorf("Error creating github client, %s", err)
			return ret
		}
		repos, err := client.GetRepos(repo.OrgName, repo.Type)
		if err != nil {
			ui.Errorf("Error getting user repos, %s", err)
			return ret
		}
		for _, r := range repos {
			ret = append(ret, typedef.Repository{
				Name:               path.Base(r),
				URL:                r,
				Cron:               repo.Cron,
				Storage:            repo.Storage,
				UseCache:           repo.UseCache,
				Type:               typedef.TypeRepo,
				AllBranches:        repo.AllBranches,
				Depth:              repo.Depth,
				DownloadReleases:   repo.DownloadReleases,
				DownloadIssues:     repo.DownloadIssues,
				DownloadWiki:       repo.DownloadWiki,
				DownloadDiscussion: repo.DownloadDiscussion,
			})
		}
	default:
		ui.Errorf("Invalid repository type %s", repo.Type)
	}
	return ret
}
