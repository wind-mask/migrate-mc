package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)

func (c *Client) SearchProjects(query string, facets []modrinth.Facet, index modrinth.Index, offset int, limit int) (modrinth.SearchResult, error) {
	log.Println("SearchProjects")
	facetsStr := "&facets=\"["
	for _, facet := range facets {
		v := reflect.ValueOf(facet)
		t := v.Type()
		parts := "["
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			value := fmt.Sprintf("%v", field.Interface())
			if value == "" {
				continue
			}
			fieldName := toLowerFirst(t.Field(i).Name)
			parts += fmt.Sprintf("\"%s%s\",", fieldName, value)
		}
		parts = strings.TrimSpace(strings.TrimSuffix(parts, ",") + "],")
		if parts == "[]," {
			continue
		} else {
			facetsStr += parts

		}
	}
	facetsStr = strings.TrimSuffix(facetsStr, ",") + "]\""
	if facetsStr == "&facets=\"[]\"" {
		facetsStr = ""
	}
	api := fmt.Sprintf("search?query=%s&index=%s&offset=%d&limit=%d%s", query, index, offset, limit, facetsStr)
	resp, err := c.get(api)
	log.Println(resp.Status)
	if err != nil {
		return modrinth.SearchResult{}, err
	}
	var result modrinth.SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
func (c *Client) ListProjectsVersions(idORslug string, loaders []string, game_versions []string) ([]modrinth.Version, error) {
	log.Println("ListProjectsVersions")
	for i, loader := range loaders {
		loaders[i] = "'" + loader + "'"
	}
	for i, game_version := range game_versions {
		game_versions[i] = "'" + game_version + "'"
	}
	query := []string{}
	if len(game_versions) != 0 {
		game_versionsStr := "game_versions=[" + strings.Join(game_versions, ",") + "]"
		query = append(query, game_versionsStr)
	}
	if len(loaders) != 0 {
		loadersStr := "loaders=[" + strings.Join(loaders, ",") + "]"
		query = append(query, loadersStr)
	}
	queryStr := strings.Join(query, "&")
	api := fmt.Sprintf("project/%s/version?%s", idORslug, queryStr)
	// log.Println(api)
	resp, err := c.get(api)
	log.Println(resp.Status)
	if err != nil {
		return nil, err
	}
	var result []modrinth.Version
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
func (c *Client) GetLatestVersionFromHash(hashAlgorithm modrinth.HashAlgorithm, hash string, getLatestVersionFromHashBody modrinth.GetLatestVersionFromHashBody) (modrinth.Version, error) {
	log.Println("GetLatestVersionFromHash")
	api := fmt.Sprintf("version_file/%s/update?algorithm=%s", hash, hashAlgorithm)
	body, err := json.Marshal(getLatestVersionFromHashBody)
	if err != nil {
		return modrinth.Version{}, err
	}
	// log.Println(string(body))
	log.Println(api)
	reader := bytes.NewReader(body)
	resp, err := c.post(api, reader)
	log.Println(resp.Status)
	if err != nil {
		return modrinth.Version{}, err
	}
	var result modrinth.Version
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
