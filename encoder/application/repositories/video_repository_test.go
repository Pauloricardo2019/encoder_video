package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewVideoRepository_Insert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)
	assert.NoError(t, err)
	assert.Equal(t, v.ID, video.ID)

}
