package pwd

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// CheckAgainstPwned hashes the provided password
// and compares the password hash to the response queried from the PwnedPasswords API
// to determine if the password is within the safety threshold.
//
// Returns nil if the password is within the threshold.
//
// See: https://haveibeenpwned.com/API/v3#PwnedPasswords
//
// See: https://www.troyhunt.com/ive-just-launched-pwned-passwords-version-2/
func CheckAgainstPwned(appUserAgent string, pwd string, threshold int64) error {
	h := sha1.New()

	_, err := h.Write([]byte(pwd))
	if err != nil {
		return err
	}

	hashed := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	hashPrefix := hashed[:5]

	reqUrl := fmt.Sprintf("https://api.pwnedpasswords.com/range/%v", hashPrefix)
	log.Print(reqUrl)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", appUserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	for result := range strings.SplitSeq(string(body), "\n") {
		parts := strings.Split(strings.TrimSpace(result), ":")

		hashSuffix := parts[0]
		count, err := strconv.ParseInt(parts[1], 10, 64)

		if err != nil {
			return err
		}

		resultHash := hashPrefix + hashSuffix

		if resultHash == hashed {
			log.Printf("found pwned %v:%d", resultHash, count)

			if count >= threshold {
				return fmt.Errorf("insecure password")
			}

			return nil
		}
	}

	return nil
}

func GenerateSalt(saltLen int) []byte {
	saltBytes := make([]byte, saltLen)
	rand.Read(saltBytes)

	return saltBytes
}

func HashPwd(pwd []byte, salt []byte, p HashParams) []byte {
	return argon2.IDKey(pwd, salt, p.Time, p.Memory, p.Threads, p.KeyLen)
}
