package main

import "testing"

func TestReplaceWords(t *testing.T) {
	type testCase struct {
		Message  string
		Expected string
	}

	cases := []testCase{
		{
			Message:  "This is a kerfuffle bro, I can't deal with this sharbert anymore, FORNAX",
			Expected: "This is a **** bro, I can't deal with this **** anymore, ****",
		},
		{
			Message:  "you are a sharbert!",
			Expected: "you are a sharbert!",
		},
	}

	toReplace := []string{"kerfuffle", "sharbert", "fornax"}
	for _, test := range cases {
		actual := replaceWords(test.Message, "****", toReplace)
		if actual != test.Expected {
			t.Errorf("expected: %s\nactual: %s", test.Expected, actual)
		}
	}
}
