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
  isLeaf() (bool,HuffBase)
  getFreq() int
  getLeft() HuffBase
  getRight() HuffBase
  getChar() rune

}

type HuffLeaf struct{
  Freq int
  Char rune
}

func (h HuffLeaf) isLeaf() (bool,HuffBase){
  return true,h
}

func (h HuffLeaf) getLeft() HuffBase {
  return nil
}

func (h HuffLeaf) getRight() HuffBase{
  return nil
}

func (h HuffLeaf) getFreq() int{
  return h.Freq
}

func (h HuffLeaf) getChar() rune{
  return h.Char
}

type HuffInternal struct{
  Freq int
  Left HuffBase
  Right HuffBase
}

func (h HuffInternal) getLeft() HuffBase {
  return h.Left
}

func (h HuffInternal) getRight() HuffBase{
  return h.Right
}
func (h HuffInternal) getFreq() int{
  return h.Freq
}

func (h HuffInternal) isLeaf() (bool,HuffBase){
  return false,h
}

func (h HuffInternal) getChar() rune{
  return ' '
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

func dfs(node HuffBase,str string, codes map[rune]string){
  isLeaf , x := node.isLeaf()
  if isLeaf {
    codes[node.getChar()] = str
    return
  }

  dfs(x.getLeft(),fmt.Sprintf("%s0",str),codes)
  dfs(x.getRight(),fmt.Sprintf("%s1",str),codes)
  
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
    for _,v := range count{
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
    huffTree := consHuffTree(huffSlice)
    
    var huffCodes map[rune]string = make(map[rune]string)
    dfs(huffTree,"",huffCodes)
    for k,v := range huffCodes{
      fmt.Printf("%q - %q\n",k,v)
    }
    
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
