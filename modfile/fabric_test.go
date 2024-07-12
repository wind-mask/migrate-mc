package modfile

import (
	"encoding/json"
	"os"
	"testing"
)

func TestExtract_fabric_mod(t *testing.T) {
	modFilePath := "../appleskin-fabric-mc1.19.2-2.5.1.jar"
	j, _, err := extract_fabric_mod(modFilePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	f, err := os.Create("fabric_mod.json")

	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	err = e.Encode(j)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

}
