package loader

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sven-seyfert/apiprobe/internal/logger"
)

// ExcludeRequestsByID returns a filtered slice of APIRequest, excluding any
// requests whose IDs are listed in the comma-separated excludeIDs string.
func ExcludeRequestsByID(requests []*APIRequest, excludeIDs string) []*APIRequest {
	if excludeIDs == "" {
		return requests
	}

	idList := strings.Split(excludeIDs, ",")
	excludeSet := make(map[string]struct{})

	for _, id := range idList {
		id = strings.TrimSpace(id)
		if id != "" {
			excludeSet[id] = struct{}{}
		}
	}

	var filteredRequests []*APIRequest

	for _, req := range requests {
		if _, found := excludeSet[req.ID]; !found {
			filteredRequests = append(filteredRequests, req)
		}
	}

	return filteredRequests
}

// FilterRequests filters the given slice of APIRequest by the '--id'
// and '--tags' flags. It returns a slice of matching requests and a
// boolean flag that is true if no requests matched the filters.
func FilterRequests(requests []*APIRequest, id string, tags string) ([]*APIRequest, bool) { //nolint:varnamelen
	if len(requests) == 0 {
		logger.Warnf(`No requests found.`)

		return requests, true
	}

	// Filter requests by ID.
	if id != "" {
		if req := filterByID(requests, id); req != nil {
			return []*APIRequest{req}, false
		}

		logger.Warnf(`No request with id (hex hash) "%s" found.`, id)

		return requests, true
	}

	// Or filter requests by tags.
	if tags != "" {
		tagsList := strings.Split(tags, ",")
		wantedTags := make([]string, 0, len(tagsList))

		for _, tag := range tagsList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				wantedTags = append(wantedTags, tag)
			}
		}

		filteredRequests := filterByTags(requests, wantedTags)
		if len(filteredRequests) > 0 {
			return filteredRequests, false
		}

		logger.Warnf(`No requests found for tags "%s".`, tags)

		return requests, true
	}

	// Or use the fallback (return all requests).
	return requests, false
}

// filterByID searches a slice of APIRequest for the given ID
// and returns the first matching object.
func filterByID(requests []*APIRequest, id string) *APIRequest {
	for _, req := range requests {
		if req.ID == id {
			return req
		}
	}

	return nil
}

// filterByTags returns all APIRequest objects whose tags intersect
// with the desired tag list.
func filterByTags(requests []*APIRequest, wantedTags []string) []*APIRequest {
	// Build a set (map) for O|1 lookup of desired tags.
	wantedSet := make(map[string]struct{}, len(wantedTags))
	for _, w := range wantedTags {
		wantedSet[w] = struct{}{}
	}

	var filteredRequests []*APIRequest

	// Check each request only once.
	for _, req := range requests {
		for _, tag := range req.Tags {
			if _, ok := wantedSet[tag]; ok {
				filteredRequests = append(filteredRequests, req)

				break
			}
		}
	}

	return filteredRequests
}

// MergePreRequests constructs a merged requests list in which, for each
// filtered request having a PreRequestID, the corresponding loaded request
// is prepended before the filtered requests. It returns the gathered/merged
// APIRequest list without duplicates.
func MergePreRequests(loadedRequests []*APIRequest, filteredRequests []*APIRequest) ([]*APIRequest, error) {
	lookupMap := make(map[string]*APIRequest, len(loadedRequests))

	for _, loadedReq := range loadedRequests {
		lookupMap[loadedReq.ID] = loadedReq
	}

	const tenCharHexHashPattern = `^[a-fA-F0-9]{10}$`

	hexPattern := regexp.MustCompile(tenCharHexHashPattern)
	requestsList := make([]*APIRequest, 0, len(loadedRequests)+len(filteredRequests))

	// Handle possible pre-requests by PreRequestID.
	for _, filteredReq := range filteredRequests {
		preID := filteredReq.PreRequestID
		if preID == "" {
			continue
		}

		if !hexPattern.MatchString(preID) {
			logger.Errorf(`PreRequestID "%s" has invalid format (not the expected ten character hex hash format).`, preID)

			return nil, errors.New("invalid format error")
		}

		prev, found := lookupMap[preID]
		if !found {
			logger.Errorf(`PreRequestID "%s" not found in loadedRequests.`, preID)

			return nil, errors.New("not found error")
		}

		requestsList = append(requestsList, prev)
	}

	// Append all filtered requests (to be behind the pre-requests).
	requestsList = append(requestsList, filteredRequests...)

	return removeDuplicates(requestsList), nil
}

// removeDuplicates returns a new slice of APIRequest pointers with
// duplicates removed, keeping only the first occurrence of each
// request based on its ID.
func removeDuplicates(requestsList []*APIRequest) []*APIRequest {
	seen := make(map[string]bool, len(requestsList))
	unique := make([]*APIRequest, 0, len(requestsList))

	for _, req := range requestsList {
		if !seen[req.ID] {
			seen[req.ID] = true
			unique = append(unique, req)
		}
	}

	return unique
}
