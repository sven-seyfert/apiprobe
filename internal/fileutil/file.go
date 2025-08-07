package fileutil

import (
	"os"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// EnsureFileExists ensures that the target file and its directory exist,
// creating them if necessary.
func EnsureFileExists(outputPath string) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err = createOutputDir(outputPath); err != nil {
			return err
		}

		if err = WriteOutputFile(outputPath, nil); err != nil {
			return err
		}
	}

	return nil
}

// WriteOutputFile writes byte content to the specified file
// with defined permissions.
func WriteOutputFile(outputPath string, output []byte) error {
	const permissions = 0o644

	if err := os.WriteFile(outputPath, output, permissions); err != nil {
		logger.Errorf(`Failed to write file "%s". Error: %v`, outputPath, err)

		return err
	}

	return nil
}
