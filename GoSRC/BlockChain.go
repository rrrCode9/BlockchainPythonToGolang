package main

import (
	"os"
	"os/exec"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"bufio"
)

//Using struct for Blockcahin generator

type BlockchainGenerator struct{

	FileName string
	NewChain bool
	DifficultyLevel	int
	BlockChain []Block
}

type Block struct{
	Data,Index, TimeStamp, PreviousHash,CurrentHash,Nounce string
}

func (b *BlockchainGenerator) HandleJson(){

	if !b.NewChain{
		filebytes,err := ioutil.ReadFile(b.FileName)
		if err != nil{
			fmt.Println("##",err) //Dont exit
		}else{

			var existingData []Block
			json.Unmarshal([]byte(filebytes),&existingData)
			b.BlockChain = existingData
		}
	}
}

//Create Block
func (b *BlockchainGenerator) CreateBlock(data string) Block{

	block := Block{}
	block.Data = data
	block.Index = strconv.Itoa(len(b.BlockChain))
	block.TimeStamp = time.Now().UTC().Format(time.RFC3339Nano)
	block.PreviousHash = func()string{
		if len(b.BlockChain) != 0 {
			return b.BlockChain[0].CurrentHash
		}else{
			return "x"
		}
	}()
	
	block.CurrentHash,block.Nounce = b.miner(block.Data+ block.Index+ block.TimeStamp+ block.PreviousHash)

	return block
}

func (b BlockchainGenerator) miner(dataString string) (string,string){

	for true{

		nounce := strconv.Itoa(rand.Intn(1e10))
		hash := sha256.Sum256([]byte(dataString+nounce))
		hash_str := base64.StdEncoding.EncodeToString(hash[:])

		if hash_str[:b.DifficultyLevel] == strings.Repeat("0",b.DifficultyLevel){

			return hash_str,nounce
		}
	}

	return "-","-"

}

// Add block

func (b *BlockchainGenerator) AddBlock(data string){

	b.BlockChain = append([]Block{b.CreateBlock(data)},b.BlockChain...)

	jsonFile,_ := json.MarshalIndent(b.BlockChain,""," ")
	_ = ioutil.WriteFile(b.FileName, jsonFile, 0644)
}


func init(){

	rand.Seed(time.Now().UnixNano()) //Random number linked to time
}


func main(){

	Transaction := BlockchainGenerator{

		FileName : "BlockChain.json",
		NewChain : false,
		DifficultyLevel: 2,

	}

	Transaction.HandleJson()

	if len(Transaction.BlockChain) == 0{

		Transaction.AddBlock("GENESIS BLOCK")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for true{

		fmt.Printf("Enter Transaction ::: ")
		scanner.Scan()
		Transaction.AddBlock(scanner.Text())
		cmd := exec.Command("powershell","start","firefox",Transaction.FileName)
		cmd.Run()
	}


}