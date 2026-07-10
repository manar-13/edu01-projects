package functions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ReplaceHex(text string) string {
	words := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(words); i++ {
		if words[i] == "(hex)" && i > 0 {
			hexVal := words[i-1]
			decimal, err := strconv.ParseInt(hexVal, 16, 64)
			if err == nil {
				result[len(result)-1] = fmt.Sprint(decimal)
				continue
			}
		}
		result = append(result, words[i])
	}

	return strings.Join(result, " ")
}

func ConvBin(text string) string {
	words := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(words); i++ {
		if words[i] == "(bin)" && i > 0 {
			binVal := words[i-1]
			decimal, err := strconv.ParseInt(binVal, 2, 64)
			if err == nil {
				result[len(result)-1] = fmt.Sprint(decimal)
				continue
			}
		}
		result = append(result, words[i])
	}

	return strings.Join(result, " ")
}

func CleanWord(word string) string {
	return strings.TrimFunc(word, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
}

func GetPunctu(word string) string {
	punct := ""
	for _, ch := range word {
		if strings.ContainsRune(".,!?;:", ch) {
			punct += string(ch)
		}
	}
	return punct
}

func ConvFormatInText(text string) string {
	words := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(words); i++ {
		cleaned := CleanWord(words[i])
		punct := GetPunctu(words[i])
		if (cleaned == "up" || cleaned == "low" || cleaned == "cap") && i > 0 {
			switch cleaned {
			case "up":
				result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			case "low":
				result[len(result)-1] = strings.ToLower(result[len(result)-1])
			case "cap":
				w := result[len(result)-1]
				if len(w) > 0 {
					result[len(result)-1] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
				}
			}

			if punct != "" {
				result[len(result)-1] += punct
			}

			continue
		}

		result = append(result, words[i])
	}

	return strings.Join(result, " ")
}

func ConvFormatWithCount(text string) string {
	words := strings.Fields(text)
	result := []string{}

	for i := 0; i < len(words); i++ {
		if strings.HasPrefix(words[i], "(") && strings.HasSuffix(words[i], ",") && i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
			cmd := strings.TrimPrefix(words[i], "(")
			cmd = strings.TrimSuffix(cmd, ",")
			countStr := strings.TrimSuffix(words[i+1], ")")
			count, err := strconv.Atoi(countStr)
			if err != nil {
				result = append(result, words[i], words[i+1])
				i += 1
				continue
			}

			start := len(result) - count
			if start < 0 {
				start = 0
			}

			for j := start; j < len(result); j++ {
				switch cmd {
				case "up":
					result[j] = strings.ToUpper(result[j])
				case "low":
					result[j] = strings.ToLower(result[j])
				case "cap":
					w := result[j]
					if len(w) > 0 {
						result[j] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
					}
				}
			}

			i += 1
			continue
		}
		result = append(result, words[i])
	}

	return strings.Join(result, " ")
}

func FixPunctuSpacing(text string) string {
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, " !", "!")
	text = strings.ReplaceAll(text, " ?", "?")
	text = strings.ReplaceAll(text, " :", ":")
	text = strings.ReplaceAll(text, " ;", ";")
	text = strings.ReplaceAll(text, "...", "§§§")
	text = strings.ReplaceAll(text, "?!", "§§")
	re := regexp.MustCompile(`([,.:;!?])([^\s§])`)
	text = re.ReplaceAllString(text, "$1 $2")
	text = strings.ReplaceAll(text, "§§§", "...")
	text = strings.ReplaceAll(text, "§§", "?!")
	return text
}

func FixDoublePunctu(text string) string {
	re := regexp.MustCompile(`([!?\.])\s+([!?\.])`)
	for {
		newText := re.ReplaceAllString(text, "$1$2")
		if newText == text {
			break
		}
		text = newText
	}
	return text
}

func FixSingleQuotes(text string) string {
	var result []string
	words := strings.Fields(text)
	inQuote := false

	for _, word := range words {
		if word == "'" {
			if inQuote {
				if len(result) > 0 {
					result[len(result)-1] = strings.TrimRight(result[len(result)-1], " ")
				}
				result = append(result, "'")
				inQuote = false
			} else {
				result = append(result, "'")
				inQuote = true
			}
		} else {
			result = append(result, word)
		}
	}

	out := strings.Join(result, " ")
	re := regexp.MustCompile(`'\s+(.*?)\s+'`)
	out = re.ReplaceAllString(out, "'$1'")

	return out
}

func FixAAnGrammar(text string) string {
	words := strings.Fields(text)

	exceptions := map[string]bool{
		"unicorn": true, "university": true, "user": true, "useful": true,
		"usual": true, "europe": true, "eulogy": true, "euphemism": true,
		"unanimous": true, "ufo": true,
	}

	forceAn := map[string]bool{
		"hour": true, "honest": true, "honor": true, "hour-long": true,
	}

	for i := 0; i < len(words)-1; i++ {
		current := strings.ToLower(words[i])
		next := strings.ToLower(words[i+1])
		first := rune(next[0])

		shouldUseAn := (strings.ContainsRune("aeiou", first) && !exceptions[next]) || forceAn[next]

		if shouldUseAn && current == "a" {
			if words[i] == "A" {
				words[i] = "An"
			} else {
				words[i] = "an"
			}
		} else if !shouldUseAn && current == "an" {
			if words[i] == "An" {
				words[i] = "A"
			} else {
				words[i] = "a"
			}
		}
	}

	return strings.Join(words, " ")
}
