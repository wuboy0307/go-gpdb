package objects

type ProductObjects struct {
	Products []ProductObjType `json:"products"`
	Links LinksType `json:"_links"`
}

type ProductObjType struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Logo_url string `json:"logo_url"`
	Links   LinksType `json:"_links"`
}