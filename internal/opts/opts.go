package opts

import (
	"io/fs"
	"strings"
)

type Opts struct {
	ignoreDirs     []string
	ignoreDotfiles bool
	allowedExts    []string
}

func NewOpts(ignoreDirs []string, ignoreDotfiles bool, allowedExts []string) *Opts {
	return &Opts{
		ignoreDirs,
		ignoreDotfiles,
		allowedExts,
	}
}

func (o *Opts) IsDirAllowed(dir fs.FileInfo) bool {
	if o.ignoreDotfiles && strings.HasPrefix(dir.Name(), ".") {
		return false
	} else if len(o.ignoreDirs) > 0 {
		for _, d := range o.ignoreDirs {
			if dir.Name() == d {
				return false
			}
		}
	}
	return true
}

func (o *Opts) IsFileAllowed(fi fs.FileInfo) bool {

	for _, ext := range o.allowedExts {
		if strings.HasSuffix(fi.Name(), ext) {
			return true
		}
	}
	return false
}
