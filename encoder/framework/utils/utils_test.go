package utils_test

import (
	"encoder/framework/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsJson(t *testing.T) {

	json := `{
			"id": "45645",
			"file_path": "path",
			"status": "pending"	
			}`

	err := utils.IsJson(json)
	assert.NoError(t, err)

	json = "wes"

	err = utils.IsJson(json)
	assert.Error(t, err)

}
