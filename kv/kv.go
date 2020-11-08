package kv

import (
	"../history"
	"fmt"
	"errors"
	"io/ioutil"
	"strings"
)

type DatabaseInterface interface {
	GetState(key string) (string, error)
	PutState(key, value string) error
}

type Database struct {
	state     KV //TODO: state DB
	tempBlock history.History //TODO: history of transactions
}

type KV map[string]string

func (db Database) GetState(key string) (string, error) {
	var ret string
	var err error

	// TODO: Get the value of an input key and return error if it has

	return ret, err
}

func (db *Database) PutState(key, value string) error {
	var err error
	// TODO: Put a key-value pair into state DB and return error if it has

	// TODO: Append history into temp block
	//err = db.tempBlock.Append(key + "\t" + value)

	return err
}

func Init(obj *Database)  {
	obj.state = make(KV)

	// TODO: Load state DB
	// statePath := "state.db"
	// data, err := ioutil.ReadFile(statePath)

	// TODO: Initialize history of transactions
	// err = obj.tempBlock.Init()
}

func Finalize(obj *Database) {
	// TODO: Store current state DB
	// statePath := "state.db"
	// s := ""
	// err := ioutil.WriteFile(statePath, []byte(s), 0644)

	// TODO: Store tempBlock
	// err = obj.tempBlock.Write()
}
