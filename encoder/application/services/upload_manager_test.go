package services_test

import (
	"encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("paulotest")
	assert.NoError(t, err)

	err = videoService.Fragment()
	assert.NoError(t, err)
	err = videoService.Encode()
	assert.NoError(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "paulotest"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload

	assert.Equal(t, result, "upload completed")

}
