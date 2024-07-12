package modfile

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	lop "github.com/samber/lo/parallel"
	"github.com/wind-mask/migrate-mc/api"
	"github.com/wind-mask/migrate-mc/api/modrinth"
)

type Mod struct {
	FabricMod FabricMod `json:"fabricMod"`
	Hashes    modrinth.VersionFileHashes `json:"hashes"`
}
type modFetch struct {
	modFileReadCloser io.ReadCloser
	modFileNaem       string
}

func (m *Mod) UpdateToFile(updateTo UpdateTo, client *api.Client) (*modFetch, error) {
	// client := api.NewModrinthClient()
	g := modrinth.GetLatestVersionFromHashBody{
		Loaders:       []string{updateTo.Loader},
		Game_versions: []string{updateTo.MinecraftVersion},
	}
	v, err := client.GetLatestVersionFromHash(modrinth.HashAlgorithmSha512, m.Hashes.Sha512, g)
	if err != nil {
		return nil, errors.Join(err, errors.New("GetLatestVersionFromHash failed:"+m.FabricMod.FabricModJson.Name))
	}
	if len(v.Files) == 0 {
		return nil, errors.New(m.FabricMod.FabricModJson.Name + " GetLatestVersionFromHash failed: no files")
	}
	rc, err := client.GetModFileFromVersionFile(v.Files[0])
	if err != nil {
		return nil, errors.Join(err, errors.New("GetModFileFromVersionFile failed:"+m.FabricMod.FabricModJson.Name))
	}
	return &modFetch{modFileReadCloser: rc, modFileNaem: v.Files[0].Filename}, nil
}

type Mods struct {
	MinecraftVersion string `json:"minecraftVersion"`
	Loader           string `json:"loader"`
	LoaderVersion    string `json:"loaderVersion"`
	ModsPath         string `json:"modsPath"`
	Mods             []Mod  `json:"mods"`
}

func Extract_mods(modpath string) (Mods, error) {
	de, err := os.ReadDir(modpath)
	Logger.Println("Mods path:", modpath)
	if err != nil {
		return Mods{}, err
	}
	var mods Mods
	var modsCount int
	for _, file := range de {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".jar") {
			modsCount++
			f, h, err := extract_fabric_mod(modpath + "/" + file.Name())
			if err != nil {
				Logger.Println("Extract_fabric_mod_json failed:", file.Name(), err)

				continue
			}
			mods.Mods = append(mods.Mods, Mod{FabricMod: FabricMod{FabricModJson: f},
				Hashes: h})
		}
	}
	Logger.Println("Mods count:", modsCount)
	return mods, nil
}

type UpdateTo struct {
	Loader           string `json:"loader"`
	MinecraftVersion string `json:"minecraftVersion"`
}

type ModUpdateStatus struct {
	Original Mod `json:"original"`
	UpdateTo UpdateTo `json:"updateTo"`
	Err     error `json:"err"`
}

func (m *Mods) UpdateToPath(loader string, mcVersion string, path string) ([]ModUpdateStatus, error) {
	err := os.MkdirAll(path, fs.ModeDir)
	if err != nil {
		return nil, errors.Join(err, errors.New("os.MkdirAll failed"))
	}
	u := UpdateTo{Loader: loader, MinecraftVersion: mcVersion}
	client := api.NewModrinthClient()
	modsUpdateStatus := lop.Map(m.Mods, func(mod Mod, _ int) ModUpdateStatus {
		modFetch, err := mod.UpdateToFile(u, client)
		if err != nil {
			return ModUpdateStatus{
				Original: mod,
				UpdateTo: u,
				Err:     errors.Join(err, errors.New("UpdateToFile failed")),
			}
		}
		f, err := os.Create(path + "/" + modFetch.modFileNaem)
		if err != nil {
			return ModUpdateStatus{
				Original: mod,
				UpdateTo: u,
				Err:     errors.Join(err, errors.New("os.Create failed")),
			}
		}
		_, err = io.Copy(f, modFetch.modFileReadCloser)
		if err != nil {

			return ModUpdateStatus{
				Original: mod,
				UpdateTo: u,
				Err:     errors.Join(err, errors.New("io.Copy failed")),
			}
		}
		modFetch.modFileReadCloser.Close()
		return ModUpdateStatus{
			Original: mod,
			UpdateTo: u,
			Err:     nil,
		}
	})
	return modsUpdateStatus, nil
}
