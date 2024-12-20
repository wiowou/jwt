//go:build go1.4
// +build go1.4

package alg_test

import (
	"crypto/rsa"
	"os"
	"strings"
	"testing"

	"github.com/wiowou/jwt/pkg/alg"
	"github.com/wiowou/jwt/pkg/pemc"
)

var rsaPSSTestData = []struct {
	name        string
	tokenString string
	alg         string
	claims      map[string]interface{}
	valid       bool
}{
	{
		"Basic PS256",
		"eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.PPG4xyDVY8ffp4CcxofNmsTDXsrVG2npdQuibLhJbv4ClyPTUtR5giNSvuxo03kB6I8VXVr0Y9X7UxhJVEoJOmULAwRWaUsDnIewQa101cVhMa6iR8X37kfFoiZ6NkS-c7henVkkQWu2HtotkEtQvN5hFlk8IevXXPmvZlhQhwzB1sGzGYnoi1zOfuL98d3BIjUjtlwii5w6gYG2AEEzp7HnHCsb3jIwUPdq86Oe6hIFjtBwduIK90ca4UqzARpcfwxHwVLMpatKask00AgGVI0ysdk0BLMjmLutquD03XbThHScC2C2_Pp4cHWgMzvbgLU2RYYZcZRKr46QeNgz9w",
		"PS256",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic PS384",
		"eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.w7-qqgj97gK4fJsq_DCqdYQiylJjzWONvD0qWWWhqEOFk2P1eDULPnqHRnjgTXoO4HAw4YIWCsZPet7nR3Xxq4ZhMqvKW8b7KlfRTb9cH8zqFvzMmybQ4jv2hKc3bXYqVow3AoR7hN_CWXI3Dv6Kd2X5xhtxRHI6IL39oTVDUQ74LACe-9t4c3QRPuj6Pq1H4FAT2E2kW_0KOc6EQhCLWEhm2Z2__OZskDC8AiPpP8Kv4k2vB7l0IKQu8Pr4RcNBlqJdq8dA5D3hk5TLxP8V5nG1Ib80MOMMqoS3FQvSLyolFX-R_jZ3-zfq6Ebsqr0yEb0AH2CfsECF7935Pa0FKQ",
		"PS384",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Basic PS512",
		"eyJhbGciOiJQUzUxMiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.GX1HWGzFaJevuSLavqqFYaW8_TpvcjQ8KfC5fXiSDzSiT9UD9nB_ikSmDNyDILNdtjZLSvVKfXxZJqCfefxAtiozEDDdJthZ-F0uO4SPFHlGiXszvKeodh7BuTWRI2wL9-ZO4mFa8nq3GMeQAfo9cx11i7nfN8n2YNQ9SHGovG7_T_AvaMZB_jT6jkDHpwGR9mz7x1sycckEo6teLdHRnH_ZdlHlxqknmyTu8Odr5Xh0sJFOL8BepWbbvIIn-P161rRHHiDWFv6nhlHwZnVzjx7HQrWSGb6-s2cdLie9QL_8XaMcUpjLkfOMKkDOfHo6AvpL7Jbwi83Z2ZTHjJWB-A",
		"PS512",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"basic PS256 invalid: foo => bar",
		"eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.PPG4xyDVY8ffp4CcxofNmsTDXsrVG2npdQuibLhJbv4ClyPTUtR5giNSvuxo03kB6I8VXVr0Y9X7UxhJVEoJOmULAwRWaUsDnIewQa101cVhMa6iR8X37kfFoiZ6NkS-c7henVkkQWu2HtotkEtQvN5hFlk8IevXXPmvZlhQhwzB1sGzGYnoi1zOfuL98d3BIjUjtlwii5w6gYG2AEEzp7HnHCsb3jIwUPdq86Oe6hIFjtBwduIK90ca4UqzARpcfwxHwVLMpatKask00AgGVI0ysdk0BLMjmLutquD03XbThHScC2C2_Pp4cHWgMzvbgLU2RYYZcZRKr46QeNgz9W",
		"PS256",
		map[string]interface{}{"foo": "bar"},
		false,
	},
}

func TestRSAPSSVerify(t *testing.T) {
	var err error

	key, _ := os.ReadFile("../00_files/sample_key.pub")
	var rsaPSSKey *rsa.PublicKey
	if rsaPSSKey, err = pemc.ToRSAPublicKey(key); err != nil {
		t.Errorf("Unable to parse RSA public key: %v", err)
	}

	for _, data := range rsaPSSTestData {
		parts := strings.Split(data.tokenString, ".")

		method, err := alg.GetAlg(data.alg)
		if err != nil {
			t.Error(err)
		}
		err = method.Verify(strings.Join(parts[0:2], "."), decodeSegment(t, parts[2]), rsaPSSKey)
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying key: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestRSAPSSSign(t *testing.T) {
	var err error

	key, _ := os.ReadFile("../00_files/sample_key")
	var rsaPSSKey *rsa.PrivateKey
	if rsaPSSKey, err = pemc.ToRSAPrivateKey(key); err != nil {
		t.Errorf("Unable to parse RSA private key: %v", err)
	}

	for _, data := range rsaPSSTestData {
		if !data.valid {
			continue
		}
		parts := strings.Split(data.tokenString, ".")
		method, err := alg.GetAlg(data.alg)
		if err != nil {
			t.Error(err)
		}
		sig, err := method.Sign(strings.Join(parts[0:2], "."), rsaPSSKey)
		if err != nil {
			t.Errorf("[%v] Error signing token: %v", data.name, err)
		}

		ssig := encodeSegment(sig)
		if ssig == parts[2] {
			t.Errorf("[%v] Signatures shouldn't match\nnew:\n%v\noriginal:\n%v", data.name, ssig, parts[2])
		}
	}
}

// func TestRSAPSSSaltLengthCompatibility(t *testing.T) {
// 	// Fails token verify, if salt length is auto.
// 	ps256SaltLengthEqualsHash := &alg.SigningMethodRSAPSS{
// 		SigningMethodRSA: alg.PS256.SigningMethodRSA,
// 		Options: &rsa.PSSOptions{
// 			SaltLength: rsa.PSSSaltLengthEqualsHash,
// 		},
// 	}

// 	// Behaves as before https://github.com/dgrijalva/jwt-go/issues/285 fix.
// 	ps256SaltLengthAuto := &alg.SigningMethodRSAPSS{
// 		SigningMethodRSA: alg.SigningMethodPS256.SigningMethodRSA,
// 		Options: &rsa.PSSOptions{
// 			SaltLength: rsa.PSSSaltLengthAuto,
// 		},
// 	}
// 	if !verify(t, alg.PS256, makeToken(ps256SaltLengthEqualsHash)) {
// 		t.Error("PS256 should accept salt length that is defined in RFC")
// 	}
// 	if !verify(t, ps256SaltLengthEqualsHash, makeToken(alg.PS256)) {
// 		t.Error("Sign by PS256 should have salt length that is defined in RFC")
// 	}
// 	if !verify(t, alg.PS256, makeToken(ps256SaltLengthAuto)) {
// 		t.Error("PS256 should accept auto salt length to be compatible with previous versions")
// 	}
// 	if !verify(t, ps256SaltLengthAuto, makeToken(alg.PS256)) {
// 		t.Error("Sign by PS256 should be accepted by previous versions")
// 	}
// 	if verify(t, ps256SaltLengthEqualsHash, makeToken(ps256SaltLengthAuto)) {
// 		t.Error("Auto salt length should be not accepted, when RFC salt length is required")
// 	}
// }

// func makeToken(method alg.ISigningAlgorithm) string {
// 	token := alg.NewWithClaims(method, alg.RegisteredClaims{
// 		Issuer:   "example",
// 		IssuedAt: alg.NewNumericDate(time.Now()),
// 	})
// 	privateKey := loadRSAPrivateKeyFromDisk("../00_files/sample_key")
// 	signed, err := token.SignedString(privateKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return signed
// }

// func verify(t *testing.T, signingMethod alg.ISigningAlgorithm, token string) bool {
// 	segments := strings.Split(token, ".")
// 	err := signingMethod.Verify(strings.Join(segments[:2], "."), decodeSegment(t, segments[2]), loadRSAPublicKeyFromDisk("../00_files/sample_key.pub"))
// 	return err == nil
// }
