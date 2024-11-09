package converter

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ConvertRSTToMarkdown converts an RST file to Markdown using Pandoc.
func ConvertRSTToMarkdown(inputPath, outputPath, pandocPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, pandocPath, "-f", "rst", "-t", "gfm", inputPath, "-o", outputPath)
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error converting %s: %v\n%s", inputPath, err, stderr)
	}
	return nil
}

// CheckPandoc verifies if Pandoc is available in the system.
func CheckPandoc(pandocPath string) error {
	cmd := exec.Command(pandocPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pandoc not found: %w", err)
	}
	return nil
}
