package arguments

type EnvObjects struct {
	Core coreType `yaml:"CORE"`
	Download downloadType `yaml:"DOWNLOAD"`
	Install installType `yaml:"INSTALL"`
}

type coreType struct {
	BaseDir string `yaml:"BASE_DIR"`
	AppName string `yaml:"APPLICATION_NAME"`
}

type downloadType struct {
	ApiToken string `yaml:"API_TOKEN"`
	DownloadDir string `yaml:"DOWNLOAD_DIR"`
}

type installType struct {
	EnvDir string `yaml:"ENV_DIR"`
	UnistallDir string `yaml:"UNINTSALL_DIR"`
	MasterDataDirectory string `yaml:"MASTER_DATA_DIRECTORY"`
	SegmentDataDirectory string `yaml:"SEGMENT_DATA_DIRECTORY"`
	TotalSegments int32 `yaml:"TOTAL_SEGMENT"`
}