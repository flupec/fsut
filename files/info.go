package files

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func GetFileProperties(filePath string) (FileProperties, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return FileProperties{}, err
	}
	return NewFileProperties(fileInfo), nil
}

func NewFileProperties(info os.FileInfo) FileProperties {
	perm := info.Mode().String()
	return FileProperties{
		Perm:           perm,
		Name:           info.Name(),
		Type:           NewFileType(perm),
		Size:           NewFileSize(info.Size()),
		LastChangeTime: info.ModTime(),
	}
}

func NewFileSize(byteSize int64) FileSize {
	if byteSize == 0 {
		return FileSize{
			SizeUnitAmount: 0,
			Factor: Byte,
		}
	}

	idx := len(orderedFactors) - 1
	log.Infof("%s", orderedFactors[idx])
	for idx != 0 || byteSize / int64(orderedFactors[idx]) == 0 {
		idx--
		log.Infof("%s", orderedFactors[idx])
	}
	factor := orderedFactors[idx]
	howManyFits := float64(byteSize / int64(factor))
	remain := float64(byteSize % int64(factor)) / float64(factor)

	return FileSize{
		SizeUnitAmount: howManyFits + remain,
		Factor: factor,
	}
}

type FileProperties struct {
	Perm           string
	Name           string
	Type           FileType
	Size           FileSize
	LastChangeTime time.Time
}

type FileSize struct {
	SizeUnitAmount float64
	Factor         SizeFactor
}

type FileType int

const (
	Directory FileType = iota
	Link
	File
	Pipe
	Socket
)

func NewFileType(perm string) FileType {
	switch perm[0: 1] {
	case "d":
		return Directory
	case "l":
		return Link
	case "s":
		return Socket
	case "-":
		return File
	case "p":
		return Pipe
	default:
		return File
	}
}

func (f FileType) String() string {
	switch f {
	case Directory:
		return "DIR"
	case Link:
		return "LINK"
	case File:
		return "FILE"
	case Socket:
		return "SOCKET"
	case Pipe:
		return "PIPE"
	default:
		return "FILE"
	}
}

type SizeFactor int64

const (
	Byte             = 1
	KByte SizeFactor = 1024
	MByte            = KByte * 1024
	GByte            = MByte * 1024
)

func (s SizeFactor) String() string {
	switch s {
	case Byte:
		return "B"
	case KByte:
		return "KB"
	case MByte:
		return "MB"
	case GByte:
		return "GB"
	default:
		return "size more than 1024 GB"
	}
}

var orderedFactors = [4]SizeFactor{Byte, KByte, MByte, GByte}