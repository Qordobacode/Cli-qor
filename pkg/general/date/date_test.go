package date

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DateFormat(t *testing.T) {
	update := int64(1555803351000)
	dateString := GetDateFromTimestamp(update)
	assert.NotNil(t, dateString)
}
