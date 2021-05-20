package tools

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"

	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + MD5V(q.URL) + q.GetQrCodeExt()
	if CheckNotExist(src) == true {
		return false
	}

	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := MD5V(q.URL) + q.GetQrCodeExt()
	src := path + name
	if CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		defer f.Close()
		if err != nil {
			return "", "", err
		}

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

type Merge struct {
	BgFilePath     string
	BgFileName     string
	QrCodeFilePath string
	QrCodeFileName string
	MergeFilePath  string
	Rect           *Rect
	Pt             *Pt
}

func NewMerge(bgFilePath, bgFileName, qrCodeFilePath, qrCodeFileName, mergeFilePath string, scaleX, scaleY int) *Merge {
	return &Merge{
		BgFilePath:     bgFilePath,
		BgFileName:     bgFileName,
		QrCodeFilePath: qrCodeFilePath,
		QrCodeFileName: qrCodeFileName,
		MergeFilePath:  mergeFilePath,
		Pt: &Pt{
			X: scaleX,
			Y: scaleY,
		},
	}
}

func (m *Merge) CheckMergedImage(mergeFile string) bool {
	if CheckNotExist(m.MergeFilePath+mergeFile) == true {
		return false
	}

	return true
}

func (m *Merge) OpenMergedImage(mergeFile string) (*os.File, error) {
	f, err := file.MustOpen(mergeFile, m.MergeFilePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (m *Merge) Generate(mergeFile string) (string, string, error) {

	if !m.CheckMergedImage(m.MergeFilePath + mergeFile) {

		mergedF, err := m.OpenMergedImage(mergeFile)
		defer mergedF.Close()
		if err != nil {
			return "", "", err
		}

		bgF, err := MustOpen(m.BgFileName, m.BgFilePath)
		defer bgF.Close()
		if err != nil {
			return "", "", err
		}

		qrF, err := MustOpen(m.QrCodeFileName, m.QrCodeFilePath)
		defer qrF.Close()
		if err != nil {
			return "", "", err
		}

		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		//jpg := image.NewRGBA(image.Rect(m.Rect.X0, m.Rect.Y0, m.Rect.X1, m.Rect.Y1))
		// rect := bgImage.Bounds()
		jpg := image.NewRGBA(bgImage.Bounds())
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(m.Pt.X, m.Pt.Y)), draw.Over)

		jpeg.Encode(mergedF, jpg, nil)
	}

	return mergeFile, m.MergeFilePath, nil
}

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
