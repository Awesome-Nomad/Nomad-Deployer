package git

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// GetGitProjectDir return directory name with prefix `baseDir`.
// eg: GetGitProjectDir("./project", "https://github.com/liemle3893/abc") => ./projects/github.com/liemle3893/abc
func GetGitProjectDir(baseDir string, gitURL *url.URL) string {
	gitPath := strings.ReplaceAll(gitURL.RequestURI(), filepath.Ext(gitURL.RequestURI()), "")
	return fmt.Sprintf("%s/%s%s", baseDir, gitURL.Host, gitPath)
}
