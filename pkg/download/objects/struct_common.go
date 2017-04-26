package objects

type LinksType struct {
	Self   HrefType `json:"self"`
	Releases   HrefType `json:"releases"`
	Product_files   HrefType `json:"product_files"`
	File_groups   HrefType `json:"file_groups"`
	Signature_file_download HrefType `json:"signature_file_download"`
	Eula_acceptance HrefType `json:"eula_acceptance"`
	User_groups HrefType `json:"user_groups"`
	Download HrefType `json:"download"`
}

type HrefType struct {
	Href string `json:"href"`
}

type eulaType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Links LinksType `json:"_links"`
}