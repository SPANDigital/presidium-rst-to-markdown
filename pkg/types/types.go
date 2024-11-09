// pkg/types/types.go
package types

// TOCItem represents an item in the table of contents.
type TOCItem struct {
	ID             string
	Name           string
	IsExternalLink bool
	URL            string // Only used if IsExternalLink is true
}

// MenuItem represents a menu item in the site configuration.
type MenuItem struct {
	Identifier string `yaml:"identifier"`
	Name       string `yaml:"name"`
	URL        string `yaml:"url,omitempty"`
	Weight     int    `yaml:"weight"`
}

// SiteConfig represents the structure of the site's configuration file.
type SiteConfig struct {
	Menu struct {
		Main []MenuItem `yaml:"main"`
	} `yaml:"menu"`
}

// Section represents a section in the markdown content.
type Section struct {
	Title   string
	Content string
}
