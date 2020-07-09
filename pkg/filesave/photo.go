package filesave

import (
	"os"
	"path/filepath"
)

var photoDir = "/app/photo/jq"

type Photo interface {
	Save(name string, src []byte) error
	Load(name string) (*os.File, error)
}

type PhotoDefault struct {
}

//func init() {
//	server.OnAppStart(Init)
//}
//func Init() {
//	if config.Config.FilePath.PhotoDir != "" {
//		photoDir = config.Config.FilePath.PhotoDir
//	}
//	_, err := os.Stat(photoDir)
//	if os.IsNotExist(err) {
//		os.MkdirAll(photoDir, 666)
//	}
//}

func (pd PhotoDefault) Save(name string, src []byte) error {
	dst, _ := os.Create(filepath.Join(photoDir, name))
	defer dst.Close()
	_, err := dst.Write(src)
	if err != nil {
		return err
	}
	return nil
}

func (pd PhotoDefault) Load(name string) (*os.File, error) {
	src, err := os.Open(filepath.Join(photoDir, name))
	if err != nil {
		return nil, err
	}
	return src, nil
}

func PhotoSave(name string, src []byte, handle Photo) error {
	return handle.Save(name, src)
}

func PhotoLoad(name string, handle Photo) (*os.File, error) {
	return handle.Load(name)
}
