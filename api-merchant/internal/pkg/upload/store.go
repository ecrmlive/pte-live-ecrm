package upload

import "mime/multipart"

// Store 素材上传存储（本地或腾讯云 COS）。
type Store interface {
	Save(scope string, fh *multipart.FileHeader) (publicURL, name string, err error)
}
