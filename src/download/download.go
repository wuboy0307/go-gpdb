package download

import (
	log "../../pkg/core/logger"
	"fmt"
	"encoding/json"
	"errors"
	"strings"
)

import (
	"../../pkg/download/objects"
	"../../pkg/download/library"
	"../../pkg/core/methods"
	"../../pkg/core/arguments"
)

// Extract all the Pivotal Network Product from the Product API page.
func extract_product() (objects.ProductObjects, error) {

	log.Println("Obtaining the product ID")

	// Get the API from the Pivotal Products URL
	ProductApiResponse, err := library.GetApi(objects.Products, false, "", 0)
	if err != nil {return objects.ProductObjects{}, err}

	// Store all the JSON on the Product struct
	json.Unmarshal(ProductApiResponse, &objects.ProductJsonType)

	// Return the struct
	return objects.ProductJsonType, nil
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

	log.Println("Obtaining the releases for product: " + objects.PivotalProduct)

	// If we do find the release URL, lets continue
	if objects.ReleaseURL != "" {

		// Extract all the releases
		ReleaseApiResponse, err := library.GetApi(objects.ReleaseURL, false, "", 0)
		if err != nil {return objects.ReleaseObjects{}, err}

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

	// Local storehouse
	var version_selected string
	var release string

	// Get all the releases from the ReleaseJson
	for _, release := range ReleaseJson.Release {
		objects.ReleaseOutputMap[release.Version] = release.Links.Self.Href
		objects.ReleaseVersion = append(objects.ReleaseVersion, release.Version)
	}

	// Check if the user provided version is on the list we have just extracted
	if methods.Contains(objects.ReleaseVersion, arguments.RequestedDownloadVersion) {
		log.Println("Found the GPDB version \""+ arguments.RequestedDownloadVersion +"\" on PivNet, continuing..")
		version_selected = objects.ReleaseOutputMap[arguments.RequestedDownloadVersion]
		release = arguments.RequestedDownloadVersion

	} else { // If its not on the list then fallback to interactive mode

		// Print warning if the user did provide a value of the version
		if arguments.RequestedDownloadVersion != "" {
			log.Warn("Unable to find the GPDB version \""+ arguments.RequestedDownloadVersion +"\" on PivNet, failing back to interactive mode..\n")
		} else { // print a blank line
			fmt.Println()
		}

		// Sort all the keys
		for index, version := range objects.ReleaseVersion {
			fmt.Printf("%d: %s\n", index+1, version)
		}

		// Total accepted values that user can enter
		objects.TotalOptions = len(objects.ReleaseVersion)

		// Ask user for choice
		users_choice := library.Prompt_choice()

		// Selected by the user
		version_selected = objects.ReleaseOutputMap[objects.ReleaseVersion[users_choice-1]]
		release = objects.ReleaseVersion[users_choice-1]
	}

	return release, version_selected, nil
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
func which_product(versionJson objects.VersionObjects) error {

	// Clearing up the buffer to ensure we are using a clean array and MAP
	objects.ProductOutputMap = make(map[string]string)
	objects.ProductOptions = []string{}

	// Not sure why from 5.0, all the files are listed inside Product file
	// Since the list is not a correct and its not what all the other product
	// follows , will ask the user to choose from it
	if len(versionJson.Product_files) >= 2 {
		for _,k := range versionJson.Product_files {
			objects.ProductOutputMap[k.Name] = k.Links.Self.Href
			objects.ProductOptions = append(objects.ProductOptions, k.Name)
		}
	} else { // This is the correct API, all the files are inside the file group MAP
		for _, k := range versionJson.File_groups {

			// GPDB Options
			if arguments.RequestedDownloadProduct == "gpdb" {
				// Default Download which is download the GPDB for Linux (no choice)
				if strings.Contains(k.Name, objects.DBServer) {
					for _, j := range k.Product_files {
						for _, name := range objects.FileNameContains {
							if strings.Contains(j.Name, name) {
								objects.ProductFileURL = j.Links.Self.Href
								objects.DownloadURL = j.Links.Download.Href
							}
						}
					}
				}
			}

			// GPCC option
			if arguments.RequestedDownloadProduct == "gpcc" {
				if strings.Contains(k.Name, objects.GPCC) {
					for _, j := range k.Product_files {
						objects.ProductOutputMap[j.Name] = j.Links.Self.Href
						objects.ProductOptions = append(objects.ProductOptions, j.Name)
					}
				}
			}

			// Others or fallback method
			if arguments.RequestedDownloadProduct == "gpextras" {
				for _, j := range k.Product_files {
					objects.ProductOutputMap[j.Name] = j.Links.Self.Href
					objects.ProductOptions = append(objects.ProductOptions, j.Name)
				}
			}

		}
	}

	// If its GPCC or GPextras, then ask users for choice.
	if (arguments.RequestedDownloadProduct == "gpextras" || arguments.RequestedDownloadProduct == "gpcc") && len(objects.ProductOptions) != 0 {
		fmt.Println()
		for index, product := range objects.ProductOptions {
			fmt.Printf("%d: %s\n", index+1, product)
		}
		objects.TotalOptions = len(objects.ProductOptions)
		users_choice := library.Prompt_choice()
		version_selected_url := objects.ProductOutputMap[objects.ProductOptions[users_choice-1]]
		objects.ProductFileURL = version_selected_url
		objects.DownloadURL = version_selected_url + "/download"
	}
	return nil
}

// extract the filename and the size of the product that the user has choosen
func extract_filename_and_size (url string) error {

	log.Println("Extracting the filename and the size of the product file.")

	// Obtain the JSON response of the product file API
	ProductFileApiResponse, err := library.GetApi(url, false , "" , 0)
	if err != nil {return err}

	// Store it on JSON
	json.Unmarshal(ProductFileApiResponse, &objects.ProductFileJsonType)

	// Get the filename and the file size
	filename := strings.Split(objects.ProductFileJsonType.Product_file.Aws_object_key, "/")
	objects.ProductFileName = filename[len(filename)-1]
	objects.ProductFileSize = objects.ProductFileJsonType.Product_file.Size

	return err

}

func Download() error {

	// Authentication validations
	log.Println("Checking if the user is a valid user")
	_, err := library.GetApi(objects.Authentication, false, "", 0)
	if err != nil {return err}

	// Product list
	productJson, err := extract_product()
	if err != nil {return err}

	// Release list
	releaseJson, err := extract_release(productJson)
	if err != nil {return err}

	// What is user's choice
	choice, choice_url, err := show_available_version(releaseJson)
	if err != nil {return err}

	// Get all the files under that version
	versionFileJson, err := extract_downloadURL(choice, choice_url)
	if err != nil {return err}

	// The users choice to what to download from that version
	err = which_product(versionFileJson)
	if err != nil {return err}

	// If we didn't find the database File, then fall back to interactive mode.
	if (arguments.RequestedDownloadProduct == "gpdb" || arguments.RequestedDownloadProduct == "gpcc") && objects.ProductFileURL == "" {
		log.Warn("Couldn't find binaries for GPDB version \"" + choice + "\", failing back to interactive mode...")
		arguments.RequestedDownloadProduct = "gpextras"
		which_product(versionFileJson)
	}

	// Extract the filename and the size of the file
	err = extract_filename_and_size(objects.ProductFileURL)
	if err != nil {return err}

	// Download the version
	log.Println("Starting downloading of file: " + objects.ProductFileName)
	if objects.DownloadURL != "" {
		fmt.Println(objects.DownloadURL)
		_, err := library.GetApi(objects.DownloadURL, true, objects.ProductFileName, objects.ProductFileSize)
		if err != nil {return err}
	} else {
		return errors.New("Download URL is blank, cannot download the product.")
	}

	return nil
}