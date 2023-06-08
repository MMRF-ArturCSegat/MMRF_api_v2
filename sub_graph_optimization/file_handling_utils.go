package sub_graph_optimization

import (
    "os"
    "time"
)


// GetLastModificationTime returns the last modification time of the given file.
func GetLastModificationTime(file *os.File) (time.Time, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}


func GetFileSize(file *os.File) int64 {
	stat, err := file.Stat()
	if err != nil {
		return 0
	}
	return stat.Size()
}


// GetFilePath returns the path of the given *os.File
func GetFilePath(file *os.File) string {
	fileInfo, err := file.Stat()
	if err != nil {
		return ""
	}
	return fileInfo.Name()
}
