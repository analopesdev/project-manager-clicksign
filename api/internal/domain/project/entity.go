package project

import (
	"net/url"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

const (
	SortName      = "name"
	SortStartDate = "start_date"
	SortEndDate   = "end_date"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Client    string    `json:"client"`
	StartDate string    `json:"start_date"` // YYYY-MM-DD
	EndDate   string    `json:"end_date"`   // YYYY-MM-DD
	CoverURL  string    `json:"cover_url"`
	Favorite  bool      `json:"favorite"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidateProjectInput(name string, client string, startDate string, endDate string, coverURL string) error {
	name = strings.TrimSpace(name)
	client = strings.TrimSpace(client)
	coverURL = strings.TrimSpace(coverURL)

	validationErrors := &ValidationError{}

	if name == "" {
		validationErrors.add("name", "name is required")
	} else if len([]rune(name)) > 120 {
		validationErrors.add("name", "name must have at most 120 characters")
	}

	if client == "" {
		validationErrors.add("client", "client is required")
	} else if len([]rune(client)) > 120 {
		validationErrors.add("client", "client must have at most 120 characters")
	}

	start, startErr := time.Parse(DateLayout, startDate)
	if startErr != nil {
		validationErrors.add("start_date", "start_date must be a valid date in YYYY-MM-DD format")
	}
	end, endErr := time.Parse(DateLayout, endDate)
	if endErr != nil {
		validationErrors.add("end_date", "end_date must be a valid date in YYYY-MM-DD format")
	}
	if startErr == nil && endErr == nil && end.Before(start) {
		validationErrors.add("end_date", "end_date must not be before start_date")
	}

	if coverURL != "" {
		u, err := url.ParseRequestURI(coverURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			validationErrors.add("cover_url", "cover_url must be a valid http or https URL")
		}
	}

	return validationErrors.orNil()
}

// ValidateSort accepts an empty sort (the default) or one of the supported orderings.
func ValidateSort(sort string) error {
	switch sort {
	case "", SortName, SortStartDate, SortEndDate:
		return nil
	}
	validationErrors := &ValidationError{}
	validationErrors.add("sort", "sort must be one of name, start_date, end_date")
	return validationErrors
}
