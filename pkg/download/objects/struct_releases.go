package objects

type ReleaseObjects struct {
	Release []releaseObjType `json:"releases"`
	Links LinksType `json:"_links"`
}

type releaseObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Release_notes_url string `json:"release_notes_url"`
	Availability string `json:"availability"`
	Description string `json:"description"`
	Eula   eulaType `json:"eula"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}