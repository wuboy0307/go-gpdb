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
func GetApi(urlLink string, download bool, filename string, filesize int64) ([]byte, error) {

	// Get the request
	req, err := http.NewRequest("GET", urlLink, nil)
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
		if err != nil {return []byte(""), errors.New("API ERROR: HTTP Status code expected ("+ strconv.Itoa(http.StatusOK) +") / received (" + strconv.Itoa(resp.StatusCode) + "), URL (" + urlLink + ")")}
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

