package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"rst2md/pkg/config"
)

var whiteSpaceReg = regexp.MustCompile(`\s+`)

// GenerateSlug creates a URL-friendly slug from a string.
func GenerateSlug(input string) string {
	return whiteSpaceReg.ReplaceAllString(strings.ToLower(strings.TrimSpace(input)), "_")
}

// IsDirEmpty checks if a directory is empty.
func IsDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// AskUserOverwrite prompts the user for overwrite confirmation.
func AskUserOverwrite() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Output directory is not empty. Overwrite existing files? (y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" {
			return true, nil
		} else if response == "n" {
			return false, nil
		} else {
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}

// IsUnderline checks if a line is an RST section underline.
func IsUnderline(line string, length int) bool {
	if len(line) != length {
		return false
	}
	validChars := "=~-`#*^\"'+_"
	for _, c := range line {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}
	return true
}

// CopyDir recursively copies a directory from src to dst.
func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, config.DirPermission); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
