package btdaemon

func IsASCIIPrintable(data []byte) bool {
	for _, b := range data {
		if b < 0x20 || b > 0x7E {
			return false
		}
	}
	return true
}
