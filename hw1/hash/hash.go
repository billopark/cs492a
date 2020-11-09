package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Get hashed string from history block
func Hash(blockID int) (string, error) {
	var ret string
	var err error

	// TODO: Get hash value of taken blockID
	if blockID == 0 {
		return ret, err
	}

	f, err := os.Open("history.block." + strconv.Itoa(blockID))
	if err != nil {
		return ret, err
	}

	defer func() {
		tempErr := f.Close()
		if err != nil {
			err = tempErr
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ret, err
	}

	ret = fmt.Sprintf("%x", h.Sum(nil))
	return ret, err
}
