package methods

import "log"

// Fatal Handler that terminates the program when its called
func Fatal_handler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


// Error Handler that return the error when called.
func Error_handler(err error) {
	if err != nil {
		log.Println(err)
	}
}