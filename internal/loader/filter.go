package loader

// FilterByID searches a slice of APIRequest for the given HexHash
// and returns the first matching object.
func FilterByID(requests []*APIRequest, id string) *APIRequest {
	for _, req := range requests {
		if req.HexHash == id {
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

	var filtered []*APIRequest

	// Check each request only once.
	for _, req := range requests {
		for _, tag := range req.Tags {
			if _, ok := wantedSet[tag]; ok {
				filtered = append(filtered, req)

				break
			}
		}
	}

	return filtered
}
