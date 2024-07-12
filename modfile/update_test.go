package modfile

import (
	"encoding/json"
	"os"
	"testing"
)


func TestExtract_mods(t *testing.T) {
	path := "../mods"
	m, err := Extract_mods(path)
	if err != nil {
		t.Errorf("Extract_mods(%s) failed: %v", path, err)
	}
	println("Mods:", len(m.Mods))
	f, err := os.Create("fabric_mods.json")
	if err != nil {
		t.Errorf("Open mods.json failed: %v", err)
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	err = e.Encode(m)
	if err != nil {
		t.Errorf("NewEncoder failed: %v", err)
	}
}
func TestUpdateToPath(t *testing.T){
	path:="../new_mods"
	mods,err:=Extract_mods("../mods")
	if err != nil {
		t.Errorf("Extract_mods failed: %v", err)
	}
	u,err:=mods.UpdateToPath("fabric","1.20.1",path)
	if err != nil {
		t.Errorf("UpdateToPath failed: %v", err)
	}

	println("Mods update status:", u)
}