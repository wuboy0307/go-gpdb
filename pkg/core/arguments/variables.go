package arguments

// Parser
var ArgOption string

// Enironment
var EnvYAML EnvObjects
var TempDir string
var DownloadDir string
var EnvFileDir string
var UninstallDir string

// Download
var AcceptedDownloadProduct = []string{"gpdb", "gpcc", "gpextras"}
var RequestedDownloadProduct string
var RequestedDownloadVersion string
var InstallAfterDownload bool

// Install
var AcceptedInstallProduct = []string{"gpdb", "gpcc"}
var RequestedInstallProduct string
var RequestedInstallVersion string
