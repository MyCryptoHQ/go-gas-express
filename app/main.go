package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/ethclient"
	root "github.com/mycryptohq/go-gas-express/pkg"
	"github.com/mycryptohq/go-gas-express/pkg/client"
	"github.com/mycryptohq/go-gas-express/pkg/helpers"
	"github.com/mycryptohq/go-gas-express/pkg/s3"
)

var (
	SAFELOW  = .35
	STANDARD = .60
	FAST     = .90
	FASTEST  = 1.0
)

func handleRequest() {
	blockEstimateNumber, err := strconv.ParseInt(root.Config.BlockEstimateCount, 10, 64)
	if err != nil {
		fmt.Printf("%d of type %T\n", blockEstimateNumber, blockEstimateNumber)
	}
	fmt.Printf("Using block interval of %d\n", blockEstimateNumber)
	ethClient := client.MakeETHClient()
	endBlockNum, err := client.GetBlockNumber(*ethClient)
	if endBlockNum == int64(0) || err != nil {
		log.Fatalf("Couldn't fetch Block Number")
	}

	startBlockNum := endBlockNum - blockEstimateNumber
	blockCacheFileEndpoint := "blockCacheFile.json"
	gasExpressFileEndpoint := "gasExpress.json"

	// Download output file from s3
	var cachedRun root.CachedItem
	if err := s3.Download(root.Config.Bucket, root.Config.Region, blockCacheFileEndpoint); err != nil {
		log.Println("\nCouldn't open output file. Assume this is first run.")
	} else {
		blockCacheFile, err := os.Open("/tmp/" + blockCacheFileEndpoint)
		if err != nil {
			log.Printf("Couldn't open %s", blockCacheFile.Name())
		}
		defer blockCacheFile.Close()
		byteOutputFileValue, _ := ioutil.ReadAll(blockCacheFile)
		_ = json.Unmarshal(byteOutputFileValue, &cachedRun)
	}

	var newCachedBlocks []root.CachedBlockItem
	var fetchedBlocks []root.CachedBlockItem
	var newCache root.CachedItem
	if cachedRun.LastBlock >= startBlockNum {
		fmt.Printf("Num of blocks to lookup: %d\n", endBlockNum-cachedRun.LastBlock)
		// If relevant blocks are already in the cached blocks, take those cachedBlocks and add to newCachedBlocks
		for _, item := range cachedRun.Blocks {
			if item.BlockNum.Int64() >= startBlockNum {
				newCachedBlocks = append(newCachedBlocks, item)
			}
		}
		//Take the other blocks still needed and look them up.
		fetchedBlocks = fetchAndProcess(ethClient, cachedRun.LastBlock+1, endBlockNum)
		newCache = root.CachedItem{
			LastBlock: endBlockNum,
			Blocks:    append(newCachedBlocks, fetchedBlocks...),
		}
	} else {
		fmt.Printf("Num of blocks to lookup: %d\n", endBlockNum-startBlockNum)
		fetchedBlocks = fetchAndProcess(ethClient, startBlockNum, endBlockNum)
		newCache = root.CachedItem{
			LastBlock: endBlockNum,
			Blocks:    fetchedBlocks,
		}
	}

	gasExpressObject := createGasExpress(newCache.Blocks, blockEstimateNumber, endBlockNum)

	// Write new cache to s3
	newCacheFile, _ := json.Marshal(newCache)
	if err := s3.Upload(root.Config.Bucket, root.Config.Region, blockCacheFileEndpoint, bytes.NewReader(newCacheFile)); err != nil {
		log.Println("Error uploading to s3", err)
	}

	// Write new gasExpress to s3
	gasExpressFile, _ := json.Marshal(gasExpressObject)
	if err := s3.Upload(root.Config.Bucket, root.Config.Region, gasExpressFileEndpoint, bytes.NewReader(gasExpressFile)); err != nil {
		log.Println("Error uploading to s3", err)
	}
}
func main() {
	//handleRequest()
	lambda.Start(handleRequest)
}

func fetchAndProcess(nodeClient *ethclient.Client, startBlock int64, endBlock int64) []root.CachedBlockItem {
	var fetchedBlocks []root.CachedBlockItem
	for i := startBlock; i <= endBlock; i++ {
		block, err := client.GetBlockByNum(*nodeClient, i)
		if err != nil {
			fmt.Println("Block couldn't be pulled from the eth node: ", i)
		} else {
			minGas := client.FetchBlockMinGas(block.Transactions())

			fetchedBlocks = append(fetchedBlocks, root.CachedBlockItem{
				BlockNum: *big.NewInt(i),
				MinGas:   *big.NewInt(minGas),
			})
		}

	}
	return fetchedBlocks
}

type byMinGas []root.CachedBlockItem

func (a byMinGas) Len() int           { return len(a) }
func (a byMinGas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMinGas) Less(i, j int) bool { return a[i].MinGas.Int64() < a[j].MinGas.Int64() }

func roundDown(val float64) int {
	if val < 0 {
		return int(val - 1.0)
	}
	return int(val)
}

func createGasExpress(blocks []root.CachedBlockItem, blockEstimateNumber int64, lastBlock int64) root.GasExpress {
	sort.Sort(byMinGas(blocks))
	return root.GasExpress{
		SafeLow:  math.Ceil(helpers.ConvertFromBase(blocks[roundDown(float64(blockEstimateNumber)*SAFELOW)].MinGas, 9)),
		Standard: math.Ceil(helpers.ConvertFromBase(blocks[roundDown(float64(blockEstimateNumber)*STANDARD)].MinGas, 9)),
		Fast:     math.Ceil(helpers.ConvertFromBase(blocks[roundDown(float64(blockEstimateNumber)*FAST)].MinGas, 9)),
		Fastest:  math.Ceil(helpers.ConvertFromBase(blocks[roundDown(float64(blockEstimateNumber)*FASTEST)].MinGas, 9)),
		BlockNum: lastBlock,
	}
}
