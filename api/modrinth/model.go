package modrinth

type Facet struct {
	Project_type       string
	Categories         string
	Versions           string
	Client_side        string
	Server_side        string
	Open_source        string
	License            string
	Title              string
	Author             string
	Follows            string
	Project_id         string
	Downloads          string
	Color              string
	Created_timestamp  string
	Modified_timestamp string
}
type Index string
type Project_type string
type Side string
type HashAlgorithm string

const (
	HashAlgorithmSha512 HashAlgorithm = "sha512"
	HashAlgorithmSha1   HashAlgorithm = "sha1"
)
const (
	Project_typeMod           Project_type = "mod"
	Project_typeModpack       Project_type = "modpack"
	Project_typeResource_pack Project_type = "resource_pack"
	Project_typeshader_pack   Project_type = "shader_pack"
)
const (
	SideRequired    Side = "required"
	SideOptional    Side = "optional"
	SideUnsupported Side = "unsupported"
)
const (
	IndexRelevance Index = "relevance"
	IndexDownloads Index = "downloads"
	IndexFollows   Index = "follows"
	IndexNewest    Index = "newest"
	IndexUpdated   Index = "updated"
)

type Project struct {
	Project_id         string       `json:"project_id"`
	Project_type       Project_type `json:"project_type"`
	Slug               string       `json:"slug"`
	Author             string       `json:"author"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Categories         []string     `json:"categories"`
	Display_categories []string     `json:"display_categories"`
	Versions           []string     `json:"versions"`
	Downloads          int          `json:"downloads"`
	Follows            int          `json:"follows"`
	Icon_url           *string      `json:"icon_url"`
	Date_created       string       `json:"date_created"`
	Date_modified      string       `json:"date_modified"`
	Latest_version     string       `json:"latest_version"`
	License            string       `json:"license"`
	Client_side        Side         `json:"client_side"`
	Server_side        Side         `json:"server_side"`
	Gallery            []string     `json:"gallery"`
	Featured_gallery   *string      `json:"featured_gallery"`
}
type SearchResult struct {
	Hits       []Project `json:"hits"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	Total_hits int       `json:"total_hits"`
}

type VersionFileHashes struct {
	Sha512 string `json:"sha512"`
	Sha1   string `json:"sha1"`
}
type VersionFile struct {
	Hashes    VersionFileHashes `json:"hashes"`
	Url       string            `json:"url"`
	Filename  string            `json:"filename"`
	Primary   bool              `json:"primary"`
	Size      int               `json:"size"`
	File_type string            `json:"file_type"`
}
type VersionDependency struct {
	Version_id      string `json:"version_id"`
	Project_id      string `json:"project_id"`
	File_name       string `json:"file_name"`
	Dependency_type string `json:"dependency_type"`
}
type BaseVersion struct {
	Name             string              `json:"name"`
	Version_number   string              `json:"version_number"`
	Changelog        string              `json:"changelog"`
	Dependencies     []VersionDependency `json:"dependencies"`
	Game_versions    []string            `json:"game_versions"`
	Version_type     string              `json:"version_type"`
	Loaders          []string            `json:"loaders"`
	Featured         bool                `json:"featured"`
	Requested_status string              `json:"requested_status"`
}
type Version struct {
	BaseVersion
	Id             string        `json:"id"`
	Project_id     string        `json:"project_id"`
	Author_id      string        `json:"author_id"`
	Date_published string        `json:"date_published"`
	Downloads      int           `json:"downloads"`
	Changelog_url  string        `json:"changelog_url"`
	Files          []VersionFile `json:"files"`
}

type GetLatestVersionFromHashBody struct {
	Loaders       []string `json:"loaders"`
	Game_versions []string `json:"game_versions"`
}
