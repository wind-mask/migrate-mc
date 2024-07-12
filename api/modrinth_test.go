package api

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)

func TestSearchProjects(t *testing.T) {

	var ModrinthClient = NewModrinthClient()

	// Search for projects
	facet := modrinth.Facet{}
	facets := []modrinth.Facet{facet}
	p, err := ModrinthClient.SearchProjects("appleskin", facets, modrinth.IndexRelevance, 0, 10)
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))

}
func TestListProjectsVersions(t *testing.T) {
	var ModrinthClient = NewModrinthClient()

	// List project versions
	loaders := []string{"fabric"}
	game_versions := []string{"1.20.1"}
	versions, err := ModrinthClient.ListProjectsVersions("appleskin", loaders, game_versions)
	if err != nil {
		t.Error(err)
	}

	b, err := json.Marshal(versions)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestGetLatestVersionFromHash(t *testing.T) {
	var ModrinthClient = NewModrinthClient()

	// Get latest version from hash
	hasher := sha512.New()
	modFilePath := "../appleskin-fabric-mc1.19.2-2.5.1.jar"
	file, err := os.Open(modFilePath)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		t.Error(err)
	}
	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))
	body := modrinth.GetLatestVersionFromHashBody{
		Loaders:       []string{"fabric"},
		Game_versions: []string{"1.20.1"},
	}
	v, err := ModrinthClient.GetLatestVersionFromHash(modrinth.HashAlgorithmSha512, hashStr, body)
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
	println(hashStr)
}
