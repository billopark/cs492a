package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type CC struct { // important: public method should start with capital letter
}

func (c *CC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	list := map[string]string{"a": "100", "b": "100", "c": "100", "d": "100"}

	for k, v := range list {
		if err := stub.PutState(k, []byte(v)); err != nil {
			return shim.Error(err.Error())
		}
	}

	// TODO: Initialize bank transactions
	// TODO: 1. remove previous all compositeTxs
	// TODO: 2. add an initialized transaction ("bank"~"1000"~"txID")
	name := "bank"
	txid := stub.GetTxID()
	compositeIndexName := "varName~value~txID"
	defer func() {
		key, err := stub.CreateCompositeKey(compositeIndexName, []string{name, "1000", txid})
		if err != nil {
			fmt.Println(err.Error())
		}
		err = stub.PutState(key, []byte{0x00})
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	states, err := stub.GetStateByPartialCompositeKey(compositeIndexName, []string{name})
	if err != nil {
		shim.Error(err.Error())
	}
	if states == nil {
		return shim.Success([]byte("OK"))
	}

	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		if state == nil {
			continue
		}

		err = stub.DelState(state.GetKey())
		if err != nil {
			shim.Error(err.Error())
		}
	}

	return shim.Success([]byte("OK"))
}

func (c *CC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	var f, args = stub.GetFunctionAndParameters()
	switch f {
	case "init":
		return c.Init(stub)
	case "send":

		money, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error(err.Error())
		}
		aValByte, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		bValByte, err := stub.GetState(args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		aVal, err := strconv.ParseFloat(string(aValByte), 64)
		if err != nil {
			return shim.Error(err.Error())
		}

		bVal, err := strconv.ParseFloat(string(bValByte), 64)
		if err != nil {
			return shim.Error(err.Error())
		}

		aVal -= money
		bVal += money

		if aVal < 0 {
			return shim.Error(errors.New("NON-POSITIVE BALANCE").Error())
		}

		aValUpdatedByte := fmt.Sprintf("%f", aVal)
		bValUpdatedByte := fmt.Sprintf("%f", bVal)

		if err = stub.PutState(args[0], []byte(aValUpdatedByte)); err != nil {
			return shim.Error(err.Error())
		}
		if err = stub.PutState(args[1], []byte(bValUpdatedByte)); err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(args[0] + ":" + aValUpdatedByte + ", " + args[1] + ":" + bValUpdatedByte))
	case "withdraw":
		money, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return shim.Error(err.Error())
		}

		feeRate, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error(err.Error())
		}

		aValByte, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}

		aVal, err := strconv.ParseFloat(string(aValByte), 64)
		if err != nil {
			return shim.Error(err.Error())
		}

		aVal -= money * (1 + feeRate)

		if aVal < 0 {
			return shim.Error(errors.New("NON-POSITIVE BALANCE").Error())
		}

		aValUpdatedByte := fmt.Sprintf("%f", aVal)

		if err = stub.PutState(args[0], []byte(aValUpdatedByte)); err != nil {
			return shim.Error(err.Error())
		}

		// TODO: Increase money*feeRate to "bank" by a delta computation
		name := "bank"
		txid := stub.GetTxID()
		compositeIndexName := "varName~value~txID"

		bankVal, err := c.getBank(stub, compositeIndexName)
		if err != nil {
			return shim.Error(err.Error())
		}

		updatedBankStr := fmt.Sprintf("%f", bankVal + money * feeRate)

		key, err := stub.CreateCompositeKey(compositeIndexName, []string{name, fmt.Sprintf("%f", money * feeRate), txid})
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(key, []byte{0x00})
		if err != nil {
			return shim.Error(err.Error())
		}

		// TODO: Return updated values of "args[0]" and "bank"
		return shim.Success([]byte(args[0] + ":" + aValUpdatedByte + ", " + "bank" + ":" + updatedBankStr))
	}
	return shim.Error("No function is supported for " + f)
}

func (c *CC) getBank(stub shim.ChaincodeStubInterface, indexName string) (float64, error) {
	var finalVal float64

	deltaResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{"bank"})
	if err != nil {
		return finalVal, err
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return finalVal, errors.New("no results for bank")
	}

	// Iterate through result set and compute final value
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return finalVal, nextErr
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return finalVal, splitKeyErr
		}

		// Retrieve the delta value and operation
		valueStr := keyParts[1]

		// Convert the value string and perform the operation
		value, convErr := strconv.ParseFloat(valueStr, 64)
		if convErr != nil {
			return finalVal, convErr
		}

		finalVal += value
	}

	return finalVal, nil
}

func main() {
	err := shim.Start(new(CC))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Start simple chaincode now")
}
