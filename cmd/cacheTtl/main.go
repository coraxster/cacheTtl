package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"github.com/coraxster/cacheTtl"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const ttl = time.Minute
const persistFile = "./data.dat"

var logger = log.New(os.Stderr, "", 0)

//simple cli wrapper
func main() {
	cache := cacheTtl.New()
	scanner := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	os.Stdout.WriteString("Hello!\n")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			logger.Println("error with parsing input: len(line) < 4")
			continue
		}
		keyStart := 4
		switch strings.ToUpper(line[:keyStart]) {
		case "SET ":
			KvLine := line[keyStart:]
			keyEnd := strings.Index(KvLine, " ")
			if keyEnd == -1 {
				logger.Println("not found key end")
				continue
			}
			key := KvLine[:keyEnd]
			value := KvLine[keyEnd+1:]
			if err := cache.Set(key, value, time.Now().Add(ttl)); err != nil {
				logger.Println("error with saving: ", err)
				continue
			}
			out.WriteString("set\n")
		case "GET ":
			key := line[keyStart:]
			val, err := cache.Get(key)
			if err != nil {
				logger.Println("error with getting: ", err)
				continue
			}
			out.WriteString(val.(string) + "\n")
		case "DEL ":
			key := line[keyStart:]
			if err := cache.Del(key); err != nil {
				logger.Println("error with deleting: ", err)
				continue
			}
			out.WriteString("del\n")
		case "LOAD":
			file, err := os.Open(persistFile)
			if err != nil {
				logger.Println("error with opening data.dat file: ", err)
				continue
			}
			dec := gob.NewDecoder(bufio.NewReader(file))
			var el PersistEl
			for {
				err := dec.Decode(&el)
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Println("error with loading data.dat file: ", err)
					continue
				}
				if err := cache.Set(el.Key, el.Val, el.Ttl); err != nil {
					logger.Println("error with saving: ", err)
					continue
				}
			}
			out.WriteString("load\n")
		case "SAVE":
			file, err := os.Create(persistFile)
			if err != nil {
				logger.Println("error with opening data.dat file: ", err)
				continue
			}
			b := bufio.NewWriter(file)
			enc := gob.NewEncoder(b)
			fn := func(key string, val interface{}, ttl time.Time) error {
				valStr := val.(string) // always string
				el := PersistEl{
					Key: key,
					Val: valStr,
					Ttl: ttl,
				}
				if err := enc.Encode(&el); err != nil {
					return errors.New("error with encoding:" + err.Error())
				}
				return nil
			}
			err = cache.Walk(fn)
			b.Flush()
			file.Close()
			if err != nil {
				logger.Println("error with saving data.dat file: ", err)
				continue
			}
			out.WriteString("save\n")
		default:
			logger.Println("error with parsing input: unknown command")
		}
		out.Flush()
	}
}

type PersistEl struct {
	Key string
	Val string
	Ttl time.Time
}
