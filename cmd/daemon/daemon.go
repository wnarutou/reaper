package daemon

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/leslieleung/reaper/internal/config"
	"github.com/leslieleung/reaper/internal/discussion"
	"github.com/leslieleung/reaper/internal/issue"
	"github.com/leslieleung/reaper/internal/release"
	"github.com/leslieleung/reaper/internal/rip"
	"github.com/leslieleung/reaper/internal/typedef"
	"github.com/leslieleung/reaper/internal/ui"
	"github.com/leslieleung/reaper/internal/wiki"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "daemon runs as a daemon to monitor git repositories",
	Run:   runDaemon,
}

func runDaemon(cmd *cobra.Command, args []string) {
	storageMap := config.GetStorageMap()

	concurrencyNum := config.GetConcurrencyNum()

	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
		gocron.WithLimitConcurrentJobs(concurrencyNum, gocron.LimitModeWait),
	)
	if err != nil {
		ui.ErrorfExit("Error creating scheduler, %s", err)
	}

	for _, repo := range rip.GetRepositories("") {
		if repo.Cron == "" {
			continue
		}
		storages := make([]typedef.MultiStorage, 0)
		for _, storage := range repo.Storage {
			if s, ok := storageMap[storage]; !ok {
				continue
			} else {
				storages = append(storages, s)
			}
		}
		_, err := s.NewJob(
			gocron.CronJob(repo.Cron, false),
			gocron.NewTask(rip.Rip, repo, storages),
		)
		if err != nil {
			ui.Errorf("Error scheduling download codes of %s, %s", repo.Name, err)
		}
		if repo.DownloadReleases {
			_, err = s.NewJob(
				gocron.CronJob(repo.Cron, false),
				gocron.NewTask(release.DownloadAllAssets, repo, storages),
			)
			if err != nil {
				ui.Errorf("Error scheduling download releases of %s, %s", repo.Name, err)
			}
		}
		if repo.DownloadIssues {
			_, err = s.NewJob(
				gocron.CronJob(repo.Cron, false),
				gocron.NewTask(issue.Sync, repo, storages),
			)
			if err != nil {
				ui.Errorf("Error scheduling download issues of %s, %s", repo.Name, err)
			}
		}
		if repo.DownloadWiki {
			_, err = s.NewJob(
				gocron.CronJob(repo.Cron, false),
				gocron.NewTask(wiki.Sync, repo, storages),
			)
			if err != nil {
				ui.Errorf("Error scheduling download wiki of %s, %s", repo.Name, err)
			}
		}
		if repo.DownloadDiscussion {
			_, err = s.NewJob(
				gocron.CronJob(repo.Cron, false),
				gocron.NewTask(discussion.Sync, repo, storages),
			)
			if err != nil {
				ui.Errorf("Error scheduling download discussion of %s, %s", repo.Name, err)
			}
		}
		ui.Printf("Scheduled %s, cron: %s", repo.Name, repo.Cron)
	}
	ui.Printf("Starting daemon")
	s.Start()
}
