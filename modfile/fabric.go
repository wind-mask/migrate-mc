package modfile

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"

	"github.com/wind-mask/migrate-mc/api/modrinth"
)

type FabricMod struct {
	FabricModJson FabricModJson
}

type Person struct {
	Name    string             `json:"name"`
	Contact ContactInformation `json:"contact"`
}
type person struct {
	Name    string             `json:"name"`
	Contact ContactInformation `json:"contact"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	if data[0] == '{' {
		var single person
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		p.Name = single.Name
		p.Contact = single.Contact
		return nil
	} else if data[0] == '"' {
		var single string
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		p.Name = single
		return nil
	}
	return errors.New("person:invalid json")
}
func (p Person) MarshalJSON() ([]byte, error) {
	if p.Contact == nil {
		return json.Marshal(p.Name)
	}
	return json.Marshal((person)(p))
}

type ContactInformation map[string]string
type FabricModJson struct {
	SchemaVersion int                  `json:"schemaVersion"`
	Id            string               `json:"id"`
	Version       string               `json:"version"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Authors       []Person             `json:"authors"`
	Contact       ContactInformation   `json:"contact"`
	License       StringAdd            `json:"license"`
	Icon          string               `json:"icon"`
	Environment   string               `json:"environment"`
	Depends       map[string]StringAdd `json:"depends"`
}


func extract_fabric_mod(modFilePath string) (FabricModJson, modrinth.VersionFileHashes, error) {
	b, h, err := extractFileFromJar(modFilePath, "fabric.mod.json")
	if err != nil {
		return FabricModJson{}, modrinth.VersionFileHashes{}, err
	}
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))
	var fabricModJson FabricModJson
	err = json.Unmarshal(b, &fabricModJson)
	if err != nil {
		return FabricModJson{}, modrinth.VersionFileHashes{}, err
	}
	return fabricModJson, h, nil
}
func getModFileHash(file io.Reader) (modrinth.VersionFileHashes, error) {
	hasherSha512 := sha512.New()
	hasherSha1 := sha1.New()
	_, err := io.Copy(io.MultiWriter(hasherSha512, hasherSha1), file)
	if err != nil {
		return modrinth.VersionFileHashes{}, err
	}
	return modrinth.VersionFileHashes{
		Sha1:   hex.EncodeToString(hasherSha1.Sum(nil)),
		Sha512: hex.EncodeToString(hasherSha512.Sum(nil)),
	}, nil
}
