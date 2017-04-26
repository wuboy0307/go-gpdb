package objects

// Core
var TotalOptions int

// Products
var ProductJsonType ProductObjects

// Release
var ReleaseJsonType ReleaseObjects
var ReleaseURL string
var PivotalProduct string
var ReleaseOutputMap = make(map[string]string)
var ReleaseVersion []string

// All files of the selected version
var VersionJsonType VersionObjects
var DowloadOutputMap = make(map[string]string)
var DownloadOption []string
var DownloadURL string
var ProductFileURL string

// Product file
var ProductFileJsonType ProductFilesObjects
var ProductFileName string
var ProductFileSize int64
