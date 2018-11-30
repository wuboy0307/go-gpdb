package main

import "fmt"

var (
	programName = "datalab"
	programVersion = "0.1"
	configFileName = fmt.Sprintf(".%s.config", programName)
)

func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
