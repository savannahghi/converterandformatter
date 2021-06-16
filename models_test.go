package converterandformatter_test

import (
	"testing"

	"github.com/savannahghi/converterandformatter"
)

func TestModelsIsEntity(t *testing.T) {

	t12 := converterandformatter.USSDSessionLog{}
	t12.IsEntity()

	t13 := converterandformatter.PhoneOptIn{}
	t13.IsEntity()
}
