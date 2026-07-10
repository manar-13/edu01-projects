# go-reloaded — functions.go

## Package and Imports

```go
package functions
```
This file belongs to the `functions` package — the second room where all the tools live.

---

```go
import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "unicode"
)
```
We are borrowing tools we need:
- `fmt` — for converting numbers to text
- `regexp` — for finding patterns in text
- `strconv` — for converting between numbers and text
- `strings` — for working with text
- `unicode` — for checking if a character is a letter or number

---

## ReplaceHex

```go
func ReplaceHex(text string) string {
    words := strings.Fields(text)
    result := []string{}
```
Start the function by splitting the text into individual words and creating an empty list called `result` to build the new text.

---

```go
for i := 0; i < len(words); i++ {
    if words[i] == "(hex)" && i > 0 {
```
Loop through every word. If the current word is `(hex)` and it is not the first word, do something special.

---

```go
        hexVal := words[i-1]
        decimal, err := strconv.ParseInt(hexVal, 16, 64)
        if err == nil {
            result[len(result)-1] = fmt.Sprint(decimal)
            continue
        }
```
- Store the word before `(hex)` in `hexVal`
- Try to convert it from hexadecimal to a decimal number
- If the conversion works, replace the previous word in `result` with the decimal number
- `continue` means skip adding `(hex)` to the result — we don't need it anymore

> `strconv.ParseInt(hexVal, 16, 64)` — the `16` means "read this as a base-16 (hex) number."

---

```go
    result = append(result, words[i])
}

return strings.Join(result, " ")
```
If the word is not `(hex)`, just add it to `result` normally. When the loop is done, join all words back into one text and return it.

---

## ConvBin

```go
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
```
This function works exactly like `ReplaceHex` but for binary numbers.
- It looks for `(bin)` instead of `(hex)`
- It converts the word before it from binary to decimal

> `strconv.ParseInt(binVal, 2, 64)` — the `2` means "read this as a base-2 (binary) number."

---

## CleanWord

```go
func CleanWord(word string) string {
    return strings.TrimFunc(word, func(r rune) bool { 
        return !unicode.IsLetter(r) && !unicode.IsNumber(r) 
    })
}
```
Takes a word and removes anything from the start and end that is not a letter or number.

For example: `"(up)"` becomes `"up"` and `"hello,"` becomes `"hello"`.

> This is used to detect instructions like `(up)`, `(low)`, `(cap)` cleanly without the brackets.

---

## GetPunctu

```go
func GetPunctu(word string) string {
    punct := ""
    for _, ch := range word {
        if strings.ContainsRune(".,!?;:", ch) {
            punct += string(ch)
        }
    }
    return punct
}
```
- Start with an empty text called `punct`
- Go through every character in the word
- If the character is a punctuation mark `.,!?;:` add it to `punct`
- Return all the punctuation found

For example: `"go(up)!"` returns `"!"`.

> This is used to save the punctuation attached to an instruction so it is not lost when we remove the instruction.
---
## ConvFormatInText

```go
func ConvFormatInText(text string) string {
    words := strings.Fields(text)
    result := []string{}
```
Split the text into words and create an empty list called `result` to build the new text.

---

```go
for i := 0; i < len(words); i++ {
    cleaned := CleanWord(words[i])
    punct := GetPunctu(words[i])
```
Loop through every word. For each word:
- `cleaned` — remove the brackets to get the instruction name. Example: `(up)` becomes `up`
- `punct` — save any punctuation attached to the word

---

```go
    if (cleaned == "up" || cleaned == "low" || cleaned == "cap") && i > 0 {
```
If the current word is an instruction (`up`, `low`, or `cap`) and it is not the first word, do something special to the word before it.

---

```go
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
```
Check which instruction it is and apply it to the last word in `result`:
- `up` — make the whole word UPPERCASE
- `low` — make the whole word lowercase
- `cap` — make only the First letter uppercase and the rest lowercase

> `result[len(result)-1]` means "the last word we added to result." This is how we reach back and change the word before the instruction.

---

```go
        if punct != "" {
            result[len(result)-1] += punct
        }

        continue
```
If the instruction had punctuation attached to it, add that punctuation back to the last word. Then `continue` — skip adding the instruction itself to `result`.

---

```go
    result = append(result, words[i])
}

return strings.Join(result, " ")
```
If the word is not an instruction, add it to `result` normally. When done, join everything back into one text and return it.

---

## ConvFormatWithCount

```go
func ConvFormatWithCount(text string) string {
    words := strings.Fields(text)
    result := []string{}
```
Split the text into words and create an empty list called `result`.

---

```go
for i := 0; i < len(words); i++ {
    if strings.HasPrefix(words[i], "(") && strings.HasSuffix(words[i], ",") && i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
```
Loop through every word. Check if the current word looks like the start of a counted instruction. For example `(up,` followed by `3)`.

All four conditions must be true:
- Current word starts with `(`
- Current word ends with `,`
- There is a next word
- The next word ends with `)`

---

```go
        cmd := strings.TrimPrefix(words[i], "(")
        cmd = strings.TrimSuffix(cmd, ",")
        countStr := strings.TrimSuffix(words[i+1], ")")
        count, err := strconv.Atoi(countStr)
```
- Remove `(` from the start and `,` from the end to get the command name. Example: `(up,` becomes `up`
- Remove `)` from the next word to get the number. Example: `3)` becomes `3`
- Convert the number from text to an actual integer

---

```go
        if err != nil {
            result = append(result, words[i], words[i+1])
            i += 1
            continue
        }
```
If the number conversion fails, it is not a real instruction. Add both words to `result` as normal and move on.

---

```go
        start := len(result) - count
        if start < 0 {
            start = 0
        }
```
Calculate which word to start from. If the count is 3, go back 3 words from the end of `result`. If there are not enough words, start from the beginning.

---

```go
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
```
Loop through the words we need to change and apply the instruction to each one. Same logic as `ConvFormatInText` but applied to multiple words.

---

```go
        i += 1
        continue
    }
    result = append(result, words[i])
}

return strings.Join(result, " ")
```
- `i += 1` — skip the next word because we already used it as the count
- `continue` — skip adding the instruction to `result`
- If the word is not an instruction, add it to `result` normally
- When done, join everything back into one text and return it

> The difference between `ConvFormatInText` and `ConvFormatWithCount` is simple. `ConvFormatInText` handles `(up)` which changes only one word. `ConvFormatWithCount` handles `(up, 3)` which changes multiple words.
---
## FixPunctuSpacing

```go
func FixPunctuSpacing(text string) string {
    text = strings.ReplaceAll(text, " ,", ",")
    text = strings.ReplaceAll(text, " .", ".")
    text = strings.ReplaceAll(text, " !", "!")
    text = strings.ReplaceAll(text, " ?", "?")
    text = strings.ReplaceAll(text, " :", ":")
    text = strings.ReplaceAll(text, " ;", ";")
```
Remove any space that appears before a punctuation mark. For example `hello ,` becomes `hello,`.

---

```go
    text = strings.ReplaceAll(text, "...", "§§§")
    text = strings.ReplaceAll(text, "?!", "§§")
```
Temporarily replace `...` and `?!` with symbols `§§§` and `§§`. This protects them from being changed by the next step.

> We use `§` because it never appears in normal English text. It is just a safe placeholder.

---

```go
    re := regexp.MustCompile(`([,.:;!?])([^\s§])`)
    text = re.ReplaceAllString(text, "$1 $2")
```
Find any punctuation mark that is directly touching the next word with no space. Add a space between them. For example `hello,world` becomes `hello, world`.

> `[^\s§]` means "any character that is not a space and not our placeholder symbol." This is why we protected `...` and `?!` first — so they don't get a space inserted inside them.

---

```go
    text = strings.ReplaceAll(text, "§§§", "...")
    text = strings.ReplaceAll(text, "§§", "?!")
    return text
}
```
Put `...` and `?!` back where they belong, then return the fixed text.

---

## FixDoublePunctu

```go
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
```
Find two punctuation marks that have a space between them and remove the space. For example `! ?` becomes `!?` and `. . .` becomes `...`.

The `for` loop keeps repeating until no more changes happen — because there could be more than two punctuation marks in a row.

> When `newText == text` it means nothing changed this time, so we stop the loop with `break`.

---

## FixSingleQuotes

```go
func FixSingleQuotes(text string) string {
    var result []string
    words := strings.Fields(text)
    inQuote := false
```
Split the text into words. Create an empty list called `result` and a flag called `inQuote` to track whether we are currently inside a pair of single quotes.

---

```go
    for _, word := range words {
        if word == "'" {
            if inQuote {
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
```
Loop through every word:
- If the word is `'` and we are not in a quote yet — open the quote, set `inQuote` to true
- If the word is `'` and we are already in a quote — close the quote, set `inQuote` to false
- If the word is anything else — add it to `result` normally

---

```go
    out := strings.Join(result, " ")
    re := regexp.MustCompile(`'\s+(.*?)\s+'`)
    out = re.ReplaceAllString(out, "'$1'")

    return out
}
```
Join all words back into one text. Then find any pattern like `' word '` that still has spaces next to the quotes and remove those spaces. For example `' awesome '` becomes `'awesome'`.

---

## FixAAnGrammar

```go
func FixAAnGrammar(text string) string {
    words := strings.Fields(text)
```
Split the text into words.

---

```go
    exceptions := map[string]bool{
        "unicorn": true, "university": true, "user": true, "useful": true,
        "usual": true, "europe": true, "eulogy": true, "euphemism": true,
        "unanimous": true, "ufo": true,
    }
```
A list of words that start with a vowel letter but are pronounced with a "y" or "w" sound — so they use "a" not "an". For example "a university" not "an university".

---

```go
    forceAn := map[string]bool{
        "hour": true, "honest": true, "honor": true, "hour-long": true,
    }
```
A list of words that start with a consonant letter but are pronounced with a vowel sound — so they use "an" not "a". For example "an hour" not "a hour".

---

```go
    for i := 0; i < len(words)-1; i++ {
        current := strings.ToLower(words[i])
        next := strings.ToLower(words[i+1])
        first := rune(next[0])
```
Loop through every word except the last one. For each word:
- `current` — the current word in lowercase
- `next` — the word after it in lowercase
- `first` — the first letter of the next word

---

```go
        shouldUseAn := (strings.ContainsRune("aeiou", first) && !exceptions[next]) || forceAn[next]
```
Decide if we should use "an":
- The next word starts with a vowel AND is not in the exceptions list
- OR the next word is in the `forceAn` list

---

```go
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
```
- If we should use "an" and the current word is "a" — change it to "an" or "An" depending on the original capitalisation
- If we should use "a" and the current word is "an" — change it to "a" or "A" depending on the original capitalisation
- When done, join all words back into one text and return it
---