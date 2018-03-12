package image

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Image download result.
type downloadResponse struct {
	name   string
	err    error
	reader io.ReadCloser
}

// Internal struct for downloading images.
type httpReader struct {
	downloadFn    func(string) (*http.Response, error)
	imageNames    []string
	imageLocation string

	loadsChan chan downloadResponse
	readyChan chan bool
	errorChan chan error

	imageReaders Bucket
}

// Initialization.
func (hi *httpReader) init(
	downloadFn func(string) (*http.Response, error),
	imageNames []string,
	imageLocation string,
) {
	hi.downloadFn = downloadFn
	hi.imageNames = imageNames
	hi.imageLocation = imageLocation

	hi.loadsChan = make(chan downloadResponse)
	hi.readyChan = make(chan bool)
	hi.errorChan = make(chan error, len(imageNames))

	hi.imageReaders = make(Bucket, 0, len(imageNames))
}

// Single download job. Gets an image from the internet and writes the result to the result channel.
func (hi *httpReader) downloadImage(name string) {
	url := hi.imageLocation + name
	response, err := hi.downloadFn(url)

	if err == nil && response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
	}

	if err != nil {
		hi.loadsChan <- downloadResponse{
			name: name,
			err:  fmt.Errorf("failed to download image %s: %s", url, err),
		}
	} else {
		hi.loadsChan <- downloadResponse{
			name:   name,
			reader: response.Body,
		}
	}
}

// Handles download jobs results and notifies listeners that all loadsChan were completed.
func (hi *httpReader) downloadObserver() {
	for range hi.imageNames {
		download := <-hi.loadsChan
		if download.err != nil {
			hi.errorChan <- download.err
		} else {
			hi.imageReaders = append(hi.imageReaders, BucketItem{
				Name:   download.name,
				Reader: download.reader,
			})
		}
	}

	hi.readyChan <- true
}

// Waits until all downloads are completed and then closes all resources.
// Is used only when some error has occurred.
func (hi *httpReader) recoverAfterError() {
	<-hi.readyChan
	hi.cleanup()
	hi.imageReaders.Close()
}

// Closes channels.
func (hi *httpReader) cleanup() {
	close(hi.readyChan)
	close(hi.errorChan)
	close(hi.loadsChan)
}

// Downloads all images in parallel. Closes all image readers if context exceeded or if there is any error.
func (hi *httpReader) downloadImages(ctx context.Context) (Bucket, error) {
	for _, name := range hi.imageNames {
		go hi.downloadImage(name)
	}

	go hi.downloadObserver()

	select {
	case <-ctx.Done():
		go hi.recoverAfterError()
		return nil, fmt.Errorf("failed to download images: context exceeded")

	case err := <-hi.errorChan:
		go hi.recoverAfterError()
		return nil, err

	case <-hi.readyChan:
		hi.cleanup()
		return hi.imageReaders, nil
	}
}
