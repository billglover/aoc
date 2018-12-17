package battlefield

import "testing"

func TestClanString(t *testing.T) {
	names := map[Clan]string{
		Elf:    "Elf",
		Goblin: "Goblin",
	}

	for k, v := range names {
		u := Unit{MemberOf: k}

		if got, want := u.MemberOf.String(), v; got != want {
			t.Errorf("unexpected clan name: got %s, want %s", got, want)
		}
	}
}
