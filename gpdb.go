package main

import (
	"./src/download"
	"./src/argParser"
)


func main() {

	argParser.ArgParser()
	download.Download()
}