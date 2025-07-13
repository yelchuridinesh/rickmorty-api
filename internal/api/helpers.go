// internal/api/helpers.go
package api

// sanitizeSortBy enforces only "id" or "name", falling back to "id" on anything else.
func sanitizeSortBy(sort string) string {
	switch sort {
	case "id", "name":
		return sort
	default:
		return "id"
	}
}
