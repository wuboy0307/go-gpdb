package library

import (
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"errors"
	"strconv"
	"fmt"
	log "../../core/logger"
	"time"
)

import (
	"../../core/methods"
	"../../core/arguments"
	"encoding/json"
)

// Progress of download
func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			// Open the file
			file, err := os.Open(path)
			methods.Fatal_handler(err)

			// Get stats of the file
			fi, err := file.Stat()
			methods.Fatal_handler(err)

			// Size now
			size := fi.Size()

			// Display Progress of download
			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100
			var bytesToMB float64 = 1024 * 1024

			fmt.Printf("Downloading file %.2f MB of %.2f MB: %.0f", float64(size)/bytesToMB, float64(total)/bytesToMB, percent)
			fmt.Println("% completed")
		}

		// Download is completed, time to terminate
		if stop {
			log.Println("Downloading completed ....")
			log.Println("Downloaded file available at: " + path)
			break
		}

		// Ask to sleep, before repainting the screen.
		time.Sleep(time.Second)
	}
}

// Get the Json from the provided URL or download the file
// if requested.
func GetApi(method string, urlLink string, download bool, filename string, filesize int64) ([]byte, error) {

	// Get the request
	req, err := http.NewRequest(method, urlLink, nil)
	methods.Fatal_handler(err)

	// Add Header to the Http Request
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Verify there is a token on config file
	if methods.IsValueEmpty(arguments.EnvYAML.Download.ApiToken) {
		return []byte(""), errors.New("Cannot find value for \"API_TOKEN\", check \"config.yml\"")
	} else {
		req.Header.Set("Authorization", "Token " + arguments.EnvYAML.Download.ApiToken)
	}

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {return []byte(""), err}

	// If the status code is not 200, then error out
	if resp.StatusCode != http.StatusOK {

		// If we get 451 even after accepting the EULA, then the user would manually
		// connect to the webpage and accept the EULA as its pivotal legal policy
		// to accept the EULA if you are downloading the product for the first time.
		if resp.StatusCode == 451 {

			// Print to why we got this error
			fmt.Println("\n\x1b[33;1mReceived Status code 451 when access the API, this means as per the pivotal " +
				"legal policy if you are attempting to download the product for the first time you are requested to " +
				"to manually accept the end user license agreement (only one time). please connect " +
				"to PivNet and accept the end user license agreement and then try again, as this step cannot be avoided. Click on the link " +
				"mentioned below to redirect you to website to accept the EULA \x1b[0m")


			// Read the error text and store it
			bodyText, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var f interface{}
			_ = json.Unmarshal(bodyText, &f)
			m := f.(map[string]interface{})

			// Obtain the URL that the user can access to accept the EULA
			for k, v := range m {
				if k == "message" {
					fmt.Println("\n\x1b[35;1m" + v.(string) + "\x1b[0m\n")
				}
			}

		}
		return []byte(""), errors.New("API ERROR: HTTP Status code expected ("+ strconv.Itoa(http.StatusOK) +") / received (" + strconv.Itoa(resp.StatusCode) + "), URL (" + urlLink + ")")
	}

	// Close the body once its done
	defer resp.Body.Close()

	// If its to download the software then download it
	if download {

		// The Size of the file
		size := filesize

		// Fully qualifies path
		filename = arguments.DownloadDir + filename

		// Create th file
		out, err := os.Create(filename)
		if err != nil {return []byte(""), err}

		// Initalize progress bar
		done := make(chan int64)
		go PrintDownloadPercent(done, filename, int64(size))
		defer out.Close()

		// Start Downloading
		n, err := io.Copy(out, resp.Body)
		if err != nil {return []byte(""), err}
		done <- n
	}

	// Read the json
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {return []byte(""), err}
	return bodyText, nil
}

