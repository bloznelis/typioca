package words

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func LoadWordSources(paths []string) ([]WordSource, error) {
	sources := make([]WordSource, len(paths))
	for _, path := range paths {
		var wordSource WordSource
		err := LoadWordList(path, wordSource)
		if err != nil {
			return sources, err
		}
	}

	return sources, nil
}

func ReadMetadata(path string) Metadata {
	var metadata Metadata

	fh, err := os.Open(path)
	defer fh.Close()
	check(err)

	decoder := json.NewDecoder(fh)
	decoder.Decode(&metadata)

	return metadata
}

func ReadMetadatas(paths []string) []Metadata {
	var acc []Metadata

	for _, elem := range paths {
		var metadata Metadata

		fh, err := os.Open(elem)
		check(err)

		decoder := json.NewDecoder(fh)
		decoder.Decode(&metadata)

		acc = append(acc, metadata)

		fh.Close()
	}

	return acc
}

func LoadWordList(filePath string, wordSourceToFill WordSource) error {
	fh, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	decoder.Decode(&wordSourceToFill)

	return nil
}

func DeleteWordList(filePath string) error {
	return os.Remove(filePath)
}

func SyncWordList(localPath string, remoteURI string) {
	DownloadFile(localPath, remoteURI)
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		return errors.New("Non 200 response")
	} else {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return err

	}
}
