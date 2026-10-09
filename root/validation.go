package root

func validText(text string) bool {
	for _, char := range text {
		if char == '\r' || char == '\n' {
			continue
		}
		if char < 32 || char > 126 {
			return false
		}
	}
	return true
}

func validStyle(style string) bool {
	switch style {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}