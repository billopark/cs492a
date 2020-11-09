package kv

import (
	"github.com/billopark/cs492a/hw1/history"
	"fmt"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type DatabaseInterface interface {
	GetState(key string) (string, error)
	PutState(key, value string) error
}

type Database struct {
	state     KV              //TODO: state DB
	tempBlock history.History //TODO: history of transactions
}

type KV map[string]string

func (db Database) GetState(key string) (string, error) {
	var ret string
	var err error

	// TODO: Get the value of an input key and return error if it has
	if val, ok := db.state[key]; !ok {
		err = errors.New("no such key")
	} else {
		ret = val
	}

	return ret, err
}

func (db *Database) PutState(key, value string) error {
	var err error
	// TODO: Put a key-value pair into state DB and return error if it has
	db.state[key] = value

	// TODO: Append history into temp block
	err = db.tempBlock.Append(key + "\t" + value)

	return err
}

func Init(obj *Database) {
	obj.state = make(KV)

	// TODO: Load state DB
	statePath := "state.db"
	data, err := ioutil.ReadFile(statePath)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			kv := strings.Split(line, "\t")
			if len(kv) != 2 {
				fmt.Println("Invalid state file")
				return
			}
			obj.state[kv[0]] = kv[1]
		}

	} else if !os.IsNotExist(err) {
		fmt.Printf("%v\n", err)
		return

	} else {
		err = nil
	}

	// TODO: Initialize history of transactions
	err = obj.tempBlock.Init()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func Finalize(obj *Database) {
	// TODO: Store current state DB
	statePath := "state.db"
	s := strings.Join(func() []string {
		var stateList []string
		for k, v := range obj.state {
			stateList = append(stateList, k+"\t"+v)
		}
		return stateList
	}(), "\n")

	err := ioutil.WriteFile(statePath, []byte(s), 0644)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// TODO: Store tempBlock
	err = obj.tempBlock.Write()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
