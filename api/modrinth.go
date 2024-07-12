package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)

func (c *Client) SearchProjects(query string, facets []modrinth.Facet, index modrinth.Index, offset int, limit int) (modrinth.SearchResult, error) {
	Logger.Println("SearchProjects")
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
	if err != nil {
		return modrinth.SearchResult{}, errors.Join(err, errors.New("get failed"))
	}
	if resp.StatusCode != 200 {
		return modrinth.SearchResult{}, errors.New("get failed: " + resp.Status)
	}
	var result modrinth.SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return modrinth.SearchResult{}, errors.Join(err, errors.New("json.Decode failed"))
	}
	return result, err
}
func (c *Client) ListProjectsVersions(idORslug string, loaders []string, game_versions []string) ([]modrinth.Version, error) {
	Logger.Println("ListProjectsVersions")
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
	if err != nil {
		return nil, errors.Join(err, errors.New("get failed"))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("get failed: " + resp.Status)
	}
	var result []modrinth.Version
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.Join(err, errors.New("json.Decode failed"))
	}
	return result, err
}
func (c *Client) GetLatestVersionFromHash(hashAlgorithm modrinth.HashAlgorithm, hash string, getLatestVersionFromHashBody modrinth.GetLatestVersionFromHashBody) (modrinth.Version, error) {
	Logger.Println("GetLatestVersionFromHash")
	api := fmt.Sprintf("version_file/%s/update?algorithm=%s", hash, hashAlgorithm)
	body, err := json.Marshal(getLatestVersionFromHashBody)
	if err != nil {
		return modrinth.Version{}, errors.Join(err, errors.New("json.Marshal failed"))
	}
	// log.Println(string(body))
	// log.Println(api)
	reader := bytes.NewReader(body)
	resp, err := c.post(api, reader)
	if err != nil {
		return modrinth.Version{}, errors.Join(err, errors.New("post failed"))
	}
	if resp.StatusCode != 200 {
		return modrinth.Version{}, errors.New("post failed: " + resp.Status)
	}
	var result modrinth.Version
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return modrinth.Version{}, errors.Join(err, errors.New("json.Decode failed"))
	}
	return result, err
}
