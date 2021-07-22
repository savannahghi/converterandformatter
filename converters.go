package converterandformatter

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/savannahghi/serverutils"
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

// ConvertInterfaceMap converts a map[string]interface{} to a map[string]string.
//
// Any conversion errors are written out to the output map instead of being
// returned as error values.
//
// New code is discouraged from using this function.
func ConvertInterfaceMap(inp map[string]interface{}) map[string]string {
	out := make(map[string]string)
	if inp == nil {
		return out
	}
	for k, v := range inp {
		val, ok := v.(string)
		if !ok {
			val = fmt.Sprintf("invalid string value: %#v", v)
			if serverutils.IsDebug() {
				log.Printf(
					"non string value in map[string]interface{} that is to be converted into map[string]string: %#v", v)
			}
		}
		out[k] = val
	}
	return out
}

// MapInterfaceToMapString converts a map with interface{} values to one with string
// values.
//
// It is used to convert a GraphQL (gqlgen) input Map to a map of strings for APIs
// that need map[string]string.
func MapInterfaceToMapString(in map[string]interface{}) (map[string]string, error) {
	out := map[string]string{}
	for k, v := range in {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%v (%T) is not a string", v, v)
		}
		out[k] = s
	}
	return out, nil
}

// ConvertStringMap converts a map[string]string to a map[string]interface{}.
//
// This is done mostly in order to conform to the gqlgen Graphql Map scalar.
func ConvertStringMap(inp map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	if inp == nil {
		return out
	}
	for k, v := range inp {
		val := interface{}(v)
		out[k] = val
	}
	return out
}
