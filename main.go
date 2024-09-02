package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

type HuffBase interface {
	isLeaf() (bool, HuffBase)
	getFreq() int
	getLeft() HuffBase
	getRight() HuffBase
	getChar() rune
}

type HuffLeaf struct {
	Freq int
	Char rune
}

func (h HuffLeaf) isLeaf() (bool, HuffBase) {
	return true, h
}

func (h HuffLeaf) getLeft() HuffBase {
	return nil
}

func (h HuffLeaf) getRight() HuffBase {
	return nil
}

func (h HuffLeaf) getFreq() int {
	return h.Freq
}

func (h HuffLeaf) getChar() rune {
	return h.Char
}

type HuffInternal struct {
	Freq  int
	Left  HuffBase
	Right HuffBase
}

func (h HuffInternal) getLeft() HuffBase {
	return h.Left
}

func (h HuffInternal) getRight() HuffBase {
	return h.Right
}
func (h HuffInternal) getFreq() int {
	return h.Freq
}

func (h HuffInternal) isLeaf() (bool, HuffBase) {
	return false, h
}

func (h HuffInternal) getChar() rune {
	return ' '
}

func bucketSort(mp map[rune]int, length int) [][]rune {
	var bucket = make([][]rune, length+1)
	for k, v := range mp {
		bucket[v] = append(bucket[v], rune(k))
	}
	return bucket
}

func reverse(arr string) string {
	sb := strings.Builder{}
	for i := len(arr) - 1; i >= 0; i-- {
		sb.WriteByte(arr[i])
	}
	return sb.String()
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
func consHuffTree(huffSlice []HuffBase) HuffInternal {

	for len(huffSlice) > 1 {
		first := huffSlice[0]
		second := huffSlice[1]
		huffSlice = huffSlice[2:]
		newNode := HuffInternal{Freq: first.getFreq() + second.getFreq(), Left: first, Right: second}
		i := 0
		for i < len(huffSlice) && huffSlice[i].getFreq() < newNode.getFreq() {
			i++
		}

		temp := append([]HuffBase{}, huffSlice[:i]...)
		temp = append(temp, newNode)
		huffSlice = slices.Concat(temp, huffSlice[i:])

	}

	if len(huffSlice) == 1 {
		return huffSlice[0].(HuffInternal)
	} else {
		return HuffInternal{}
	}
}

func binaryToInt(str string) int {
	str = reverse(str)
	var ans int
	for i, x := range str {
		ans += powInt(2, i) * int(x-'0')
	}
	return ans
}

func dfs(node HuffBase, codes map[rune]string, curr string) {
	isLeaf, x := node.isLeaf()
	if isLeaf {
		codes[node.getChar()] = curr
		return
	}

	dfs(x.getLeft(), codes, curr+"0")
	dfs(x.getRight(), codes, curr+"1")

}

func huffmanCoding(fd io.Reader) (map[rune]string, error) {

	rd := bufio.NewReader(fd)
	buff := make([]byte, 1024*4)
	count := make(map[rune]int)
	for {
		_, err := rd.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		for _, x := range buff {
			count[rune(x)] += 1
		}

	}
	max := 0
	for _, v := range count {
		if v > max {
			max = v
		}
	}
	buckets := bucketSort(count, max)

	huffSlice := []HuffBase{}
	for i, bucket := range buckets {
		for _, x := range bucket {
			huffSlice = append(huffSlice, HuffLeaf{Freq: i, Char: rune(x)})
		}
	}
	huffTree := consHuffTree(huffSlice)

	var huffCodes map[rune]string = make(map[rune]string)
	dfs(huffTree, huffCodes, "")
	for k, v := range huffCodes {
		fmt.Printf("%q - %s\n", k, v)
	}

	return huffCodes, nil
}

func Compress(fd *os.File, codes map[rune]string) error {
	buff := make([]byte, 1024*1024*4)
	newFile, _ := os.Create(fmt.Sprintf("%s_compressed.huff", fd.Name()))
	for {
		n, err := fd.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		bitBuffer := uint8(0)
		bitCount := uint8(0)
		for _, x := range buff[:n] {
			code := codes[rune(x)]
			for _, bit := range code {
				if bit == '1' {
					bitBuffer |= 1 << (7 - bitCount)
				}
				bitCount++
				if bitCount == 8 {
					_, err = newFile.Write([]byte{bitBuffer})
					if err != nil {
						return err
					}
					bitBuffer = 0
					bitCount = 0
				}
			}

		}
		if bitCount > 0 {
			_, err = newFile.Write([]byte{bitBuffer})
			if err != nil {
				return err
			}
		}

		return nil

	}
	return nil
}

func main() {
	args := os.Args
	l := len(args)
	switch {
	case l < 3:
		fmt.Fprintln(os.Stderr, "please give filename...")
		os.Exit(1)
	default:
		if args[1] == "--file" {
			fileName := args[2]

			fd, err := os.Open(fileName)
			if err != nil {
				log.Fatalln("file does not exits")
			}
			codes, err := huffmanCoding(fd)
			fd.Close()
			if err != nil {
				log.Fatal("Failed...")
			}
			fd, err = os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			Compress(fd, codes)

		}
	}
}
