package media

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// LocalPath describes a file resolved on the local filesystem.
type LocalPath struct {
	Path           string
	InputReference string
}

// LocalPathResolver resolves user-provided paths against the current working
// directory and the local ComfyUI input directory.
type LocalPathResolver struct {
	roots    []string
	inputDir string
}

// NewLocalPathResolver creates a resolver. The first root that contains a
// matching relative file wins.
func NewLocalPathResolver(inputDir string, roots ...string) *LocalPathResolver {
	result := &LocalPathResolver{inputDir: filepath.Clean(inputDir)}
	for _, root := range roots {
		if strings.TrimSpace(root) == "" {
			continue
		}
		result.roots = append(result.roots, filepath.Clean(root))
	}
	return result
}

// NewDefaultLocalPathResolver uses the current working directory and the
// standard ComfyUI Desktop input directory for the current user.
func NewDefaultLocalPathResolver() (*LocalPathResolver, error) {
	return NewDefaultLocalPathResolverWithRoots()
}

// NewDefaultLocalPathResolverWithRoots uses the standard ComfyUI Desktop
// input directory and working directory, plus any caller-provided roots.
// Caller-provided roots are searched first so packaged workflow assets can
// take precedence over an unrelated file with the same name in the cwd.
func NewDefaultLocalPathResolverWithRoots(extraRoots ...string) (*LocalPathResolver, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	inputDir := filepath.Join(home, "Documents", "ComfyUI", "input")
	roots := append(extraRoots, workingDir)
	return NewLocalPathResolver(inputDir, roots...), nil
}

// Resolve returns a local regular file when requested resolves to one. A
// missing or non-file value returns found=false so callers can leave normal
// ComfyUI-relative names and literal strings unchanged.
func (r *LocalPathResolver) Resolve(requested string) (resolved LocalPath, found bool, err error) {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return LocalPath{}, false, nil
	}

	candidates := make([]string, 0, len(r.roots)+1)
	if filepath.IsAbs(requested) {
		candidates = append(candidates, requested)
	} else {
		for _, root := range r.roots {
			candidates = append(candidates, filepath.Join(root, requested))
		}
		if r.inputDir != "" {
			candidates = append(candidates, filepath.Join(r.inputDir, requested))
		}
	}

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		candidate = filepath.Clean(candidate)
		key := strings.ToLower(candidate)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		info, statErr := os.Stat(candidate)
		if ignorableStatError(statErr) {
			continue
		}
		if statErr != nil {
			return LocalPath{}, false, statErr
		}
		if !info.Mode().IsRegular() {
			continue
		}

		resolved = LocalPath{Path: candidate}
		if relative, ok := relativePath(r.inputDir, candidate); ok {
			resolved.InputReference = filepath.ToSlash(relative)
		}
		return resolved, true, nil
	}

	return LocalPath{}, false, nil
}

func ignorableStatError(err error) bool {
	if errors.Is(err, os.ErrNotExist) {
		return true
	}

	// Windows reports ordinary literal strings that are too long or contain
	// invalid filename characters as path errors rather than ErrNotExist. They
	// should remain literal values, just like any other non-file input.
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case syscall.Errno(123), // ERROR_INVALID_NAME
			syscall.Errno(206): // ERROR_FILENAME_EXCED_RANGE
			return true
		}
	}
	return false
}

func relativePath(root, candidate string) (string, bool) {
	if strings.TrimSpace(root) == "" {
		return "", false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(candidate))
	if err != nil || relative == "." || relative == "" {
		return "", false
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	if filepath.IsAbs(relative) {
		return "", false
	}
	return relative, true
}
