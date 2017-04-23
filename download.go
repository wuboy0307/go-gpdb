package main

// Import Modules
import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"os"
	"io"
)

// All the Pivnet Url's
const (
	Authentication = "https://network.pivotal.io/api/v2/authentication"
	Products =  "https://network.pivotal.io/api/v2/products"
	Releases = "https://network.pivotal.io/api/v2/products/{prod_slug}/releases"
	fileID = "https://network.pivotal.io/api/v2/products/{prod_slug}/releases/{rel_id}"
	download = "https://network.pivotal.io/api/v2/products/{prod_slug}/releases/{rel_id}/product_files/{prod_id}/download"
)

// Product of our interest
const ProductName = "Pivotal Greenplum"


// Function when called replaces the text string with proper values.
func UrlReplacer(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

// Get the Json
func getApi(urlLink string, download bool) []byte {
	req, err := http.NewRequest("GET", urlLink, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.ExpandEnv("Token ${API_TOKEN}"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("output.txt")
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Upstream ERROR: %d (%s)", resp.StatusCode, urlLink)
		// log.Fatal("Authenication Failure")
	}

	defer resp.Body.Close()
	if download {
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Panic(err)
		}
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return bodyText
}

// String to Json Convertor
func String2JsonConv(ApiResponse []byte) map[string]interface{} {
	var dat map[string]interface{}
	if err := json.Unmarshal(ApiResponse, &dat); err != nil {
		panic(err)
	}
	return dat
}

// Main Program
func main() {

	// Store house for all the Api values
	ApiValues := make(map[string]string)
	ReleaseValues := make(map[string]float64)
	ProductValue := make(map[string]interface{})

	// Authentication validations
	log.Println("Checking if the user is a valid user")
	getApi(Authentication, false)

	// Product list
	log.Println("Obtaining the product ID")
	ProductApiResponse := getApi(Products, false)
	dat := String2JsonConv(ProductApiResponse)
	for _, v := range dat["products"].([]interface{}) {
		ProductJSON := v.(map[string]interface{})
		if ProductJSON["name"] == ProductName {
			ApiValues["ProdSlug"] = ProductJSON["slug"].(string)
		}
	}

	// Releases
	log.Println("Obtaining all the releases of the product")
	ReleaseApiResponse := getApi(UrlReplacer(Releases, "{prod_slug}", ApiValues["ProdSlug"]), false)
	dat = String2JsonConv(ReleaseApiResponse)
	for _, v := range dat["releases"].([]interface{}) {
		ReleaseJSON := v.(map[string]interface{})
		ReleaseValues[ReleaseJSON["version"].(string)] = ReleaseJSON["id"].(float64)
	}

	// File ID
	log.Println("Obtaining all the file list of this release")
	FileApiResponse := getApi(UrlReplacer(fileID, "{prod_slug}", ApiValues["ProdSlug"], "{rel_id}", "210"), false)
	dat = String2JsonConv(FileApiResponse)
	for _, v := range dat["file_groups"].([]interface{}) {
		FileIDJSON := v.(map[string]interface{})
		ProductValues := FileIDJSON["product_files"].([]interface {})
		for _,p := range ProductValues {
			ProductList := p.(map[string]interface{})
			ProductValue[ProductList["name"].(string)] = ProductList["id"].(float64)
		}
	}

	// Downloading the software
	log.Println("Downloading the requested software")
	DownloadUrl := UrlReplacer(download, "{prod_slug}", ApiValues["ProdSlug"], "{rel_id}", "210", "{prod_id}", "4233")
	getApi(DownloadUrl, true)

}