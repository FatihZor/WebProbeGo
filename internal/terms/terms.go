package terms

import (
	"bufio"
	"os"
	"strings"
)

// Load loads terms from a single string and/or a file path.
func Load(single, path string) ([]string, error) {
	var out []string
	if single != "" {
		out = append(out, single)
	}
	if path == "" {
		return out, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	return out, sc.Err()
}
