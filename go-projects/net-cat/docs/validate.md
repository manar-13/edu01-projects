# net-cat — internal/chat/validate.go

## IsValidPort

```go
func IsValidPort(port string) bool {
	if len(port) == 0 {
		return false
	}
	var v int
	for _, c := range port {
		if c < '0' || c > '9' {
			return false
		}
		v = v*10 + int(c-'0')
		if v > 65535 {
			return false
		}
	}
	return v >= 1
}
```
Checks if a port string is a valid port number.
- If the string is empty, return false immediately
- Loops through every character in the string
- If any character is not a digit between `0` and `9`, return false — letters or symbols are not allowed
- Builds the number digit by digit — multiplies the current value by 10 and adds the new digit
- If the number goes above 65535 at any point, return false — that is the maximum valid port number
- At the end, returns true only if the number is 1 or above — port 0 is not allowed

> Valid port numbers are between 1 and 65535. Port 0 is reserved by the operating system. Ports below 1024 are reserved for system services like HTTP (80) and HTTPS (443) — but the code allows them technically. The default port for this project is 8989.

> `int(c-'0')` converts a digit character to its number value. For example the character `'5'` has ASCII value 53 and `'0'` has ASCII value 48, so `53-48 = 5`. This is a common Go trick to convert a digit character to an integer without using `strconv`.
---