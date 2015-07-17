package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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
	dat := strings.Split(string(fi), "\n")
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
func main() {
	var dic = make(map[string][]string)
	err := ReadSignFile(&dic)
	IsErr(err)
	s, _ := ReadFile("filesign.exe")
	for ext, binL := range dic {
		for _, bin := range binL {
			if strings.EqualFold(bin, s[:len(bin)]) {
				fmt.Printf("%s", ext)
			}
		}
	}
}
