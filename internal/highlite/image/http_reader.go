package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const highliteImageLocation = "http://www.highlite.nl/var/StorageHighlite/ProduktBilder/"

// Image download result.
type downloadResponse struct {
	name   string
	err    error
	reader io.ReadCloser
}

// Internal struct for downloading images.
type httpReader struct {
	downloadFn   func(string) (*http.Response, error)
	downloads    chan downloadResponse
	ready        chan bool
	imageReaders map[string]io.ReadCloser
	imageNames   []string
	err          error
}

// Single download job. Gets an image from the internet and writes the result to the result channel.
func (hi *httpReader) downloadImage(name string) {
	response, err := hi.downloadFn(highliteImageLocation + name)
	if err != nil {
		hi.downloads <- downloadResponse{
			name: name,
			err:  err,
		}
	} else {
		hi.downloads <- downloadResponse{
			name:   name,
			reader: response.Body,
		}
	}
}

// Handles download jobs results and notifies listeners that all downloads were completed.
func (hi *httpReader) downloadObserver() {
	for range hi.imageNames {
		download := <-hi.downloads
		if download.err != nil {
			hi.err = download.err
		} else {
			hi.imageReaders[download.name] = download.reader
		}
	}

	close(hi.downloads)
	close(hi.ready)
}

// Closes all readers.
func (hi *httpReader) closeReaders() {
	for _, reader := range hi.imageReaders {
		reader.Close()
	}
}

// Downloads all images in parallel. Closes all image readers if context exceeded or if there is any error.
func (hi *httpReader) downloadImages(ctx context.Context) (map[string]io.ReadCloser, error) {
	for _, name := range hi.imageNames {
		go hi.downloadImage(name)
	}

	go hi.downloadObserver()

	select {

	case <-ctx.Done():
		go func() {
			<-hi.ready
			hi.closeReaders()
		}()

		return nil, fmt.Errorf("Context exceeded while image dowloading")

	case <-hi.ready:
		if hi.err != nil {
			hi.closeReaders()

			return nil, hi.err
		}

		return hi.imageReaders, nil
	}
}
