package auth

import (
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/loader"
	"github.com/sven-seyfert/apiprobe/internal/logger"
	"github.com/sven-seyfert/apiprobe/internal/util"
)

// RepaceAuthTokenPlaceholderInRequestHeader replaces the <auth-token> placeholder
// in request headers with the corresponding token from the token store, if available.
// Returns nothing.
func RepaceAuthTokenPlaceholderInRequestHeader(req *loader.APIRequest, tokenStore *TokenStore) {
	const headerReplacementIndicator = "<auth-token>"

	lookupID := req.PreRequestID

	for idx, header := range req.Request.Headers {
		if !strings.Contains(header, headerReplacementIndicator) {
			continue
		}

		if token, found := tokenStore.Get(lookupID); found {
			lastTokenChars := token[util.Max(0, len(token)-12):] //nolint:mnd

			logger.Debugf(`Token "...%s" found for auth request "%s".`, lastTokenChars, lookupID)

			req.Request.Headers[idx] = strings.ReplaceAll(header, headerReplacementIndicator, token)

			break
		}

		logger.Warnf(`No token found for auth request "%s".`, lookupID)
	}
}

// AddAuthTokenToTokenStore attempts to add the token to the provided token store
// using the request ID as the key. Returns nothing.
func AddAuthTokenToTokenStore(result []byte, tokenStore *TokenStore, req *loader.APIRequest) {
	token := util.TrimQuotes(string(result))
	lastTokenChars := token[util.Max(0, len(token)-12):] //nolint:mnd

	if added := tokenStore.Add(req.ID, token); added {
		logger.Debugf(`Token "...%s" for auth request "%s" added to token store.`, lastTokenChars, req.ID)
	} else {
		logger.Warnf(`Token "...%s" for auth request "%s" already exists in token store.`, lastTokenChars, req.ID)
	}
}
