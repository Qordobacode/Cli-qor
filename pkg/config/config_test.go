package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test struct {
	BaseURL string `mapstructure:"base_url"`
}

func Test_Deccode(t *testing.T) {
	input := map[string]interface{}{
		"base_url": "test",
	}

	var result Test
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result.BaseURL)
	assert.Equal(t, "test", result.BaseURL)

}
