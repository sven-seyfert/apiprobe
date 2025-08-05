package loader

import (
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
