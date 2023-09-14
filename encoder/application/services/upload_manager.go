package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objPath string, client *storage.Client, ctx context.Context) error {

	path := strings.Split(objPath, os.Getenv("localStoragePath")+"/")

	f, err := os.Open(objPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

func (vu *VideoUpload) loadPaths() error {

	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {

	in := make(chan int, runtime.NumCPU()) // Qual o arquivo baseado na posição do slice paths
	returnChannel := make(chan string)

	err := vu.loadPaths() //carregas todos os paths
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload() //Pegar o client do google cloud, e o context
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ { //Iniciar diversas rotinas de acordo com a concurrency
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx) //Ficará lendo o canal in, e quando tiver algo, vai chamar o uploadWorker
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ { //Percorrer todos os paths, e atribui a posição do slice no canal in
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel { //Lendo o canal returnChannel, a cada upload que retornar
		if r != "" { //Caso tenha dado erro, adiciona a message no doneUpload e para toda a execução
			doneUpload <- r
			break
		}
	}
	return nil
}

func (vu *VideoUpload) uploadWorker(in <-chan int, returnChannel chan string, uploadClient *storage.Client, ctx context.Context) {

	for x := range in { //Lendo o channel in, e atribuindo a posição do slice paths
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx) //Fazendo o upload do arquivo
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during the upload: %v. error: %v", vu.Paths[x], err)
			returnChannel <- err.Error() //Se der erro, retorna o erro no channel returnChannel
		}
		returnChannel <- "" //Se não der erro, retorna uma string vazia no channel returnChannel
	}

	returnChannel <- "upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
