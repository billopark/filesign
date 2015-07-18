package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const NEWLINE = "\r\n"

func IsErr(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func ReadFile(addr string) (string, error) {
	file, err := os.Open(addr)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	buf := make([]byte, 1024)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

func ReadSignFile(di *map[string][]string) error {
	dic := *di
	fi, err := ioutil.ReadFile("sign.data")
	if err != nil {
		return err
	}
	dat := strings.Split(string(fi), NEWLINE)
	var key []string
	var flag int = 0
	for _, i := range dat {
		flag = 0
		d := strings.Split(i, "=")
		if len(d) != 2 {
			err = errors.New("Data was wrong.")
			if err != nil {
				return err
			}
		}
		for _, j := range key {
			if d[0] == j {
				dic[d[0]] = append(dic[d[0]], d[1])
				flag = 1
				break
			}
		}
		if flag == 1 {
			continue
		}
		v := []string{d[1]}
		dic[d[0]] = v
		key = append(key, d[0])
	}
	return nil
}

func FindExt(di *map[string][]string, hexData string) string {
	dic := *di
	for ext, binL := range dic {
		for _, bin := range binL {
			temp := bin
			for strings.IndexRune(temp, 'n') != -1 {
				i := strings.IndexRune(temp, 'n')
				temp = strings.Join([]string{temp[:i], string(hexData[i]), temp[i+1:]}, "")
			}
			if strings.EqualFold(hexData[:utf8.RuneCountInString(temp)], temp) {
				return ext
			}
		}
	}
	return ""
}

func main() {
	filePtr := flag.String("u", "", "file address")
	flag.StringVar(filePtr, "U", "", "file address")
	flag.Parse()

	for *filePtr == "" {
		fmt.Printf("Input file address:")
		fmt.Scanf("%s\n", filePtr)
	}
	var dic = make(map[string][]string)
	err := ReadSignFile(&dic)

	IsErr(err)
	s, _ := ReadFile(*filePtr)
	ext := FindExt(&dic, s)
	if ext == "" {
		fmt.Printf("No matching Extension")
		os.Exit(1)
	}

	fmt.Printf("%s", ext)
}
