package art

import "os"

func WriteToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
