package general

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DateFormat(t *testing.T) {
	update := int64(1555803351000)
	dateString := GetDateFromTimestamp(update)
	assert.Equal(t, dateString, "2019-04-21 02:35:51")
}
