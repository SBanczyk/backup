package backend

type Backend interface {
	DownloadBackupFilesFile() (string, error)
	UploadFile(filePath string, cloudName string) error
	DownloadFile(cloudName string, filePath string) error
	RemoveFile(cloudName string) error
}
