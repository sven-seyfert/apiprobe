package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// BuildOutputFilePath computes the output file path for a
// given APIRequest and optional test case index, by inserting
// '-test-case-XX' into the JSON file name and nesting under
// './data/output'.
func BuildOutputFilePath(req *loader.APIRequest, testCaseIndex *int) string {
	outputDir := "./data/output"

	fileExt := filepath.Ext(req.JSONFilePath)
	file := req.JSONFilePath

	if testCaseIndex != nil {
		file = strings.Replace(file, fileExt, fmt.Sprintf("-test-case-%02d%s", *testCaseIndex+1, fileExt), 1)
	} else {
		file = strings.Replace(file, fileExt, fmt.Sprintf("-test-case-%02d%s", 0, fileExt), 1)
	}

	return filepath.Join(outputDir, file)
}

// createOutputDir ensures that the parent directory for the given
// output path exists. If necessary, it creates all missing directories.
func createOutputDir(outputPath string) error {
	const permissions = 0o755

	if err := os.MkdirAll(filepath.Dir(outputPath), permissions); err != nil {
		logger.Errorf(`Failed to create output directory "%s". Error: %v"`, outputPath, err)

		return err
	}

	return nil
}
