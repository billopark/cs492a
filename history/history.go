package history

import (
	"../hash"
	"io/ioutil"
	"strconv"
	"strings"
)

type History struct {
	txs            [10]string
	currentTxId    int
	currentBlockId int
}

func (h *History) Append(tx string) error {
	var err error

	// TODO: store tx into h.txs

	// TODO: write h.txs if h.currentTxId >len(h.txs)

	return err
}

func (h *History) Init() error {
	var err error

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}

	maxID := 1
	for _, file := range files {
		s := file.Name()
		if strings.HasPrefix(s, "history.block.") {
			s = strings.TrimPrefix(s, "history.block.")
			i, err := strconv.Atoi(s)

			if err != nil {
				return err
			}
			maxID = i + 1 // next history
		}
	}

	h.currentBlockId = maxID

	return err
}

func (h *History) Write() error {
	defer func() {
		h.txs = [10]string{}
		h.currentTxId = 0
		h.currentBlockId++
	}()

	var err error

	// TODO: Write the hash value of previous block
	//blockPath := "history.block." + strconv.Itoa(h.currentBlockId)
	//hashValue, err := hash.Hash(h.currentBlockId - 1)
	//s := ""
	//err = ioutil.WriteFile(blockPath, []byte(s), 0644)
	return err
}
