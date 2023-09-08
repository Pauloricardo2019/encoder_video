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

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "Converted", video)
	assert.NoError(t, err)

	repoJob := repositories.NewJobRepository(db)

	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, j.ID)
	assert.Equal(t, j.ID, job.ID)

}
