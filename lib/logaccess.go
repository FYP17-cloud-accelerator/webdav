package lib

import (
	"strings"
	"time"
)

type LogAccess struct {
	Username   string
	FileName   string
	Path       string
	FullPath   string
	Extension  string
	AccessTime time.Time
	ModTime    time.Time
	Size       int64
}

func NewLogAccess(username string, full_path string, mod_time time.Time, size int64) *LogAccess {
	res := strings.Split(full_path, "/")
	filename := res[len(res)-1]
	path := strings.Replace(full_path, filename, "", 1)
	res = strings.Split(filename, ".")
	extension := res[len(res)-1]
	return &LogAccess{
		Username:   username,
		FileName:   filename,
		Extension:  extension,
		AccessTime: time.Now(),
		ModTime:    mod_time,
		Path:       path,
		FullPath:   full_path,
		Size:       size,
	}
}
