package utils

// ValidatePagination validates and normalizes pagination parameters
// Returns normalized page and limit values
func ValidatePagination(page, limit int) (int, int) {
	// Validate page
	if page < 1 {
		page = 1
	}

	// Validate limit
	if limit < 1 {
		limit = 10 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit to prevent abuse
	}

	return page, limit
}
