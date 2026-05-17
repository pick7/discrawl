package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/openclaw/crawlkit/releasecheck"
	"github.com/openclaw/discrawl/internal/config"
)

const discrawlUpgradeHint = "brew upgrade openclaw/tap/discrawl"

func discrawlReleaseCheckOptions(force bool) releasecheck.Options {
	cfg := config.Default()
	return releasecheck.Options{
		AppName:        "discrawl",
		Owner:          "openclaw",
		Repo:           "discrawl",
		CurrentVersion: version,
		CacheDir:       cfg.CacheDir,
		Force:          force,
	}
}

func (r *runtime) maybeNotifyRelease(args []string) {
	_, _ = releasecheck.Notify(r.ctx, releasecheck.NotifyOptions{
		Options:     discrawlReleaseCheckOptions(false),
		Stderr:      r.stderr,
		InstallHint: discrawlUpgradeHint,
		Args:        args,
		JSONOutput:  r.json,
		IsTerminal:  releasecheck.StderrIsTerminal(),
	})
}

func (r *runtime) runCheckUpdate(args []string) error {
	fs := flag.NewFlagSet("check-update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	jsonOut := fs.Bool("json", false, "")
	force := fs.Bool("force", false, "")
	if err := fs.Parse(args); err != nil {
		return usageErr(err)
	}
	if fs.NArg() != 0 {
		return usageErr(errors.New("check-update takes flags only"))
	}
	if *jsonOut {
		r.json = true
	}
	result, err := releasecheck.Check(r.ctx, discrawlReleaseCheckOptions(*force))
	if err != nil && !errors.Is(err, releasecheck.ErrSkipped) {
		return err
	}
	if r.json {
		return r.print(result)
	}
	_, err = fmt.Fprint(r.stdout, releasecheck.StatusText("discrawl", discrawlUpgradeHint, result))
	return err
}
