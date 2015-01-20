package lrps

import (
	"flag"
	"os"
	"time"

	"github.com/pivotal-cf-experimental/veritas/common"
	"github.com/pivotal-cf-experimental/veritas/config_finder"
	"github.com/pivotal-cf-experimental/veritas/lrps/set_domain"
	"github.com/pivotal-cf-experimental/veritas/say"
)

func SetDomainCommand() common.Command {
	var (
		etcdClusterFlag string
	)

	flagSet := flag.NewFlagSet("set-domain", flag.ExitOnError)
	flagSet.StringVar(&etcdClusterFlag, "etcdCluster", "", "comma-separated etcd cluster urls")

	return common.Command{
		Name:        "set-domain",
		Description: "domain ttl - updates the domain in the BBS (ttl is a duration)",
		FlagSet:     flagSet,
		Run: func(args []string) {
			etcdCluster, err := config_finder.FindETCDCluster(etcdClusterFlag)
			common.ExitIfError("Could not find etcd cluster", err)
			if len(args) != 2 {
				say.Fprintln(os.Stderr, 0, say.Red("Expected domain and ttl"))
				os.Exit(1)
			}
			ttl, err := time.ParseDuration(args[1])
			common.ExitIfError("Failed to parse TTL", err)

			err = set_domain.SetDomain(etcdCluster, args[0], ttl)
			common.ExitIfError("Failed to submit lrp", err)
		},
	}
}
