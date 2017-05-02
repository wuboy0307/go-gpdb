package methods

import (
	"os"
	"io"
	"fmt"
)

func CreateFile(path string) error {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil { return err}
		defer file.Close()
	}

	return nil
}

func WriteFile(path string, args []string) error {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil { return err}
	defer file.Close()

	// write some text to file
	for _, k := range args {
		text := string(k)
		_, err = file.WriteString(text + "\n")
		if err != nil { return err}
	}

	// save changes
	err = file.Sync()
	if err != nil { return err}

	return nil
}

func ReadFile(path string) error {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil { return err}
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if err != nil { return err}
		}
		if n == 0 {
			break
		}
	}
	fmt.Println(string(text))
	if err != nil { return err}

	return nil
}

func DeleteFile(path string) error {

	// delete file
	var err = os.Remove(path)
	if err != nil { return err}

	return nil
}
