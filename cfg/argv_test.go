package cfg

import "testing"

func TestGetFirstArgv(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		res, err := getFirstArg([]string{"asdf", "qwert"})
		if err != nil {
			t.Errorf("Expect err to be nil, but got %s\n", err)
		}
		if res != "qwert" {
			t.Errorf("Expect res to be \"qwert\", but got %s\n", res)
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, err := getFirstArg([]string{"asdf"})
		if err != ErrEnoughArgs {
			t.Errorf("Expect err to be ErrEnoughArgs, but got %s\n", err)
		}
	})
}
