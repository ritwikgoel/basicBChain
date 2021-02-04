package main

import(
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"fmt"
	"time"
)
var gcounter=0;

type Block struct{
	Index int
	Timestamp string
	Data string
	Hash string//SHA256 identifier representing data record
	PrevHash string 
}
var BlockChain []Block//slice of a block
/*
	PrevHash of the next block should be same as the Hash of the current block to ensure the integrity
	of the BC.
	Data is hashed to maintain integrity and to also save space 
*/

func HashIt(b Block) string{
	// rec string
	rec:=string(b.Index)+b.Timestamp+b.Data+b.PrevHash
	hh:=sha256.New()
	hh.Write([]byte(rec))
	hashed:=hh.Sum(nil)
	return hex.EncodeToString(hashed)

}

func createBlock(oldBlock Block, data string)Block{
	var newBlock Block
	t:=time.Now()
	newBlock.Timestamp=t.String()
	newBlock.Index=oldBlock.Index+1
	newBlock.Data=data
	newBlock.PrevHash=oldBlock.Hash
	newBlock.Hash=HashIt(newBlock)
	gcounter++
	BlockChain=append(BlockChain,newBlock)
	return newBlock
}


func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if HashIt(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func view(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BChain server Init\n")
	//var localCounter=0
	for i,v := range BlockChain {
		fmt.Fprintf(w,"%d",v.Index)
		i=i
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BChain server write \n")
	
	if(isBlockValid(createBlock(BlockChain[gcounter],"Data Here"),BlockChain[gcounter])) {
		fmt.Fprintf(w,"Written")
	} else {
		fmt.Fprintf(w,"error")
	}

}

func main() {
	var initBlock Block
	initBlock.Index=0
	t:=time.Now()
	initBlock.Timestamp=t.String()
	initBlock.Data="INIT"
	initBlock.PrevHash=""
	initBlock.Hash=HashIt(initBlock)
	BlockChain=append(BlockChain,initBlock)
	http.HandleFunc("/view/",view)
	http.HandleFunc("/add/",add)
	http.ListenAndServe(":8080",nil)
}