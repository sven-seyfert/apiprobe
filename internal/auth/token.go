package auth

// TokenStore maintains a map of request IDs to API tokens.
// Each key is a 10 character hex hash, and each value
// is the corresponding token.
type TokenStore struct {
	data map[string]string
}

// NewTokenStore initializes and returns a new TokenStore.
func NewTokenStore() *TokenStore {
	return &TokenStore{
		data: make(map[string]string),
	}
}

// Add inserts a token for the given id if it does not already exist.
// Returns true if the token was added, false if the id already exists.
func (t *TokenStore) Add(id, token string) bool {
	if _, exists := t.data[id]; exists {
		return false
	}

	t.data[id] = token

	return true
}

// Get retrieves the token for the given id. Returns the token and
// true if found, or "" and false otherwise.
func (t *TokenStore) Get(id string) (string, bool) {
	if t, found := t.data[id]; found {
		return t, true
	}

	return "", false
}
