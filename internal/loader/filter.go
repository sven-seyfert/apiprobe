package loader

import "strings"

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

// FilterByID searches a slice of APIRequest for the given ID
// and returns the first matching object.
func FilterByID(requests []*APIRequest, id string) *APIRequest {
	for _, req := range requests {
		if req.ID == id {
			return req
		}
	}

	return nil
}

// FilterByTags returns all APIRequest objects whose tags intersect
// with the desired tag list.
func FilterByTags(requests []*APIRequest, wantedTags []string) []*APIRequest {
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
