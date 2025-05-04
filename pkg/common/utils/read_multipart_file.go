package utils

import (
	"io"
	"mime/multipart"
)

func ReadMultipartFile(fileHeader *multipart.FileHeader, limit int64) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := io.LimitReader(file, limit)

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
