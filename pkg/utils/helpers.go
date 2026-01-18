package utils

import "strconv"

// DefaultPage returns default page number
func DefaultPage(page string) int {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		return 1
	}
	return p
}

// DefaultLimit returns default limit number
func DefaultLimit(limit string) int {
	l, err := strconv.Atoi(limit)
	if err != nil || l < 1 {
		return 10
	}
	if l > 100 {
		return 100
	}
	return l
}

// GetSkipNumber calculates offset for pagination
func GetSkipNumber(page, limit int) int {
	return (page - 1) * limit
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// UintPtr returns a pointer to a uint
func UintPtr(i uint) *uint {
	return &i
}
