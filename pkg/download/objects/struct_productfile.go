package objects

type ProductFilesObjects struct {
	Product_file ProductFilesObjType `json:"product_file"`
}

type ProductFilesObjType struct {
	Id int `json:"id"`
	Aws_object_key string `json:"aws_object_key"`
	Description string `json:"description"`
	Docs_url string `json:"docs_url"`
	File_transfer_status string `json:"file_transfer_status"`
	File_type string `json:"file_version"`
	Has_signature_file string `json:"has_signature_file"`
	Included_files []string `json:"included_files"`
	Md5 string `json:"md5"`
	Sha256 string `json:"sha256"`
	Name string `json:"name"`
	Ready_to_serve bool `json:"ready_to_serve"`
	Released_at string `json:"released_at"`
	Size int64 `json:"size"`
	System_requirements []string `json:"system_requirements"`
	Links LinksType `json:"_links"`
}