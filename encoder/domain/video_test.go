package domain_test

import (
	"encoder/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {

	video := domain.NewVideo()
	err := video.Validate()
	assert.Error(t, err)

}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "abc"
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	assert.Error(t, err)

}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.New().String()
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	assert.NoError(t, err)

}
