package main

import (
	"bytes"
	"fmt"
	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
	"github.com/mholt/archiver"
	"github.com/ryanuber/columnize"
	"os"
	"regexp"
	"strings"
	"time"
	"strconv"
)

// Function that checks if the string is available on a array.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Check is the value is empty
func IsValueEmpty(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}

// exists returns whether the given file or directory exists or not
func doesFileOrDirExists(path string) (bool, error) {
	Debugf("Checking if the directory %s exists", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Create directory
func CreateDir(path string) {
	// Check if the path or directory exists
	exists, err := doesFileOrDirExists(path)
	if err != nil {
		Fatalf("Failed to check the directory status, the error: %v", err)
	}
	// If not exists then create one
	if !exists {
		Warnf("Directory \"%s\" does not exists, creating one", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			Fatalf("Failed to create the directory, the error: %v", err)
		}
	}
}

// Print the data in tabular format
func printOnScreen(message string, content []string) {

	// Message before the table
	fmt.Printf("\n%s\n\n", message)

	// Print the table format
	result := columnize.SimpleFormat(content)

	// Print the results
	fmt.Println(result + "\n")

}

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
			if err != nil {
				Fatalf("Error in opening the file, err: %v", err)
			}

			// Get stats of the file
			fi, err := file.Stat()
			if err != nil {
				Fatalf("Error in obtaining the stats of the file, err: %v", err)
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

		// Download is completed, time to terminate
		if stop {
			Info("Downloading completed ....")
			Infof("Downloaded file available at: %s", path)
			break
		}

		// Ask to sleep, before repainting the screen.
		time.Sleep(time.Second)
	}
}

// Unzip the binaries.
func unzip(search string) string {
	// Check if we can find the binaries in the directory
	allfiles, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%s.zip", search))

	// Did we find any
	if len(allfiles) > 0 {
		binary := allfiles[0]
		Infof("Found & unzip the binary for the version %s: %s", cmdOptions.Version, binary)
		removeFiles(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("*%s*.bin*", cmdOptions.Version))
		err := archiver.Unarchive(binary, Config.DOWNLOAD.DOWNLOADDIR)
		if err != nil {
			Fatalf("Couldn't unzip the file, err: %v", err)
		}
		Debugf("Unzipped the file %s completed successfully", binary)

		// Get the binary file name
		binFile, _ := FilterDirsGlob(Config.DOWNLOAD.DOWNLOADDIR, fmt.Sprintf("%s.bin", search))
		if len(binFile) > 0 {
			return binFile[0]
		} else {
			Fatalf("No binaries found for the product %s with version %s under directory %s", cmdOptions.Product, cmdOptions.Version, Config.DOWNLOAD.DOWNLOADDIR)
		}
	} else {
		Fatalf("No binary zip found for the product %s with version %s under directory %s", cmdOptions.Product, cmdOptions.Version, Config.DOWNLOAD.DOWNLOADDIR)
	}
	return ""
}

// Extract the contents that we are interested
func contentExtractor(contents []byte, src string, vars []string) bytes.Buffer {

	// Create a parser
	prog, err := parser.ParseProgram([]byte(src), nil)
	if err != nil {
		Fatalf("Failed to parse the program: %s", src)
	}

	// The configuration
	var buf bytes.Buffer
	config := &interp.Config{
		Stdin:  bytes.NewReader([]byte(contents)),
		Vars:   vars,
		Output: &buf,
	}

	// Execute the program
	_, err = interp.ExecProgram(prog, config)
	if err != nil {
		Fatalf("Failure in executing the goawk script: %v", err)
	}

	return buf
}

// Remove blank lines from the contentExtractor
func removeBlanks(s string) string {
	regex, err := regexp.Compile("\n$")
	if err != nil {
		Fatalf("Failure in removing blank lines, err: %v", err)
	}
	s = strings.TrimSpace(regex.ReplaceAllString(s, ""))
	return s
}

// is the port out of range
func outOfRangePort(port string) bool {
	if strToInt(port) > 63000 {
		return true
	}
	return false
}

// string to init
func strToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}