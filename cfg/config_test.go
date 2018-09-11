package cfg

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("show version", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		err := Get([]string{os.Args[0], "version"}, &struct{}{})

		if err != ErrPrintVersion {
			t.Errorf("Wrong error %v", err)
		}
	})
}
