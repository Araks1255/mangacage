package titles

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

func FillTitleRequestBody(
	name, englishName, originalName, ageLimit, titleType, publishingStatus, yearOfRelease, description *string,
	authorID *uint,
	genresIDs []uint, tagsIDs []uint, cover []byte,
) (
	res *bytes.Buffer,
	ContentType string,
	err error,
) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if name != nil {
		if err := writer.WriteField("name", *name); err != nil {
			return nil, "", err
		}
	}

	if englishName != nil {
		if err := writer.WriteField("englishName", *englishName); err != nil {
			return nil, "", err
		}
	}

	if originalName != nil {
		if err := writer.WriteField("originalName", *originalName); err != nil {
			return nil, "", err
		}
	}

	if ageLimit != nil {
		if err := writer.WriteField("ageLimit", *ageLimit); err != nil {
			return nil, "", err
		}
	}

	if titleType != nil {
		if err := writer.WriteField("type", *titleType); err != nil {
			return nil, "", err
		}
	}

	if publishingStatus != nil {
		if err := writer.WriteField("publishingStatus", *publishingStatus); err != nil {
			return nil, "", err
		}
	}

	if yearOfRelease != nil {
		if err := writer.WriteField("yearOfRelease", *yearOfRelease); err != nil {
			return nil, "", err
		}
	}

	if authorID != nil {
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", *authorID)); err != nil {
			return nil, "", err
		}
	}

	if description != nil {
		if err := writer.WriteField("description", *description); err != nil {
			return nil, "", err
		}
	}

	if genresIDs != nil && len(genresIDs) != 0 {
		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", fmt.Sprintf("%d", genresIDs[i])); err != nil {
				return nil, "", err
			}
		}
	}

	if tagsIDs != nil && len(tagsIDs) != 0 {
		for i := 0; i < len(tagsIDs); i++ {
			if err := writer.WriteField("tagsIds", fmt.Sprintf("%d", tagsIDs[i])); err != nil {
				return nil, "", err
			}
		}
	}

	if cover != nil {
		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			return nil, "", err
		}
		if _, err := part.Write(cover); err != nil {
			return nil, "", err
		}
	}

	writer.Close()

	return &body, writer.FormDataContentType(), nil
}
