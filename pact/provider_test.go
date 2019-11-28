package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"github.mpi-internal.com/Yapo/events-router/pkg/infrastructure"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)

type PactConf struct {
	BrokerHost   string `env:"BROKER_HOST" envDefault:"http://3.229.36.112"`
	BrokerPort   string `env:"BROKER_PORT" envDefault:"80"`
	ProviderHost string `env:"PROVIDER_HOST" envDefault:"http://localhost"`
	ProviderPort string `env:"PROVIDER_PORT" envDefault:"8080"`
}

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider
func TestProvider(t *testing.T) {
	fmt.Printf("Pact directory: %+v", pactDir)
	var conf PactConf
	infrastructure.LoadFromEnv(&conf)
	var pact = &dsl.Pact{
		Consumer: "events-router",
		Provider: "profile-ms",
		Port:     6666,
	}
	files, err := IOReadDir(pactDir)
	if err != nil {
		fmt.Printf("Error in reading files. Error %+v", err)

	}
	for _, file := range files {
		// Verify the Provider with local Pact Files
		h := types.VerifyRequest{
			ProviderBaseURL:       conf.ProviderHost + ":" + conf.ProviderPort,
			PactURLs:              []string{pactDir + "/" + file},
			CustomProviderHeaders: []string{"Authorization: basic e5e5e5e5e5e5e5"},
		}
		_, err := pact.VerifyProvider(t, h)
		if err != nil {
			fmt.Printf("Error verifying the provider.Error %+v\n", err)
			return
		}
		pactPublisher := &dsl.Publisher{}
		err = pactPublisher.Publish(types.PublishRequest{
			PactURLs:        []string{"./pacts/events-router.json"},
			PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
			ConsumerVersion: "1.0.0",
			Tags:            []string{"events-router"},
		})
		if err != nil {
			fmt.Printf("Error with the Pact Broker server. Error %+v\n", err)
			return
		}
	}
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
