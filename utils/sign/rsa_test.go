package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestAuthlibSign(t *testing.T) {
	rsa2048, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}
	as := NewAuthlibSignWithKey(rsa2048)
	pri, err := as.GetPriKey()
	require.Nil(t, err)
	require.NotEmpty(t, pri)
	pub, err := as.GetPKIXPubKey()
	require.Nil(t, err)
	require.NotEmpty(t, pub)

	sign, err := as.Sign([]byte("xmdhs"))
	require.Nil(t, err)

	hashed := sha1.Sum([]byte("xmdhs"))

	err = rsa.VerifyPKCS1v15(&as.key.PublicKey, crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(sign)))
	require.Nil(t, err)

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

func TestMojangSign(t *testing.T) {
	b, err := os.ReadFile("../../service/yggdrasil/yggdrasil_session_pubkey.der")
	if err != nil {
		t.Fatal(err)
	}
	pub := lo.Must(x509.ParsePKIXPublicKey(b)).(*rsa.PublicKey)

	const value = "ewogICJ0aW1lc3RhbXAiIDogMTY5NDc5NDE2NzU3NSwKICAicHJvZmlsZUlkIiA6ICI5ZjUxNTczYTVlYzU0NTgyOGMyYjA5ZjdmMDg0OTdiMSIsCiAgInByb2ZpbGVOYW1lIiA6ICJ4bWRocyIsCiAgInNpZ25hdHVyZVJlcXVpcmVkIiA6IHRydWUsCiAgInRleHR1cmVzIiA6IHsKICAgICJTS0lOIiA6IHsKICAgICAgInVybCIgOiAiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9iNDc4YmQxY2FlODJhNGQwN2M4NWU4Y2FjYTc3NTA3MWQ5ZjYxZWYyYTkzNDQzOWU2NTFhZTcwYzAzZDBmMTVkIgogICAgfSwKICAgICJDQVBFIiA6IHsKICAgICAgInVybCIgOiAiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS8yMzQwYzBlMDNkZDI0YTExYjE1YThiMzNjMmE3ZTllMzJhYmIyMDUxYjI0ODFkMGJhN2RlZmQ2MzVjYTdhOTMzIgogICAgfQogIH0KfQ=="
	const signature = "h1dqL2jie5uRBt1xVOHACsn3YrdOWidSYvXq+eQXJIv9cyQAlns9ASE20sRHaQCL6Gil0nfnyGDaHWbXJ3Vi9+bCyvPNRgFFKLgNDfcWEgEhPducWH30Pl3zDyAXWUZQbT2ecmHDTfkzb1UR3SGEnVxmJpV3++RDjgotRy6bWLE1Hx8U8OZdpTwhZb9Y3m/otPBpYgHl0SQAEGQvr1dn4/SAMM13GNXqynlWzZ2X94I9DvO8oLEn6VtIIgw+0kHKmQfdepgeLMDDOOUaXskGo1liV7efEMTGlTLZrgEzgo/rNsWn1O98Wc+3mjshLPP2PIJqSfGpXvPeE6Z2wCKJdfguQEYNRomP/8gCxEzfXG1eg/XE7FQBEVi6Eoath3aYopTqcDwKL4v0f520JNcPtTfXfIZiYGpJ9JiZFqjL8q51Y4SUIcDN7vX2/OzdPiJ5xI1MEK1AsLDUaSWAzR6ZwjmOoFv0m68U44c5GDnnvt4kN+oWM0jUAzGAU1QXutrGiVee/1jryEgXJVM43x+D9ZYJvWBXDFqoCLY9hLvRKtt3ohv/aQTDFwTLZzhM82kQ74RdqwtRCtrNsCIqjLlwxrQTY8xYWjLWMvYI82x40+Zk4aEmfc6PEp8gwlLV/gTtYlzsR7uJC+lpmN9Is/LiC9bMj4iHjP3Dk4ykwA/s/5k="

	hashed := sha1.Sum([]byte(value))
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], lo.Must(base64.StdEncoding.DecodeString(signature)))
	if err != nil {
		t.Fatal(err)
	}

}

func TestMojangPubKey(t *testing.T) {
	b, err := os.ReadFile("../../service/yggdrasil/yggdrasil_session_pubkey.der")
	if err != nil {
		t.Fatal(err)
	}
	pub := lo.Must(x509.ParsePKIXPublicKey(b)).(*rsa.PublicKey)
	derBytes, err := x509.MarshalPKIXPublicKey(pub)
	require.Nil(t, err)

	const key = `MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAylB4B6m5lz7jwrcFz6Fd/fnfUhcvlxsTSn5kIK/2aGG1C3kMy4VjhwlxF6BFUSnfxhNswPjh3ZitkBxEAFY25uzkJFRwHwVA9mdwjashXILtR6OqdLXXFVyUPIURLOSWqGNBtb08EN5fMnG8iFLgEJIBMxs9BvF3s3/FhuHyPKiVTZmXY0WY4ZyYqvoKR+XjaTRPPvBsDa4WI2u1zxXMeHlodT3lnCzVvyOYBLXL6CJgByuOxccJ8hnXfF9yY4F0aeL080Jz/3+EBNG8RO4ByhtBf4Ny8NQ6stWsjfeUIvH7bU/4zCYcYOq4WrInXHqS8qruDmIl7P5XXGcabuzQstPf/h2CRAUpP/PlHXcMlvewjmGU6MfDK+lifScNYwjPxRo4nKTGFZf/0aqHCh/EAsQyLKrOIYRE0lDG3bzBh8ogIMLAugsAfBb6M3mqCqKaTMAf/VAjh5FFJnjS+7bE+bZEV0qwax1CEoPPJL1fIQjOS8zj086gjpGRCtSy9+bTPTfTR/SJ+VUB5G2IeCItkNHpJX2ygojFZ9n5Fnj7R9ZnOM+L8nyIjPu3aePvtcrXlyLhH/hvOfIOjPxOlqW+O5QwSFP4OEcyLAUgDdUgyW36Z5mB285uKW/ighzZsOTevVUG2QwDItObIV6i8RCxFbN2oDHyPaO5j1tTaBNyVt8CAwEAAQ==`
	require.Equal(t, key, base64.StdEncoding.EncodeToString(derBytes))
}
