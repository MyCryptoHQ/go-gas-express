package client

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	root "github.com/mycryptohq/go-gas-express/pkg"
)

func MakeETHClient() *ethclient.Client {
	configEndpoint := root.Config.NodeEndpoint
	fmt.Println(configEndpoint)
	client, err := ethclient.Dial(configEndpoint)
	if err != nil {
		log.Fatalf("Could not connect to eth client")
	}
	return client
}
