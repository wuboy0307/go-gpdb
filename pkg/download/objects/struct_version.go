package objects

type VersionObjects struct {
	VersionObjType
}

type VersionObjType struct {
	Id int `json:"id"`
	Version string `json:"version"`
	Release_type string `json:"release_type"`
	Release_date string `json:"release_date"`
	Availability string `json:"availability"`
	Eula   eulaType `json:"eula"`
	End_of_support_date string `json:"end_of_support_date"`
	End_of_guidance_date string `json:"end_of_guidance_date"`
	Eccn string `json:"eccn"`
	License_exception string `json:"license_exception"`
	Controlled bool `json:"controlled"`
	Product_files []verProdType `json:"product_files"`
	File_groups []verFileGroupType `json:"file_groups"`
	Updated_at   string `json:"updated_at"`
	Software_files_updated_at string `json:"software_files_updated_at"`
	Links LinksType `json:"_links"`
}

type verProdType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	File_version string `json:"file_version"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}

type verFileGroupType struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Product_files []verProdType `json:"product_files"`
}