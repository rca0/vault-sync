package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type Vault struct {
	SourceAddr       string
	SourceToken      string
	DestinationAddr  string
	DestinationToken string
	authType         string
}

var (
	githubRef  string
	githubSHA  string
	httpClient = &http.Client{Timeout: 30 * time.Second}
	tokenPath  = "/.github-token"
)

func getToken(tokenPath string) (string, error) {
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", errors.Wrap(err, "get token")
	}
	return strings.TrimSpace(string(token)), nil
}

func login(addr, token, authType string) (*api.Logical, error) {
	config := &api.Config{
		Address:    addr,
		HttpClient: httpClient,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	switch authType {
	case "github":
		data, err := client.Logical().Write(
			"/auth/github/login",
			map[string]interface{}{"token": token},
		)
		if err != nil {
			return nil, err
		}
		client.SetToken(data.Auth.ClientToken)
	case "token":
		client.SetToken(token)
	}

	return client.Logical(), nil
}

func getPaths(client *api.Logical, currentPath string) ([]string, error) {
	var tmpValue string
	var results []string

	secret, err := client.List(currentPath)
	if err != nil {
		return []string{""}, err
	}
	if secret == nil {
		return []string{currentPath}, nil
	}

	for _, v := range secret.Data["keys"].([]interface{}) {
		tmpValue = v.(string)
		tmpValue = fmt.Sprintf("%s%s", currentPath, tmpValue)
		innerResults, err := getPaths(client, tmpValue)
		if err != nil {
			return results, err
		}
		results = append(results, innerResults...)
	}
	return results, nil
}

func copy(srcClient, dstClient *api.Logical, paths []string, src, dst string, wg *sync.WaitGroup) {
	for _, p := range paths {
		if strings.HasPrefix(p, "kv-v2") {
			p = strings.ReplaceAll(p, "metadata", "data")
		}
		secret, err := srcClient.Read(p)
		if err != nil {
			errors.Wrapf(err, "read path")
		}

		_, err = dstClient.Write(p, secret.Data)
		if err != nil {
			errors.Wrapf(err, "copy to destination")
		}

		fmt.Printf("~> from [%s] to [%s] - writting secrets on path: %s\n", src, dst, p)
	}
	wg.Done()
}

func parseToken(v Vault) Vault {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	if v.SourceToken == "" || v.DestinationToken == "" {
		token, err := getToken(home + tokenPath)
		if err != nil {
			log.Fatal(err)
		}
		v.SourceToken = token
		v.DestinationToken = token
		return v
	}
	return v
}

func run(v Vault) error {
	v = parseToken(v)
	var wg sync.WaitGroup

	srcClient, err := login(v.SourceAddr, v.SourceToken, v.authType)
	if err != nil {
		return err
	}

	dstClient, err := login(v.DestinationAddr, v.DestinationToken, v.authType)
	if err != nil {
		return err
	}

	kvPaths, err := getPaths(srcClient, "kv-v2/metadata/")
	if err != nil {
		return err
	}

	secretPaths, err := getPaths(srcClient, "secret/")
	if err != nil {
		return err
	}

	wg.Add(1)
	go copy(srcClient, dstClient, kvPaths, v.SourceAddr, v.DestinationAddr, &wg)

	wg.Add(1)
	go copy(srcClient, dstClient, secretPaths, v.SourceAddr, v.DestinationAddr, &wg)

	wg.Wait()
	return nil
}

func main() {
	var v Vault

	app := &cli.App{
		Name:    "vault-sync",
		Version: fmt.Sprintf("%s commit: %s", githubRef, githubSHA),
		Usage:   "synchronize vault data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "srcaddr",
				Usage:       "Source Vault Address",
				Required:    true,
				Destination: &v.SourceAddr,
			},
			&cli.StringFlag{
				Name:        "srctoken",
				Usage:       "Source Vault Token",
				Required:    false,
				Destination: &v.SourceToken,
			},
			&cli.StringFlag{
				Name:        "dstaddr",
				Usage:       "Destination Vault Address",
				Required:    true,
				Destination: &v.DestinationAddr,
			},
			&cli.StringFlag{
				Name:        "dsttoken",
				Usage:       "Destination Vault Token",
				Required:    false,
				Destination: &v.DestinationToken,
			},
			&cli.StringFlag{
				Name:        "method",
				Usage:       "Define auth method (github/token)",
				Required:    false,
				Value:       "github",
				Destination: &v.authType,
			},
		},
		Action: func(c *cli.Context) error {
			return run(v)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
