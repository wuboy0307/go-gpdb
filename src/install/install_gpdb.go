package install


import (
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
)


func Install() {

	// If the install is called from download command the set default values
	if !methods.IsValueEmpty(arguments.RequestedDownloadVersion) {
		arguments.RequestedInstallVersion = arguments.RequestedDownloadVersion
	}

	// Unzip the binaries, if its file is zipped
	UnzipBinary()

}
