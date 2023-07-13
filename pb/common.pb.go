package pb

import "mime/multipart"

type ID struct {
	ID string
}

type Empty struct {
}

type File struct {
	FileHeader *multipart.FileHeader `json:"file_header" form:"file_header"`
}
