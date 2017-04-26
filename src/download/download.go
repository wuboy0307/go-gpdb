package download

import (
	"log"
	"fmt"
	"encoding/json"
	"errors"
	"strings"
)

import (
	"../../pkg/download/objects"
	"../../pkg/download/library"
	"../../pkg/core/methods"
)

// Extract all the Pivotal Network Product from the Product API page.
func extract_product() objects.ProductObjects {
	log.Println("Obtaining the product ID")

	// Get the API from the Pivotal Products URL
	ProductApiResponse, err := library.GetApi(objects.Products, false, "", 0)
	methods.Fatal_handler(err)

	// Store all the JSON on the Product struct
	json.Unmarshal(ProductApiResponse, &objects.ProductJsonType)

	// Return the struct
	return objects.ProductJsonType
}


// Extract all the Releases of the product with slug : pivotal-gpdb
func extract_release(productJson objects.ProductObjects) (objects.ReleaseObjects, error) {

	// Check what is the URL for the all the releases for product of our interest
	for _, product := range productJson.Products {
		if product.Slug == objects.ProductName {
			objects.ReleaseURL = product.Links.Releases.Href
			objects.PivotalProduct = product.Name
		}
	}

	log.Println("Obtaining the releases for product: " + objects.PivotalProduct + "\n")

	// If we do find the release URL, lets continue
	if objects.ReleaseURL != "" {

		// Extract all the releases
		ReleaseApiResponse, err := library.GetApi(objects.ReleaseURL, false, "", 0)
		methods.Fatal_handler(err)

		// Store all the releases on the release struct
		json.Unmarshal(ReleaseApiResponse, &objects.ReleaseJsonType)

	} else { // Else lets error out
		return objects.ReleaseJsonType, errors.New("Cannot find any Release URL for slug ID: " + objects.ProductName)
	}

	// Return the release struct
	return objects.ReleaseJsonType, nil
}

// provide choice of which version to download
func show_available_version(ReleaseJson objects.ReleaseObjects) (string, string, error) {

	// Get all the releases from the ReleaseJson
	for _, release := range ReleaseJson.Release {
		objects.ReleaseOutputMap[release.Version] = release.Links.Self.Href
		objects.ReleaseVersion = append(objects.ReleaseVersion, release.Version)
	}

	// Sort all the keys
	for index, version := range objects.ReleaseVersion {
		fmt.Printf("%d: %s\n", index+1, version)

	}

	// Total accepted values that user can enter
	objects.TotalOptions = len(objects.ReleaseVersion)

	// Ask user for choice
	users_choice := library.Prompt_choice()
	version_selected := objects.ReleaseOutputMap[objects.ReleaseVersion[users_choice-1]]

	return objects.ReleaseVersion[users_choice-1], version_selected, nil
}

// From the user choice extract all the files available on that version
func extract_downloadURL(ver string, url string) (objects.VersionObjects, error){

	log.Println("Obtaining the files under the greenplum version: " + ver)

	// Extract all the files from the API
	VersionApiResponse, err := library.GetApi(url, false, "", 0)
	methods.Fatal_handler(err)

	// Load it to struct
	json.Unmarshal(VersionApiResponse, &objects.VersionJsonType)

	// Return the result
	return objects.VersionJsonType, nil
}

// Ask user what file in that version are they interested in downloading
// Default is to download GPDB, GPCC and other with a option from parser
func which_product(versionJson objects.VersionObjects) {

	for _, k := range versionJson.File_groups {

		// Default Download which is download the GPDB for Linux (no choice)
		if strings.Contains(k.Name, objects.DBServer) {
			for _, j := range k.Product_files {
				if strings.Contains(j.Name, objects.FileNameContains) {
					objects.ProductFileURL = j.Links.Self.Href
					objects.DownloadURL = j.Links.Download.Href
				}
			}
		}
		// GPCC option

		// Other or fallback method

	}

}

// extract the filename and the size of the product that the user has choosen
func extract_filename_and_size (url string) {

	log.Println("Extracting the filename and the size of the product file.")

	// Obtain the JSON response of the product file API
	ProductFileApiResponse, err := library.GetApi(url, false , "" , 0)
	methods.Fatal_handler(err)

	// Store it on JSON
	json.Unmarshal(ProductFileApiResponse, &objects.ProductFileJsonType)

	// Get the filename and the file size
	filename := strings.Split(objects.ProductFileJsonType.Product_file.Aws_object_key, "/")
	objects.ProductFileName = filename[len(filename)-1]
	objects.ProductFileSize = objects.ProductFileJsonType.Product_file.Size

}

func Download() {

	// Authentication validations
	log.Println("Checking if the user is a valid user")
	library.GetApi(objects.Authentication, false, "", 0)

	// Product list
	productJson := extract_product()

	// Release list
	releaseJson, err := extract_release(productJson)
	methods.Fatal_handler(err)

	// What is user's choice
	choice, choice_url, err := show_available_version(releaseJson)
	methods.Error_handler(err)

	// Get all the files under that version
	versionFileJson, _ := extract_downloadURL(choice, choice_url)

	// The users choice to what to download from that version
	which_product(versionFileJson)

	// Extract the filename and the size of the file
	extract_filename_and_size(objects.ProductFileURL)

	// Download the version
	library.GetApi(objects.DownloadURL, true, objects.ProductFileName, objects.ProductFileSize)

}