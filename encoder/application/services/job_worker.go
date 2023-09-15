package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"os"
	"sync"
	"time"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

var mu = &sync.Mutex{}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, job domain.Job, workerID int) {

	for message := range messageChannel {

		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}
		mu.Lock()
		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}
		jobService.VideoService.Video.ID = uuid.New().String()
		mu.Unlock()
		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		mu.Lock()
		err = jobService.VideoService.InsertVideo()
		mu.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.ID = uuid.New().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		mu.Lock()
		_, err = jobService.JobRepository.Insert(&job)
		mu.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		err = jobService.Start()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChannel <- returnJobResult(job, message, nil)

	}

}

func returnJobResult(job domain.Job, message amqp.Delivery, error error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   error,
	}
	return result
}
