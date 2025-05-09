package utils

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(title string, maxLength int) (string, error) {
	if maxLength <= 0 {
		maxLength = 100 // Default max length
	}

	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace special characters and multiple dashes
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove duplicate dashes
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim dashes from beginning and end
	slug = strings.Trim(slug, "-")

	// Set maximum length for slug
	if len(slug) > maxLength {
		slug = slug[:maxLength]
	}

	return slug, nil
}

// EnsureUniqueSlug makes sure a slug is unique in the specified table and column
// If the slug already exists, it appends a number to make it unique
func EnsureUniqueSlug(db *gorm.DB, baseSlug, tableName, columnName string, excludeID ...uint) (string, error) {
	slug := baseSlug

	for i := 1; ; i++ {
		// Check if slug exists
		var count int64
		query := db.Table(tableName).Where(columnName+" = ?", slug)

		// If excluding a specific ID (for updates)
		if len(excludeID) > 0 && excludeID[0] > 0 {
			query = query.Where("id != ?", excludeID[0])
		}

		if err := query.Count(&count).Error; err != nil {
			return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
		}

		if count == 0 {
			break // Slug is unique, we can use it
		}

		// If slug exists, append a number and try again
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	return slug, nil
}
