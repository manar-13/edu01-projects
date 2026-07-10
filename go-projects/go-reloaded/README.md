# go-reloaded

## What go-reloaded does?

Imagine you have a text file with messy, unformatted text full of special
instructions like `(up)`, `(hex)`, `(bin)` etc.

Your program **reads that file, fixes everything, and writes the clean result
into a new file.**

Like a **find and replace tool** — but smarter.

---

## The jobs your program does:

| Instruction | What it does |
|---|---|
| `(hex)` | Converts the number before it from hexadecimal to decimal |
| `(bin)` | Converts the number before it from binary to decimal |
| `(up)` | Makes the word before it UPPERCASE |
| `(low)` | Makes the word before it lowercase |
| `(cap)` | Capitalizes the word before it |
| `(up, 3)` | Makes the 3 words before it UPPERCASE |
| Punctuation | Fixes spaces around `. , ! ? : ;` |
| `' word '` | Removes spaces inside single quotes |
| `a` before vowel | Changes "a" to "an" automatically and vice versa |

---

## Author

**Manar Mohamed**