package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/sven-seyfert/apiprobe/internal/config"
	"github.com/sven-seyfert/apiprobe/internal/crypto"
	"github.com/sven-seyfert/apiprobe/internal/db"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"zombiezen.com/go/sqlite"
)

type CLIFlags struct {
	ID        *string
	Tags      *string
	Exclude   *string
	NewID     *bool
	AddSecret *string
}

// Init defines and parses the CLI flags and returning their values.
func Init() *CLIFlags {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, config.Version+"\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

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

	addSecretUsage := "Stores a secret (e.g., API request token, api-key, a bearer token or\n" +
		"other request secrets) in the database and return a placeholder such as \"<secret-b29ff12b50>\".\n" +
		"Use this placeholder in your JSON definition (input) file instead of the actual secret value.\n" +
		"Example: --add-secret \"ThisIsMySecretText\"\n"

	cliFlags := &CLIFlags{
		ID:        flag.String("id", "", idUsage),
		Tags:      flag.String("tags", "", tagUsage),
		Exclude:   flag.String("exclude", "", excludeUsage),
		NewID:     flag.Bool("new-id", false, newIDUsage),
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
