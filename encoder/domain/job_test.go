package domain_test

import (
	"encoder/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("output_path", "Converted", video)
	assert.NotNil(t, job)
	assert.NoError(t, err)

}
