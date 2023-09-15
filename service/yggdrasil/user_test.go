package yggdrasil

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/samber/lo"
)

const (
	publicKeySignatureV2Str = "sB+BdmutEFqkIoIDOSPOrle8tCQl+lv8pCKqviivreZ0Et8zNgUtFagR9jvWiGfhSbOtQXf1ooDZDQx/TO58KA9k1gSoKoB3IWvxUooW5fiAftH+VvFluO+c9bulWiZ/Lk28EOk/moaGtovgYzMtt6F4BIJYB/iZLSkuRp/MiJqsUJcDN5DWViizWSjBsuAkiV+6+/T1QW9Z6IKrOvlLfJqhm/HSXgTuxS4CbzzqxGWZL49TDdeZFs0q1O+NgK2k5zLeL1JjfFdS6hu5X13pa9ygmTLrFIwuml9rT65mlQV5KSXGkO7uajKVpKBfqoM8gZ42wc20kyItDiSVIw4afoEpdr5AtlOL4VIK+Qnphty3YS18sbZzIukIMxMN/s/i90QA9xuMy/U0a8RwOq0haFJ5OYES9EFIYJjqa7uZJ1+riSJvXWQZwbh4l09lQ5p4/TBCDIqa3cp+YH2kVM8JCxhsLNmkkoXlV9LSJQ2zyAe1ItkrIi/HYx2qLQ56CIel5qM4j60+kPpoGfmvgeEiHDERloH9lekHHUeRXawXV6RE1wwf41Qi9+3UycJQWGOi03qprE42YjsRMi3BUToEQqK1QOPJo8x2F/YfGVRJpLMBOS1g25uqw/kIOBAzMWfGRtQ5iKuPzFGhLx+rJjp58pyr0939KMvQtq/71ADQllM="
	publicKeySignatureStr   = "kxjk7pdzEyQ5Pj+OOza2hmMFFZ1RHSqAHv1DcHaF+um5ZLRBvuwRTmbmvwuq4IFafdPvwaGqohNF7U9EXv+/v2JrEkZQWRXTfhzx+WSpogfDFtIOe20QFKOfZ1kBgPYaqCjSx67idvgaCMxafjWZDmZEMfJLv5TL1rO2LHnyeMIlFdVRhVbRiwa8f7RSPcNhLL2I0zHl40EG7lryIW08bctd4Ksgd0eDYw/tJCfY79NV1PurKE+9YCr2fK1ErhSfirkdfmXptEQ9iEZ+YRwkgYZRI6bQiwwpMqMDFLNOL5STxuEZz6ynygRI2WL+n09FYt066Eci3RC9lRVt3VUrRO2TWx3stwzI8zktrIkplYUC1l/ECtt+RBgN4Fc7QU4RX1ZBirUPDykK2gObS5OtowTi0o3tTPOnl0Jy/a1vr+Nqc1TPlDomE6xOBDtCD61sMuPedNM8IoBl71OphR91t4oxDDe5mZMZBCc27J6apjy3rfzL4Xpi9g72047QKYrQEYNWpF9ddkX/Ed8roC/e7Zaa77ARhSDfBbKs9gRvILqFZ8eprMBGzn55K4QI23l8G1X7IS4squaWkRNi+LXRq8M936DqYRIZ9bloJa5W+XJPGn29dWGxPlto2ELzGLI/69vemqNiKLmoXocBRByVAi7TU/Tqr9ZiJU1RBjsdh2w="
	publikKey               = "-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvUP6vXby++vuXGONEve5grn9cu2Zv3KQ\nXimz2kdIlN8DZ3u4ZUHkS6b1vWpOhS/IWE5TIp2iA7HOVH6WOTK3uViBhMXaMZ5eFuqASfA/jPF8\nxaGmHiH4dBmGzbO7B+DRvL38PYLCTMzgwxp1qEuW7VVmcd6Glr0HG9q3YDB0rb9gjPhvj6seT6sf\ntoYC0wHX29qUe4Y89bgSsF9WDf0lQ+qdjLwlcP/Aqyc+VaCYLV67LiRDjQc8q2JmPoy4KEwxOYuD\nEfHSfHzu48xn2mfvVMJXRU09VDLvvXcteEt7KLF6a1QizNx51AXwUq8xre42dgQkPF4Dnv4JxysD\niMVHzwIDAQAB\n-----END RSA PUBLIC KEY-----\n"
	expiresAt               = "2023-09-17T03:06:05.251232611Z"
)

func Test_publicKeySignature(t *testing.T) {
	b, err := os.ReadFile("yggdrasil_session_pubkey.der")
	if err != nil {
		t.Fatal(err)
	}
	pub := lo.Must(x509.ParsePKIXPublicKey(b)).(*rsa.PublicKey)
	pubByte, _ := pem.Decode([]byte(publikKey))
	if pubByte == nil {
		t.FailNow()
	}

	userPub := lo.Must(x509.ParsePKIXPublicKey(pubByte.Bytes)).(*rsa.PublicKey)
	timeUinx := lo.Must(time.Parse(time.RFC3339Nano, expiresAt)).UnixMilli()

	t.Run("publicKeySignatureV2", func(t *testing.T) {
		signByte, err := publicKeySignatureV2(userPub, "9f51573a5ec545828c2b09f7f08497b1", timeUinx)
		if err != nil {
			t.Fatal(err)
		}

		hashed := sha1.Sum(signByte)

		err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(publicKeySignatureV2Str)))
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("publicKeySignature", func(t *testing.T) {
		b := publicKeySignature(publikKey, timeUinx)

		hashed := sha1.Sum(b)

		err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(publicKeySignatureStr)))
		if err != nil {
			t.Fatal(err)
		}
	})
}
