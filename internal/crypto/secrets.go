package crypto

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"zombiezen.com/go/sqlite"
)

// HandleSecrets iterates over each APIRequest in filteredRequests, finds all
// placeholders '<secret-<hash>>' in PostBody, BasicAuth, Params, Headers and TestCases,
// retrieves the real secret from the database, deobfuscates it, and replaces the
// placeholder. Returns an error immediately if any DB lookup fails.
func HandleSecrets(filteredRequests []*loader.APIRequest, conn *sqlite.Conn) ([]*loader.APIRequest, error) {
	for _, req := range filteredRequests {
		newBody, err := replaceSecretInString(req.PostBody, conn)
		if err != nil {
			return nil, err
		}

		req.PostBody = newBody

		newAuth, err := replaceSecretInString(req.BasicAuth, conn)
		if err != nil {
			return nil, err
		}

		req.BasicAuth = newAuth

		if err := replaceSecretInSlice(req.Params, conn); err != nil {
			return nil, err
		}

		if err := replaceSecretInSlice(req.Headers, conn); err != nil {
			return nil, err
		}

		if err := replaceSecretInSlice(req.TestCases, conn); err != nil {
			return nil, err
		}
	}

	return filteredRequests, nil
}

// replaceSecretInString searches a single string for '<secret-<hash>>'
// patterns. For each found hash, it retrieves the secret from the database,
// deobfuscates it, and replaces the placeholder in the string.
// Returns an error if DB lookup fails.
func replaceSecretInString(str string, conn *sqlite.Conn) (string, error) {
	const secretPrefix = "<secret-"

	if !strings.Contains(str, secretPrefix) {
		return str, nil
	}

	secretHash := extractSecretHash(str)

	if secretHash == "" {
		logger.Warnf("No valid secret hash found in string: %s", str)

		return str, nil
	}

	secret, err := db.SelectHash(conn, secretHash)
	if err != nil {
		logger.Debugf(`Failed to retrieve secret for hash "%s": %v`, secretHash, err)

		return "", err
	}

	if secret != "" {
		from := fmt.Sprintf("%s%s>", secretPrefix, secretHash)
		to := Deobfuscate(secret)

		return strings.ReplaceAll(str, from, to), nil
	}

	logger.Warnf(`Secret value "%s" not found`, secretHash)

	return str, nil
}

// replaceSecretInSlice iterates over a slice of strings, calls replaceSecretInString
// on each element, and updates the slice in-place.
// Returns the first error encountered, if any.
func replaceSecretInSlice(reqSlice []string, conn *sqlite.Conn) error {
	for idx, val := range reqSlice {
		newVal, err := replaceSecretInString(val, conn)
		if err != nil {
			return err
		}

		reqSlice[idx] = newVal
	}

	return nil
}

// extractSecretHash uses a precompiled regex to extract the hash from
// a '<secret-<hash>>' placeholder. Returns the hash without angle brackets
// or prefix or an empty string if no match is found.
func extractSecretHash(input string) string {
	pattern := regexp.MustCompile(`<secret-([^>]+)>`)
	matches := pattern.FindStringSubmatch(input)

	if len(matches) < 2 { //nolint:mnd
		logger.Warnf(`"No secret hash found: "%s"`, input)

		return ""
	}

	return matches[1]
}
