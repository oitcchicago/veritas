package store

import (
	"flag"
	"io"
	"time"

	"github.com/cloudfoundry-incubator/bbs"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/store/fetch_store"
	"github.com/pivotal-cf-experimental/veritas/store/print_store"
)

func DumpStoreCommand() common.Command {
	var (
		bbsConfig config_finder.BBSConfig
		tasks     bool
		lrps      bool
		rate      time.Duration
		verbose   bool
	)

	flagSet := flag.NewFlagSet("dump-store", flag.ExitOnError)
	bbsConfig.PopulateFlags(flagSet)
	flagSet.BoolVar(&tasks, "tasks", true, "print tasks")
	flagSet.BoolVar(&lrps, "lrps", true, "print lrps")
	flagSet.BoolVar(&verbose, "v", false, "be verbose")
	flagSet.DurationVar(&rate, "rate", time.Duration(0), "rate at which to poll the store")

	return common.Command{
		Name:        "dump-store",
		Description: "- Fetch and print contents of the BBS",
		FlagSet:     flagSet,
		Run: func(args []string) {
			bbsClient, err := config_finder.NewBBS(bbsConfig)
			common.ExitIfError("Could not construct BBS", err)

			if rate == 0 {
				err = dump(bbsClient, verbose, tasks, lrps, false)
				common.ExitIfError("Failed to dump", err)
				return
			}

			ticker := time.NewTicker(rate)
			for {
				<-ticker.C
				err = dump(bbsClient, verbose, tasks, lrps, true)
				if err != nil {
					say.Println(0, say.Red("Failed to dump: %s", err.Error()))
				}
			}
		},
	}
}

func dump(bbsClient bbs.Client, verbose bool, tasks bool, lrps bool, clear bool) error {
	reader, writer := io.Pipe()

	errs := make(chan error)
	go func() {
		err := fetch_store.Fetch(bbsClient, writer)
		errs <- err
	}()
	go func() {
		err := print_store.PrintStore(verbose, tasks, lrps, clear, reader)
		errs <- err
	}()

	err1 := <-errs
	if err1 != nil {
		return err1
	}
	err2 := <-errs
	if err2 != nil {
		return err2
	}
	return nil
}
