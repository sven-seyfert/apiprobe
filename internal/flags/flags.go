package flags

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zombiezen.com/go/sqlite"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/logger"
)

type CLIFlags struct {
	Name      *string
	ID        *string
	Tags      *string
	Exclude   *string
	NewID     *bool
	NewFile   *bool
	AddSecret *string
}

// Init defines and parses the CLI flags and returning their values.
func Init() *CLIFlags {
	flag.Usage = func() { //nolint:reassign
		fmt.Fprintf(os.Stderr, config.Version+"\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	nameUsage := "Custom name for this test run (for this execution). Shown in the final notification to help identify the run.\n" +
		"Example: --name \"Environment: PROD\"\n"

	idUsage := "Specify the ten-character hex hash (id) of the request to run.\n" +
		"The hash must match the JSON \"id\" value, in the JSON definition (input) files.\n" +
		"In combination with the --exclude flag, exclude will be prioritized.\n" +
		"Example: --id \"ff00fceb61\"\n"

	tagUsage := "Specify a comma-separated list of tags to select which requests to run.\n" +
		"Tags must match the JSON \"tags\" value, in the JSON definition (input) files.\n" +
		"In combination with the --exclude flag, exclude will be prioritized.\n" +
		"Example: --tags \"reqres, booker\"\n"

	excludeUsage := "Specify a comma-separated list of IDs (hashes) to exclude from the execution.\n" +
		"The IDs must match the JSON \"id\" value, in the JSON definition (input) files.\n" +
		"Example: --exclude \"bb5599abcd, ff00fceb61\"\n"

	newIDUsage := "Generate a new ten-character hex hash (id) for the \n" +
		"JSON \"id\" value, in the JSON definition (input) file.\n" +
		"Example: --new-id\n"

	newFileUsage := "Generate a new JSON definition template file.\n" +
		"Then enter the request values/data and done.\n" +
		"Example: --new-file\n"

	addSecretUsage := "Stores a secret (e.g., API request token, api-key, a bearer token or\n" +
		"other request secrets) in the database and return a placeholder such as \"<secret-b29ff12b50>\".\n" +
		"Use this placeholder in your JSON definition (input) file instead of the actual secret value.\n" +
		"Example: --add-secret \"ThisIsMySecretText\"\n"

	cliFlags := &CLIFlags{
		Name:      flag.String("name", "", nameUsage),
		ID:        flag.String("id", "", idUsage),
		Tags:      flag.String("tags", "", tagUsage),
		Exclude:   flag.String("exclude", "", excludeUsage),
		NewID:     flag.Bool("new-id", false, newIDUsage),
		NewFile:   flag.Bool("new-file", false, newFileUsage),
		AddSecret: flag.String("add-secret", "", addSecretUsage),
	}

	flag.Parse()

	return cliFlags
}

// IsNewID checks whether a new ID should be generated, and if so,
// produces a cryptographically secure hex hash and prints it and
// returns an instruction to exit the program or not.
func IsNewID(isNewID bool) (bool, error) {
	complete := false

	if !isNewID {
		return complete, nil
	}

	hash, err := crypto.HexHash()
	if err != nil {
		logger.Errorf("Failed to generate new ID. Error: %v", err)

		return complete, err
	}

	fmt.Printf(`Use this ID "%s" in your JSON file, key "id".`, hash) //nolint:forbidigo

	complete = true

	return complete, nil
}

// IsNewFile checks if a new file should be created. If true, it generates an ID,
// writes a new template JSON file, and returns true on success. Returns false
// and an error if any step fails.
func IsNewFile(isNewFile bool) (bool, error) {
	complete := false

	if !isNewFile {
		return complete, nil
	}

	hash, err := crypto.HexHash()
	if err != nil {
		logger.Errorf("Failed to generate new ID. Error: %v", err)

		return complete, err
	}

	if err = writeNewTemplateJSONFile(hash); err != nil {
		return complete, err
	}

	complete = true

	return complete, nil
}

// writeNewTemplateJSONFile creates a new JSON definition file (a template)
// with a given ID as content. Returns an error if directory creation
// or file writing fails.
func writeNewTemplateJSONFile(hash string) error {
	content := `[
    {
        "id": "${ID}",
		"isActive": true,
        "isAuthRequest": false,
        "preRequestId": "",
        "request": {
            "description": "...",
            "method": "GET",
            "url": "https://...",
            "endpoint": "/...",
            "basicAuth": "",
            "headers": [],
            "params": [],
            "postBody": {},
			"name": ""
        },
        "testCases": [
            {
                "name": "",
                "paramsData": "",
                "postBodyData": {}
            }
        ],
        "tags": [
            "env-prod"
        ],
        "jq": ""
    }
]`

	const (
		path              = "./data/input/"
		file              = "new-template.json"
		createPermissions = 0o755
		writePermissions  = 0o644
	)

	err := os.MkdirAll(filepath.Dir(path), createPermissions)
	if err != nil {
		logger.Errorf(`Failed to create data/input directory "%s". Error: %v`, file, err)

		return err
	}

	filePath := filepath.Join(path, file)
	content = strings.Replace(content, "${ID}", hash, 1)

	err = os.WriteFile(filePath, []byte(content), writePermissions)
	if err != nil {
		logger.Errorf(`Failed to write file "%s". Error: %v`, filePath, err)

		return err
	}

	return nil
}

// IsAddSecret validates the provided secret string and, if non-empty,
// generates a cryptographically secure hex hash to serve as a placeholder
// and prints it and returns an instruction to exit the program or not.
func IsAddSecret(givenSecret string, conn *sqlite.Conn) (bool, error) {
	complete := false

	if givenSecret == "" {
		return complete, nil
	}

	hash, err := crypto.HexHash()
	if err != nil {
		logger.Errorf("Failed to generate new ID. Error: %v", err)

		return complete, err
	}

	DBValidSecret := crypto.Obfuscate(givenSecret)

	countBefore, err := db.GetTableEntryCount(conn)
	if err != nil {
		return complete, err
	}

	if err = db.InsertSecret(conn, hash, DBValidSecret); err != nil {
		return complete, err
	}

	countAfter, err := db.GetTableEntryCount(conn)
	if err != nil {
		return complete, err
	}

	fmt.Printf("%d ==> %d\n"+ //nolint:forbidigo
		"Use this placeholder \"<secret-%s>\" in your JSON file "+
		"instead of the actual secret value.", countBefore, countAfter, hash)

	complete = true

	return complete, nil
}
