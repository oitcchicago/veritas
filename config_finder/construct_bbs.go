package config_finder

import (
	"errors"
	"flag"
	"net/url"
	"os"

	"github.com/cloudfoundry-incubator/bbs"
)

func NewBBS(c BBSConfig) (bbs.Client, error) {
	c.PopulateFromEnv()
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	if c.IsSecure() {
		return bbs.NewSecureSkipVerifyClient(c.URL, c.CertFile, c.KeyFile, 0, 0)
	} else {
		return bbs.NewClient(c.URL), nil
	}
}

type BBSConfig struct {
	URL      string
	CertFile string
	KeyFile  string
}

func (c *BBSConfig) PopulateFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.URL, "bbsEndpoint", "", "BBS url")
	flagSet.StringVar(&c.CertFile, "bbsCertFile", "", "path to BBS TLS cert file")
	flagSet.StringVar(&c.KeyFile, "bbsKeyFile", "", "path to BBS TLS key file")
}

func (c *BBSConfig) IsSecure() bool {
	u, err := url.Parse(c.URL)
	if err != nil {
		panic("crap")
	}
	return u.Scheme == "https"
}

func (c *BBSConfig) PopulateFromEnv() {
	if c.URL == "" {
		c.URL = os.Getenv("BBS_ENDPOINT")
	}
	if c.CertFile == "" {
		c.CertFile = os.Getenv("BBS_CERT_FILE")
	}
	if c.KeyFile == "" {
		c.KeyFile = os.Getenv("BBS_KEY_FILE")
	}
}

func (c *BBSConfig) Validate() error {
	if c.URL == "" {
		return errors.New("You must either specify --bbsEndpoint or set BBS_ENDPOINT")
	}
	if c.IsSecure() {
		if c.CertFile == "" {
			return errors.New("You must either specify --bbsCertFile or set BBS_CERT_FILE")
		}
		if c.KeyFile == "" {
			return errors.New("You must either specify --bbsKeyFile or set BBS_KEY_FILE")
		}
	}

	return nil
}
