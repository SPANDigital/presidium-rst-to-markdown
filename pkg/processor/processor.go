package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"rst2md/pkg/config"
	"rst2md/pkg/converter"
	"rst2md/pkg/types"
	"rst2md/pkg/utils"

	"gopkg.in/yaml.v2"
)

// Run orchestrates the main workflow of the application.
func Run(cfg config.Config) error {
	// Check for Pandoc
	if err := converter.CheckPandoc(cfg.PandocPath); err != nil {
		return fmt.Errorf("pandoc not found: %w", err)
	}

	// Process directories
	if err := ProcessDirectories(cfg); err != nil {
		return err
	}

	// Process index.rst and parse TOC
	toc, err := ProcessIndexAndGetTOC(cfg)
	if err != nil {
		return err
	}

	// Process external links
	if err := ProcessExternalLinks(cfg.OutputDir, toc); err != nil {
		return err
	}

	// Convert other RST files to Markdown
	if err := ConvertAllRSTFiles(cfg); err != nil {
		return err
	}

	// Process index.rst separately
	if err := ProcessIndexRST(cfg); err != nil {
		return fmt.Errorf("error processing index.rst: %w", err)
	}

	// Create config.yaml
	if err := CreateConfigYAML(cfg.OutputDir, toc); err != nil {
		return err
	}

	// Cleanup intermediate files
	if err := CleanupIntermediateFiles(cfg.OutputDir); err != nil {
		return err
	}

	return nil
}

// ProcessDirectories handles input and output directory setup.
func ProcessDirectories(cfg config.Config) error {
	// Check if input directory exists
	if _, err := os.Stat(cfg.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(cfg.OutputDir, config.DirPermission); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if output directory is empty
	empty, err := utils.IsDirEmpty(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to check if output directory is empty: %w", err)
	}

	if !empty && !cfg.Force {
		overwrite, err := utils.AskUserOverwrite()
		if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}
		if !overwrite {
			return fmt.Errorf("output directory is not empty and user chose not to overwrite existing files")
		}
	}

	// Copy images directory if it exists
	imagesDir := filepath.Join(cfg.InputDir, "images")
	if _, err := os.Stat(imagesDir); err == nil {
		if err := utils.CopyDir(imagesDir, filepath.Join(cfg.OutputDir, "images")); err != nil {
			return fmt.Errorf("failed to copy images directory: %w", err)
		}
	}

	return nil
}

// ProcessIndexAndGetTOC processes index.rst and extracts the table of contents.
func ProcessIndexAndGetTOC(cfg config.Config) ([]types.TOCItem, error) {
	indexPath := filepath.Join(cfg.InputDir, "index.rst")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("index.rst not found in input directory")
	}

	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index.rst: %w", err)
	}

	toc, err := ParseTableOfContents(string(indexContent), cfg.InputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table of contents: %w", err)
	}

	return toc, nil
}

// ParseTableOfContents parses the toctree in index.rst and returns a slice of TOCItem.
func ParseTableOfContents(content string, inputDir string) ([]types.TOCItem, error) {
	var toc []types.TOCItem
	lines := strings.Split(content, "\n")

	// Find the start of the toctree
	tocStart := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == ".. toctree::" {
			tocStart = i + 1
			break
		}
	}

	if tocStart == -1 {
		return nil, fmt.Errorf("table of contents not found in the content")
	}

	// Process toctree entries
	for _, line := range lines[tocStart:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		if strings.Contains(line, "<") && strings.Contains(line, ">") {
			// External link
			parts := strings.SplitN(line, "<", 2)
			name := strings.TrimSpace(parts[0])
			url := strings.TrimSuffix(parts[1], ">")
			toc = append(toc, types.TOCItem{
				ID:             utils.GenerateSlug(name),
				Name:           name,
				IsExternalLink: true,
				URL:            url,
			})
		} else {
			// Regular file
			filePath := filepath.Join(inputDir, line+".rst")
			name, err := GetTopLevelHeading(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to get top-level heading for %s: %w", line, err)
			}
			toc = append(toc, types.TOCItem{
				ID:             line,
				Name:           name,
				IsExternalLink: false,
			})
		}
	}

	return toc, nil
}

// GetTopLevelHeading extracts the top-level heading from an RST file.
func GetTopLevelHeading(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	boldRegex := regexp.MustCompile(`^\*\*(.*)\*\*$`)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "..") || strings.HasPrefix(line, "<") {
			continue
		}
		if i+1 < len(lines) {
			underline := strings.TrimSpace(lines[i+1])
			if utils.IsUnderline(underline, len(line)) {
				// Remove bold formatting if present
				if boldMatches := boldRegex.FindStringSubmatch(line); boldMatches != nil {
					line = boldMatches[1]
				}
				return line, nil
			}
		}
	}

	return "", fmt.Errorf("no top-level heading found in %s", filePath)
}

// ProcessExternalLinks creates directories and _index.md files for external links.
func ProcessExternalLinks(outputDir string, toc []types.TOCItem) error {
	for _, item := range toc {
		if item.IsExternalLink {
			// Create directory for the external link
			dirPath := filepath.Join(outputDir, item.ID)
			if err := os.MkdirAll(dirPath, config.DirPermission); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", item.ID, err)
			}

			// Create _index.md file content
			content := fmt.Sprintf(`---
title: %s
---
%s
`, item.Name, item.URL)

			// Write content to _index.md file
			filePath := filepath.Join(dirPath, "_index.md")
			if err := os.WriteFile(filePath, []byte(content), config.FilePermission); err != nil {
				return fmt.Errorf("failed to create _index.md for %s: %w", item.ID, err)
			}
		}
	}
	return nil
}

// ProcessIndexRST processes the index.rst file separately.
func ProcessIndexRST(cfg config.Config) error {
	inputPath := filepath.Join(cfg.InputDir, "index.rst")
	overviewDir := filepath.Join(cfg.OutputDir, "overview")
	if err := os.MkdirAll(overviewDir, config.DirPermission); err != nil {
		return fmt.Errorf("failed to create overview directory: %w", err)
	}
	outputPath := filepath.Join(overviewDir, "_index.md")

	if err := converter.ConvertRSTToMarkdown(inputPath, outputPath, cfg.PandocPath); err != nil {
		return err
	}

	// Read the converted markdown file
	content, err := os.ReadFile(outputPath)
	if err != nil {
		return err
	}

	// Remove the TOC div
	tocRe := regexp.MustCompile(`(?s)<div class="toctree".*?</div>`)
	content = tocRe.ReplaceAll(content, []byte{})

	// Remove the level 1 heading
	headingRe := regexp.MustCompile(`(?m)^# .+$`)
	content = headingRe.ReplaceAll(content, []byte{})

	// Prepare front matter with fixed "Overview" title
	frontMatter := "---\ntitle: Overview\n---\n\n"

	// Combine front matter and content
	newContent := []byte(frontMatter + string(content))

	// Write the modified content back to the file
	if err := os.WriteFile(outputPath, newContent, config.FilePermission); err != nil {
		return fmt.Errorf("failed to write overview _index.md: %w", err)
	}

	return nil
}

// CreateConfigYAML generates the config.yaml file based on the TOC.
func CreateConfigYAML(outputDir string, toc []types.TOCItem) error {
	var siteConfig types.SiteConfig

	// Add the Overview section
	overviewItem := types.MenuItem{
		Identifier: "overview",
		Name:       "Overview",
		URL:        "/overview/",
		Weight:     10,
	}
	siteConfig.Menu.Main = append(siteConfig.Menu.Main, overviewItem)

	// Add the rest of the TOC items
	for i, entry := range toc {
		item := types.MenuItem{
			Identifier: entry.ID,
			Name:       entry.Name,
			URL:        "/" + entry.ID + "/",
			Weight:     (i + 2) * 10, // Start at 20 and increment by 10
		}
		siteConfig.Menu.Main = append(siteConfig.Menu.Main, item)
	}

	// Marshal the config to YAML
	yamlData, err := yaml.Marshal(&siteConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write the YAML to a file
	configPath := filepath.Join(outputDir, "config.yaml")
	if err := os.WriteFile(configPath, yamlData, config.FilePermission); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

	return nil
}

// CleanupIntermediateFiles removes temporary files from the output directory.
func CleanupIntermediateFiles(outputDir string) error {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") && entry.Name() != "_index.md" {
			if err := os.Remove(filepath.Join(outputDir, entry.Name())); err != nil {
				return fmt.Errorf("failed to remove intermediate file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// ConvertAllRSTFiles converts all RST files in the input directory to Markdown.
func ConvertAllRSTFiles(cfg config.Config) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.MaxParallel)

	errChan := make(chan error, 1)
	done := make(chan bool)

	go func() {
		defer close(done)
		err := filepath.Walk(cfg.InputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// Exclude certain directories like images
				if info.Name() == "images" {
					return filepath.SkipDir
				}
				return nil
			}

			if filepath.Ext(path) == ".rst" {
				relPath, err := filepath.Rel(cfg.InputDir, path)
				if err != nil {
					return err
				}

				// Skip processing index.rst
				if relPath == "index.rst" {
					return nil
				}

				wg.Add(1)
				semaphore <- struct{}{}
				go func(path, relPath string) {
					defer wg.Done()
					defer func() { <-semaphore }()

					outputPath := filepath.Join(cfg.OutputDir, strings.TrimSuffix(relPath, ".rst")+".md")
					if err := os.MkdirAll(filepath.Dir(outputPath), config.DirPermission); err != nil {
						errChan <- fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
						return
					}

					if err := converter.ConvertRSTToMarkdown(path, outputPath, cfg.PandocPath); err != nil {
						errChan <- err
						return
					}

					if err := PostProcessMarkdown(outputPath, cfg.Depth); err != nil {
						errChan <- err
						return
					}
				}(path, relPath)
			}

			return nil
		})
		if err != nil {
			errChan <- err
		}
		wg.Wait()
	}()

	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}

// PostProcessMarkdown splits the Markdown content into sections and creates files accordingly.
func PostProcessMarkdown(filePath string, maxDepth int) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// Split content into sections based on headers
	sections := SplitIntoSections(string(content), maxDepth)

	if len(sections) == 0 {
		return fmt.Errorf("no sections found in %s", filePath)
	}

	// Create directory for the file
	dirName := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	if err := os.MkdirAll(dirName, config.DirPermission); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirName, err)
	}

	// Create _index.md with front matter
	indexContent := fmt.Sprintf("---\ntitle: %s\n---\n\n%s", sections[0].Title, sections[0].Content)
	if err := os.WriteFile(filepath.Join(dirName, "_index.md"), []byte(indexContent), config.FilePermission); err != nil {
		return fmt.Errorf("failed to write _index.md in %s: %w", dirName, err)
	}

	// Create separate files for each section
	for i, section := range sections[1:] {
		fileName := utils.GenerateSlug(section.Title) + ".md"
		sectionContent := fmt.Sprintf("---\ntitle: %s\nweight: %d\n---\n\n%s", section.Title, (i+1)*10, section.Content)
		if err := os.WriteFile(filepath.Join(dirName, fileName), []byte(sectionContent), config.FilePermission); err != nil {
			return fmt.Errorf("failed to write %s in %s: %w", fileName, dirName, err)
		}
	}

	return nil
}

// SplitIntoSections splits the content into sections based on markdown headers.
func SplitIntoSections(content string, maxDepth int) []types.Section {
	var sections []types.Section
	lines := strings.Split(content, "\n")

	var currentSection types.Section
	headingRegex := regexp.MustCompile(`^(#{1,6})\s*(.*)$`)
	boldRegex := regexp.MustCompile(`^\*\*(.*)\*\*$`)

	for _, line := range lines {
		if matches := headingRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			title := matches[2]

			// Remove bold formatting if present
			if boldMatches := boldRegex.FindStringSubmatch(title); boldMatches != nil {
				title = boldMatches[1]
			}

			if level <= maxDepth {
				if currentSection.Title != "" {
					sections = append(sections, currentSection)
				}
				currentSection = types.Section{
					Title:   title,
					Content: "",
				}
			} else {
				// Include the heading in the current section content
				currentSection.Content += line + "\n"
			}
		} else {
			if currentSection.Title != "" {
				currentSection.Content += line + "\n"
			}
		}
	}

	if currentSection.Title != "" {
		sections = append(sections, currentSection)
	}

	return sections
}
