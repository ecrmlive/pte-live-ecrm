package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Local struct {
	Dir        string // 物理目录，如 data/uploads
	PublicBase string // URL 前缀，如 /uploads
}

var _ Store = Local{}

func (l Local) Save(scope string, fh *multipart.FileHeader) (publicURL, name string, err error) {
	if l.Dir == "" {
		l.Dir = "data/uploads"
	}
	if l.PublicBase == "" {
		l.PublicBase = "/uploads"
	}
	if fh == nil {
		return "", "", fmt.Errorf("empty file")
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".mp4", ".mov", ".webm":
	default:
		return "", "", fmt.Errorf("unsupported type %s", ext)
	}
	maxSize := int64(10 << 20)
	errMsg := "图片不能超过 10MB"
	if ext == ".mp4" || ext == ".mov" || ext == ".webm" {
		maxSize = 100 << 20
		errMsg = "视频不能超过 100MB"
	}
	if fh.Size > maxSize {
		return "", "", fmt.Errorf("%s", errMsg)
	}
	day := time.Now().Format("20060102")
	relDir := filepath.Join(scope, day)
	absDir := filepath.Join(l.Dir, relDir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return "", "", err
	}
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", "", err
	}
	fname := hex.EncodeToString(b[:]) + ext
	abs := filepath.Join(absDir, fname)
	src, err := fh.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()
	dst, err := os.Create(abs)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", "", err
	}
	publicURL = strings.TrimRight(l.PublicBase, "/") + "/" + filepath.ToSlash(filepath.Join(relDir, fname))
	name = fh.Filename
	if name == "" {
		name = fname
	}
	return publicURL, name, nil
}
