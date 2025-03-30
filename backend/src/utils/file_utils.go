package utils

import "os"

func FilesEnsureDirExists(path string) error {
    _, err := os.Stat(path)

    if (os.IsNotExist(err)) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
            return err
		}
    }

    return nil
}
