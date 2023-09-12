package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestAuthlibSign(t *testing.T) {
	rsa2048, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	as := NewAuthlibSignWithKey(rsa2048)
	pri, err := as.GetPriKey()
	if err != nil {
		t.Fatal(err)
	}
	pub, err := as.GetPKIXPubKey()
	if err != nil {
		t.Fatal(err)
	}
	sign, err := as.Sign([]byte("xmdhs"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pri)
	fmt.Println(pub)
	fmt.Println(sign)

	hashed := sha1.Sum([]byte("xmdhs"))

	err = rsa.VerifyPKCS1v15(&as.key.PublicKey, crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(sign)))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLittleskinSign(t *testing.T) {
	const pubkey = "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArGcNOOFIqLJSqoE3u0hj\ntOEnOcET3wj9Drss1BE6sBqgPo0bMulOULhqjkc/uH/wyosYnzw3xaazJt87jTHh\nJ8BPMxCeQMoyEdRoS3Jnj1G0Kezj4A2b61PJJM1DpvDAcqQBYsrSdpBJ+52MjoGS\nvJoeQO5XUlJVQm21/HmJnqsPhzcA6HgY71RHYE5xnhpWJiPxLKUPtmt6CNYUQQoS\no2v36XWgMmLBZhAbNOPxYX+1ioxKamjhLO29UhwtgY9U6PWEO7/SBfXzyRPTzhPV\n2nHq7KJqd8IIrltslv6i/4FEM81ivS/mm+PN3hYlIYK6z6Ymii1nrQAplsJ67OGq\nYHtWKOvpfTzOollugsRihkAG4OB6hM0Pr45jjC3TIc7eO7kOgIcGUGUQGuuugDEz\nJ1N9FFWnN/H6P9ukFeg5SmGC5+wmUPZZCtNBLr8o8sI5H7QhK7NgwCaGFoYuiAGL\ngz3k/3YwJ40BbwQayQ2gIqenz+XOFIAlajv+/nyfcDvZH9vGNKP9lVcHXUT5YRnS\nZSHo5lwvVrYUrqEAbh/zDz8QMEyiujWvUkPhZs9fh6fimUGxtm8mFIPCtPJVXjeY\nwD3Lvt3aIB1JHdUTJR3eEc4eIaTKMwMPyJRzVn5zKsitaZz3nn/cOA/wZC9oqyEU\nmc9h6ZMRTRUEE4TtaJyg9lMCAwEAAQ==\n-----END PUBLIC KEY-----\n"

	const signed = `QcQVFjzn1iIts/HV1PakfAjI+2b9mYnX7SwY5Kqlbqt96aa3x+3c8H7J6PDl2Hb/M3q1rm5LHgda51EMbsyhm1zbXX7Xakkp5aoAm+jsZo9MZuL7KRW2IVwo44sWaUH37NM72nIyn4YBxja1Wg5pKspxDdg2iF8BTR/27vIXLejS+ZAFvo+pERJL+ye+2CQcpknbuZN9a8ns5ylSr71kJ0xo6lVpAgOO5M7RiD7SXrICSLGItopIankq2Ra64pR/GR86AfHgfbY8D0jPmtc0p6/KH6wbQtNSgxRnf6A0VFSoUvj+G7kwgOb3Fh2hnfnDsWlLl3MxVSqNrpoL/y0p9OjytHyVZmoTiCVE3dX5qehebo5wmNvHYV1sS0zCsjqomu9JipKmt5uQtujez6mymXtwL1oVsKpkaxZYb4FoDRnuMF8rWElilMiMyZNDcACr0fDTtUC/w+u2IMTiSMisNL/XELzPoZdVwf3z+1Eklyn3kIhlvIH4mnj4r8C1FtEMYPMVFmXIFNFv10qh32vHTSGoA3ZNduPd72rERkfC4wZYAGfkWQ0fT2f2BwRoeJm4ixDrMoHTCi5MgXC7t9Cij/tyuoaYx80DIxc8sA5itmpLV4o1LqO3DC8n8QWDcN3sRtxEP08ToyR5Q35HjyImPMcLkf/wEAepAujKScQmvO0=`

	b, _ := pem.Decode([]byte(pubkey))
	if b == nil {
		t.FailNow()
	}
	pub, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	hashed := sha1.Sum([]byte("skin,cape"))
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(signed)))
	if err != nil {
		t.Fatal(err)
	}

}

func TestAuthlibNew(t *testing.T) {
	rsa2048, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	as := NewAuthlibSignWithKey(rsa2048)

	_, err = NewAuthlibSign([]byte(lo.Must1(as.GetPriKey())))
	if err != nil {
		t.Fatal(err)
	}
}
