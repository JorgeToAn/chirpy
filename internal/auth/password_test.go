package auth

import "testing"

func TestHashPassword(t *testing.T) {
	passwords := [5]string{"123testing", "securepa$$w0rd", "MyPaS$", "admin", "password"}
	hashes := []string{}

	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}
		if hash == password {
			t.Error("expected hash and password to not be the same")
		}
		hashes = append(hashes, hash)
	}

	for i, hash := range hashes {
		match, err := CheckPasswordHash(passwords[i], hash)
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Error("expected match between password and hash")
		}
	}
}
