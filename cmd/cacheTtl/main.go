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

//simple cli wrapper as example
func main() {
	cache := cacheTtl.New()
	scanner := bufio.NewScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	logger.Println("hello")
	for scanner.Scan() {
		if err := processLine(scanner.Text(), cache, out); err != nil {
			logger.Println(err)
		}
	}
}

func processLine(line string, cache *cacheTtl.Cache, out *bufio.Writer) error {
	if len(line) < 4 {
		return errors.New("error with parsing input: len(line) < 4")
	}
	keyStart := 4
	switch strings.ToUpper(line[:keyStart]) {
	case "SET ":
		kvLine := line[keyStart:]
		if err := set(kvLine, cache); err != nil {
			return errors.New("error with set: " + err.Error())
		}
		logger.Println("set")
	case "GET ":
		key := line[keyStart:]
		val, err := cache.Get(key)
		if err != nil {
			return errors.New("error with getting: " + err.Error())
		}
		out.WriteString(val.(string) + "\n") // val is always string
		out.Flush()
	case "DEL ":
		key := line[keyStart:]
		if err := cache.Del(key); err != nil {
			return errors.New("error with deleting: " + err.Error())
		}
		logger.Println("del")
	case "LOAD":
		if err := load(cache); err != nil {
			return errors.New("error with loading: " + err.Error())
		}
		logger.Println("load")
	case "SAVE":
		if err := save(cache); err != nil {
			return errors.New("error with saving: " + err.Error())
		}
		logger.Println("save")
	default:
		return errors.New("error with parsing input: unknown command")
	}
	return nil
}

func set(kvLine string, cache *cacheTtl.Cache) error {
	keyEnd := strings.Index(kvLine, " ")
	if keyEnd == -1 {
		return errors.New("not found key end")
	}
	key := kvLine[:keyEnd]
	value := kvLine[keyEnd+1:]
	return cache.Set(key, value, time.Now().Add(ttl))
}

type PersistEl struct {
	Key string
	Val string
	Ttl time.Time
}

func load(cache *cacheTtl.Cache) error {
	file, err := os.Open(persistFile)
	if err != nil {
		return errors.New("error with opening data.dat file: " + err.Error())
	}
	dec := gob.NewDecoder(bufio.NewReader(file))
	var el PersistEl
	for {
		err := dec.Decode(&el)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.New("error with loading data.dat file: " + err.Error())
		}
		if err := cache.Set(el.Key, el.Val, el.Ttl); err != nil {
			return errors.New("error with set:" + err.Error())
		}
	}
}

func save(cache *cacheTtl.Cache) error {
	file, err := os.Create(persistFile)
	if err != nil {
		return errors.New("error with opening data.dat file: " + err.Error())
	}
	defer file.Close()
	b := bufio.NewWriter(file)
	defer b.Flush()
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
	return cache.Walk(fn)
}
