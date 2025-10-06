package records

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)


func GetRecordNames(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCameraNotFound
		} else {
			return nil, ErrGetRecordsFailed
		}
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".mp4"){
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	return names, nil
}


func GetRecordInfoWithDateFilter(
	name string, 
	dirPath string, 
	useDateFilter bool, 
	filterDate time.Time,
	cameraName string,
	baseUrl string,
	) (Record, error) {

	filePath := filepath.Join(dirPath, name)
	info, err := os.Stat(filePath)
	if err != nil {
		return Record{}, ErrGetFileInfo
	}

	if useDateFilter {
		y, m, d := info.ModTime().Date()
		fy, fm, fd := filterDate.Date()
		if y != fy || m != fm || d != fd {
			return Record{}, ErrDateNotMatch
		}
	}

	url := fmt.Sprintf("%s/%s/%s", strings.TrimRight(baseUrl, "/"), cameraName, name)
	
	return Record{
		Name: name,
		Url: url,
		Timestamp: info.ModTime().UTC(),
		Size: info.Size(),
	}, nil

}
