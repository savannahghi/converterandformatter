package converters_and_formatters_test

import (
	"testing"

	convertersandformatters "github.com/savannahghi/converters_and_formatters"
)

func TestModelsIsEntity(t *testing.T) {

	t12 := convertersandformatters.USSDSessionLog{}
	t12.IsEntity()

	t13 := convertersandformatters.PhoneOptIn{}
	t13.IsEntity()
}
