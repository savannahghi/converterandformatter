package converterandformatter

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// StructToMap converts an object (struct) to a map.
//
// WARNING: int inputs are converted to floats in the output map. This is an
// unintended consequence of converting through JSON.
//
// In future, this should be deprecated.
func StructToMap(item interface{}) (map[string]interface{}, error) {
	bs, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal to JSON: %v", err)
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal from JSON to map: %v", err)
	}
	return res, nil
}

// GenerateRandomWithNDigits - given a digit generate random numbers
func GenerateRandomWithNDigits(numberOfDigits int) (string, error) {
	rangeEnd := int64(math.Pow10(numberOfDigits) - 1)
	value, _ := rand.Int(rand.Reader, big.NewInt(rangeEnd))
	return strconv.FormatInt(value.Int64(), 10), nil
}

// GenerateRandomEmail allows us to get "unique" emails while still keeping
// one main be.well@bewell.co.ke email account
func GenerateRandomEmail() string {
	return fmt.Sprintf("be.well+%v@bewell.co.ke", time.Now().Unix())
}
