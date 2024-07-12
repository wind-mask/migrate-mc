package modfile

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)
var Logger = log.Default()
type StringAdd []string

func (sos *StringAdd) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		var single string
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*sos = []string{single}
	} else {
		var slice []string
		if err := json.Unmarshal(data, &slice); err != nil {
			return err
		}
		*sos = slice
	}
	return nil
}
func (sos StringAdd) MarshalJSON() ([]byte, error) {
	if len(sos) == 1 {
		return json.Marshal(sos[0])
	}
	return json.Marshal([]string(sos))
}
// Extract file from jar
func extractFileFromJar(jarPath, fileName string) ([]byte, modrinth.VersionFileHashes, error) {
	file, err := os.Open(jarPath)
	if err != nil {
		return nil, modrinth.VersionFileHashes{}, err
	}
	defer file.Close()
	fileSize, err := file.Stat()
	if err != nil {
		return nil, modrinth.VersionFileHashes{}, err
	}
	// r, err := zip.OpenReader(jarPath)
	hash, err := getModFileHash(file)
	if err != nil {
		return nil, modrinth.VersionFileHashes{}, err
	}
	r, err := zip.NewReader(file, fileSize.Size())
	if err != nil {

		return nil, modrinth.VersionFileHashes{}, err
	}
	for _, f := range r.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return nil, modrinth.VersionFileHashes{}, err
			}
			defer rc.Close()
			buf := make([]byte, f.UncompressedSize64)
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				return nil, modrinth.VersionFileHashes{}, err
			}
			return buf, hash, nil
		}
	}
	return nil, modrinth.VersionFileHashes{}, errors.New("file not found")
}