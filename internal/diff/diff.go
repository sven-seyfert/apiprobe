package diff

import (
	"crypto/sha256"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/fileutil"
	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// HasFileContentChanged compares the SHA256 checksum of the given output
// bytes against the current contents of outputPath. If they differ,
// writes the new content to file and returns true;
// otherwise logs 'No change' and returns false.
func HasFileContentChanged(output []byte, outputPath string) (bool, error) {
	err := fileutil.EnsureFileExists(outputPath)
	if err != nil {
		return false, err
	}

	newHash := sha256.Sum256(output)

	var prevHash [32]byte

	existing, err := os.ReadFile(outputPath)
	if err != nil {
		logger.Errorf(`Failed to read file "%s"`, outputPath)

		return false, err
	}

	prevHash = sha256.Sum256(existing)

	if newHash == prevHash {
		logger.Infof(`No change for "%s"`, outputPath)

		return false, nil
	}

	logger.Infof(`Detected change (diff) in "%s"`, outputPath)

	if err := fileutil.WriteOutputFile(outputPath, output); err != nil {
		return true, err
	}

	return true, nil
}
