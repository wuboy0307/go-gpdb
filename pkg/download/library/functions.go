package library

import (
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"errors"
	"bufio"
	"strconv"
	"fmt"
	"log"
	"time"
)

import (
	"../../core/methods"
	"../objects"
)


func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			// Open the file
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			// Get stats of the file
			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

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

		if stop {
			log.Println("Downloading completed ....")
			break
		}

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
	req.Header.Set("Authorization", os.ExpandEnv("Token ${API_TOKEN}"))

	// Handle the request
	resp, err := http.DefaultClient.Do(req)
	methods.Fatal_handler(err)

	// If the status code is not 200, then error out
	if resp.StatusCode != http.StatusOK {
		methods.Fatal_handler(errors.New("Upstream ERROR: " + string(resp.StatusCode) + "(" + urlLink + ")"))
	}

	// Close the body once its done
	defer resp.Body.Close()

	// If its to download the software then download it
	if download {

		// The Size of the file
		size := filesize

		// Create th file
		out, err := os.Create(filename)

		// Initalize progress bar
		done := make(chan int64)
		go PrintDownloadPercent(done, filename, int64(size))
		methods.Fatal_handler(err)
		defer out.Close()

		// Start Downloading
		n, err := io.Copy(out, resp.Body)
		methods.Fatal_handler(err)
		done <- n
	}

	// Read the json
	bodyText, err := ioutil.ReadAll(resp.Body)
	methods.Fatal_handler(err)
	return bodyText, err
}

func Prompt_choice() int {

	var choice_entered int
	fmt.Print("\nEnter choose of version that you want to download from the above list (eg.s 1 or 2 etc): ")

	// Start the new scanner to get the user input
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choice_entered, err := strconv.Atoi(input.Text())

		// If user enters a string instead of a integer then ask to re-enter
		if err != nil {
			fmt.Println("Incorrect value: Please choose a integer (eg.s 1 or 2 etc) from the above list")
			return Prompt_choice()
		}

		// If its a valid value move on
		if choice_entered > 0 && choice_entered <= objects.TotalOptions {
			return choice_entered
		} else { // Else ask for re-entering the selection
			fmt.Println("Invalid Choice: The choice you entered is not on the list above, try again.")
			return Prompt_choice()
		}
	}

	return choice_entered
}

