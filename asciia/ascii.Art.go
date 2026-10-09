package asciiArt

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var bannerHashes = map[string]string{
	"standard.txt":   "c3ec7584fb7ecfbd739e6b3f6f63fd1fe557d2ae3e24f870730d9cf8b2559e94",
	"shadow.txt":     "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73",
	"thinkertoy.txt": "e3c7a11f41a473d9b0d3bf2132a8f6dabb754bd16efa3897fa835a432d3b9caa",
}

func AsciiArt(text string, style string) (string, error) {
	input := text

	bannerFile, err := getBannerFile(style)
	if err != nil {
		return "", err
	}

	fontMap, err := loadFont(bannerFile)
	if err != nil {
		return "", fmt.Errorf("load banner %q: %w", style, err)
	}

	finalArt := renderText(input, fontMap)

	return finalArt, nil
}

func getBannerFile(style string) (string, error) {
	switch style {
	case "shadow":
		return filepath.Join("style", "shadow.txt"), nil
	case "thinkertoy":
		return filepath.Join("style", "thinkertoy.txt"), nil
	case "standard":
		return filepath.Join("style", "standard.txt"), nil
	default:
		return "", fmt.Errorf("unsupported banner %q", style)
	}
}

func loadFont(fileName string) (map[rune][]string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	bannerName := filepath.Base(fileName)
	expectedHash, ok := bannerHashes[bannerName]
	if !ok {
		return nil, fmt.Errorf("no hash registered for banner %q", bannerName)
	}

	// Git on Windows may check the banners out with CRLF line endings, so
	// hash the LF-normalised content to get the same result on every OS.
	content := strings.ReplaceAll(string(data), "\r", "")
	actualHash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	if actualHash != expectedHash {
		return nil, fmt.Errorf("banner %q was modified", bannerName)
	}

	lines := strings.Split(content, "\n")

	font := make(map[rune][]string)
	ascii := 32

	for i := 0; i+8 < len(lines); i += 9 {
		font[rune(ascii)] = lines[i+1 : i+9]
		ascii++
	}

	return font, nil
}

func renderText(input string, font map[rune][]string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	var finalArt strings.Builder
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if line == "" {
			finalArt.WriteByte('\n')
			continue
		}

		for row := 0; row < 8; row++ {
			for _, char := range line {
				if art, ok := font[char]; ok && len(art) == 8 {
					finalArt.WriteString(art[row])
				} else {
					finalArt.WriteString("        ")
				}
			}
			finalArt.WriteByte('\n')
		}
	}
	return finalArt.String()
}
