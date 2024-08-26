package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

type HuffBase interface{
  isLeaf() bool
  getFreq() int

}

type HuffLeaf struct{
  Freq int
  Char rune
}

func (h HuffLeaf) isLeaf() bool{
  return true
}


func (h HuffLeaf) getFreq() int{
  return h.Freq
}

type HuffInternal struct{
  Freq int
  Left HuffBase
  Right HuffBase
}

func (h HuffInternal) getFreq() int{
  return h.Freq
}

func (h HuffInternal) isLeaf() bool{
  return false
}

func bucketSort(mp map[rune]int, length int) [][]rune {
  var bucket = make([][]rune,length+1)
  for k,v := range mp{
     bucket[v] = append(bucket[v],rune(k))
  }
  return bucket
} 

func reverse[T any](arr *[]T){
  l := 0
  r := len(*arr)-1
  for l < r{
    temp := (*arr)[l]
    (*arr)[l] = (*arr)[r]
    (*arr)[r] = temp
    l++
    r--
  }
}

func consHuffTree(huffSlice []HuffBase) HuffInternal{

  for len(huffSlice) > 1{
    first := huffSlice[0]
    second := huffSlice[1]
    huffSlice = huffSlice[2:]
    newNode := HuffInternal{Freq: first.getFreq() + second.getFreq() , Left: first , Right: second}
    i := 0
    for i < len(huffSlice) && huffSlice[i].getFreq() < newNode.getFreq(){
      i++
    }

    temp := append([]HuffBase{},huffSlice[:i]...)
    temp = append(temp,newNode)
    huffSlice = slices.Concat(temp,huffSlice[i:])

  }
  
  if len(huffSlice) == 1{
    return huffSlice[0].(HuffInternal)
  }else{
    return HuffInternal{}
  }
}

func huffmanCoding(fd io.Reader) error{

    rd := bufio.NewReader(fd)
    buff := make([]byte,1024*4)
    count := make(map[rune]int)
    for{
      _, err := rd.Read(buff)
      if err != nil{
        if err == io.EOF{
          break
        }
        return err
      }
      for _,x := range buff{
        count[rune(x)] += 1
      }

    }
    max := 0
    for k,v := range count{
      fmt.Printf("%q - %d\n",k,v)
      if v > max{
        max = v
      }
    }
    buckets := bucketSort(count,max)
    
    huffSlice := []HuffBase{}
    for i,bucket := range buckets{
      for _,x := range bucket{
         huffSlice = append(huffSlice, HuffLeaf{Freq: i,Char: rune(x)})
      }
    }
    consHuffTree(huffSlice)
    
    return nil
}

func main(){
  args := os.Args
  l := len(args)
  switch {
  case l < 3:
    fmt.Fprintln(os.Stderr,"please give filename...")
    os.Exit(1)
  default:
    if args[1] == "--file"{
      fileName := args[2]

      fd , err := os.Open(fileName)
      if err != nil{
        log.Fatalln("file does not exits")
      }
      err = huffmanCoding(fd)
      if err != nil{
        log.Fatal("Failed...")
      }
    }
  }
}
