package internal

import "os"

func FileExists(p string) bool {
	s, _ := os.Stat(p)
	return s != nil
}
