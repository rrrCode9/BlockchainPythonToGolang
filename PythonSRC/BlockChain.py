import os
from hashlib import sha256
from random import randint
from datetime import datetime
import json

class BlockChainGenerator():

    def __init__(self,fileName, newChain=True,difficultyLevel=2):
        self.fileName = fileName
        self.newChain = newChain
        self.difficultyLevel = difficultyLevel
        self.BlockChain = []
        self.HandleJson()
        self.AddBlock("GENESIS BLOCK") if self.BlockChain == [] else None

    def HandleJson(self):
        if not self.newChain:
            try:
                with open(self.fileName) as f:
                    self.BlockChain = json.load(f)
            except Exception as e:
                print(f"Error: {e}")
    
    #Create block
    def CreateBlock(self,data):
        block = {}
        block["data"] = data
        block["index"]= str(len(self.BlockChain))
        block["timeStamp"] = str(datetime.utcnow())
        block["previousHash"] = self.BlockChain[0]["currentHash"] if self.BlockChain != [] else "x"
        block["currentHash"], block["nounce"] = self.miner(
            block["data"]+block["index"]+block["timeStamp"]+block["previousHash"]
        )

        return block

    def miner(self,dataString):

        while True:
            nounce = str(randint(0,1E10))
            hash = sha256(str(dataString+nounce).encode()).hexdigest()
            if hash[:self.difficultyLevel]=="0"*self.difficultyLevel:
                return hash,nounce

    #add block to BlockChain
    def AddBlock(self,data):
        self.BlockChain = [self.CreateBlock(data)]+self.BlockChain 

        with open(self.fileName,"w") as f:
            json.dump(self.BlockChain,f)


#---------------------------------------------------
BlockChainFileName = "BlockChain.json"
Transaction = BlockChainGenerator(BlockChainFileName,newChain=False,difficultyLevel=3)
while True:
    Transaction.AddBlock(input("Enter Transaction ::: "))
    os.system(BlockChainFileName)
