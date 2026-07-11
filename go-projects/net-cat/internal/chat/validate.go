package chat

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
