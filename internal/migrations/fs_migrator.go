package migrations

import (
	"os"
)

func fsMigrate(pathToMediaDir string) error {
	errs := make([]error, 0, 8)

	errs = append(errs, os.MkdirAll(pathToMediaDir+"/chapters", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/titles", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/teams", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/users", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/titles_on_moderation", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/chapters_on_moderation", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/teams_on_moderation", 0644))
	errs = append(errs, os.MkdirAll(pathToMediaDir+"/users_on_moderation", 0644))

	for i := 0; i < len(errs); i++ {
		if err := errs[i]; err != nil {
			return err
		}
	}

	return nil
}
