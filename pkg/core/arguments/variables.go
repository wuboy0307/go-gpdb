package arguments

// Parser
var ArgOption string

// Enironment
var EnvYAML EnvObjects
var TempDir string
var DownloadDir string
var EnvFileDir string
var UninstallDir string
var FutureRefDir string

// Download
var AcceptedDownloadProduct = []string{"gpdb", "gpcc", "gpextras"}
var RequestedDownloadProduct string
var RequestedDownloadVersion string
var InstallAfterDownload bool = false

// Install
var AcceptedInstallProduct = []string{"gpdb", "gpcc"}
var RequestedInstallProduct string
var RequestedInstallVersion string
var RequestedCCInstallVersion string

// Core
var YesOrNo = map[string]string{"y":"y", "ye":"y", "yes":"y", "n":"n", "no":"n" }
