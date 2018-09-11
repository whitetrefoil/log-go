package server

import "testing"

func TestMatcher(t *testing.T) {
	t.Parallel()
	t.Run("Get IDs", func(t *testing.T) {
		if matcher == nil {
			t.Fatal("matcher failed to compile")
		}
		m1 := matcher.FindStringSubmatch("/api/harvesters/1/2")
		if m1 == nil || m1[1] != "1" || m1[2] != "2" {
			t.Fatal("Wrong match:", m1)
		}
		m2 := matcher.FindStringSubmatch("/api/harvesters/1/2/3")
		if m2 != nil {
			t.Fatal("Wrong match:", m2)
		}
		m3 := matcher.FindStringSubmatch("/api/harvesters/1")
		if m3 != nil {
			t.Fatal("Wrong match:", m3)
		}
	})
}
