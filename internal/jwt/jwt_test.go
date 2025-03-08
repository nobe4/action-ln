// Important: All keys here were generated only for those tests, and don't
// expose any real credentials.
//
//nolint:gosec // Used only for testing.
package jwt

import (
	"errors"
	"testing"
)

const (
	now = 1000
	id  = "id"

	notPEMKey = "not PEM"

	// Generated with `openssl genrsa 512`.
	invalidKey = `
-----BEGIN PRIVATE KEY-----
MIIBVgIBADANBgkqhkiG9w0BAQEFAASCAUAwggE8AgEAAkEAwEZEAeUsBCGiC7GN
Q5bPaL8WZWoaJGUhhRSO9uY9JBkwneyob4qGIBXU+NILSUz6sJRsUbQmejuwMHRu
w8IoNQIDAQABAkEAhVe1jkLqtarFgKqPt1H9YT00QPzGSHtCNdK+GwgtWrxRBRDv
/bqnBy6YC6jvzTWofJPGUvWWqjyQqwmjay3PQQIhAO13Nuik2IyAJ69FkoWjntoD
MUyDYsK+IrojmMfLlHIdAiEAz0gWuYZNDs3ZWy+CGhwF/G+hW+a3+7PJxgG3D35z
svkCIQCLb17siC8ngPDMeBurIQJbnVhLRzKsixy1E8XYO2/0+QIhAKPic1UsAjD6
QCgAX/UUwwbbm9B1knHHrHiJUptFd2TBAiBf1s4FqIqEue9pVE9ZmvtRxU/KXsN2
f0wNsIEvMabLaw==
-----END PRIVATE KEY-----
`

	// Generated with `openssl rsa -in <(openssl genrsa 512) -traditional`.
	// https://pkg.go.dev/crypto/rsa#hdr-Minimum_key_size
	insecureKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALYmy5huoBwewANiDHhsa6ItURNHBDHXqWtYG2z0co4OrJXBuq5L
K20YkZ8Ah5Or4xJNHS3zdHuYeya4KbllBRsCAwEAAQJASQzX++LpzT07zl+FFsqM
g/bem/eQJBkUddtY7GJAit2EK9lcZ0BC0AMyCUT6caHDEfB/NfrcGHJzBZ3iBIAE
uQIhAOWQpcMPcz9/Uf0fKJOSFGNW5Ex+HUiE2B2AdGaALcj9AiEAyyBx0vjddBBk
Xqu+AhqfN06svc5LF3K6cvDWCyQUTfcCIAXPqpKMgpNZ6r5omoNZ0FBPc8oH6z/Z
tQrSJKAvoHkVAiEAl0kr3XDDJ02KI8SP+OsxCDVNjPRXkzg8y5y6HoQZp1MCIQCX
tlGFO2/1I0VhlFWE4hOxbMSxCGo27DYJTnfztOcWKw==
-----END RSA PRIVATE KEY-----
`

	// Generated with `openssl rsa -in <(openssl genrsa) -traditional`.
	validKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4S+DIySITr1EtEccoDdFW/g4/q6xCb/gBb/KXs2KftHhRLwO
fA96O69Uv/ffcNHvE1cBvjtpQ81d/K2KDJxGydMLg+cJgRFWxrk/EKXqnjJSPe3p
Wg4LDdMgZwmK6orr0bJXRXl2qXMiTBsgBLO8ayzw5B15BoSrMXNEa7bKrANhK/ir
dz0tDs1wT2QcK16sPkX2LDi5oNkhUMSJuFrhJ5cDuwdJWl64WqAYorGCfnjkwXJl
TPWs3ctLQ53AoyecFS2caXWdqiKPt+vSAc+VAR8Z3ceJYLHWputkC7+AUzAAj0I8
LjDE9pCcJ4HVn1luZEQFdromkS5wfK4OPuDnfQIDAQABAoIBAByDEY2fkIq3uD1D
S8KVfPi6Iy1MHSpo2wqfUBZU3BZWqLk1PnhC1W09M/PihK4aRrUiLRxTFW76T84w
guw51VS0nh6jYDaxZVVgGzYWa/B/2p3ww70dThUACHhDYw1zbYxtklM/n+CwrUUJ
ojI2N8MyO4YGnU5P+gUW3TDuQhcpUXwnmn7eUkiVIM4n2K+ValnpcM26B64MDdVU
fZ9rAoaHZncxbgw2GxQ2AOkvYVAcIErZd6YsyjI1aQV8jCEl3pFhHcit2Cbl+RWg
8773CnqCOH/ah600UR3Y1H0a8RUSRcdnObbCajeIS6I8gnkSa2/TnJ5jCF2uXvsw
1N4VR1kCgYEA/h/8BmEKUWqLPq1/PfRWcbIl4BoOrFCkrcHJ0xCWxhFiujTH7JLb
H/YHvnvc5PWeeiVb7yULilEIIPkhvvW73KvxiZRcubm+Q1yTz4GdFiV/PxfY+AC7
JhBJfJzNtDZrnH+NSSG/SmLCx7lJV1xBkYVpCv2eJYkXeCdBzV0WtB8CgYEA4tjd
R8eA/5X2QhEvAeVl1QHbwq0SQOS6mCqkfKicP5EN8JgVHsfBkMs7tBNyfpN7sSqQ
biVPQKEi05KQsD2eGuBETkHjfiuPTgj4k6AiaJj3Ec8IMEPEJXI66edHveusLXcm
QsR+xrgBzuw2ANXqX3w5TnLzMi+OV51Pn0Ha0OMCgYEAhA8JACOjogWVEOBGVGLK
HVFvn1LLNz69JVKkWBuxzoIwZQWSs1zppGVNRu7FLvJ5BY6uhMsigSF08PWmVL8M
fjOYVF+WBCoDNqxAX8BCasTXqGjzJoXyu2gRWEGAIFt7dptOR6fS6YwDHpkqBMz7
gezrVnvPmD/yw0zbRCZQ6w8CgYA39phgpO9GHpDqK6MVLKq6qgK1PE1MhSEjeSGr
P02MwRRXTq2nMlCmj/ziqAmPAIN7aazH/5xVrWsSFw5q7EidCMbRJ6Af+E8aSUxJ
3y+d7l7FnfW/MnipZEz0d4JTcFjBvqtJvYApNiv8CHoqKpvvgo4AtIsaznCnXL/P
4kdBUQKBgDaLS8c19wG70alGeQ2fr9xqZwPE1nvbP3BG3juiZNT4PYvuMCYlOLBd
og6wflp3gAUJwYCHxaCSagN+yK4CKbYmWzkuIvel0UVmj+4uQBu3y5vg3aXSYAG/
qvvvQXkidoag4IXUQDcm8pRGE+ZAvXwHfuYPrWNVqS1j/2TfbJX3
-----END RSA PRIVATE KEY-----
`

	// Corresponds to the validKey test.
	wantSignature = "eyJ0eXAiOiJKV1QiLCAiYWxnIjoiUlMyNTYifQ." +
		"eyJpYXQiOjk0MCwiZXhwIjoxNjAwLCJpc3MiOiJpZCJ9." +
		"B91XuqtyYTpSXCPxgh0lTRl2BFLEylW8UPlCNh0IEth5jM5L-3QIHLpEGqkR-dRD2waCL-67_7uX8ZDqqHFcLVydEDO-zr3Kq6TelpBujKW0e6RO3X3CM3_y40d0MtmadB01wUoAvIbCo5NDvcSoBXFMz9mQHUFYP0WsXB0QBmDfZiaHkZsuSOB6ioHVotvDPoa6uM4CQJKLfV-dVLNOh2leahd0IS72ME7hb2l4hTX1Wd3jVXeQwVHnZ9p8Os4S3_YN_iH0mRDxXbdb3XbKSYH0ixB2GEpnx4IXvEdRXhliOiJUcHVL8_zyzSBK7vPBkZT-0K6XMas1Q08traZIEw"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("invalid PEM", func(t *testing.T) {
		t.Parallel()

		_, err := New(now, id, notPEMKey)
		if !errors.Is(err, ErrNotPEM) {
			t.Errorf("expected %q, got %q", ErrNotPEM, err)
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		t.Parallel()

		_, err := New(now, id, invalidKey)
		if !errors.Is(err, ErrInvalidKey) {
			t.Errorf("expected %q, got %q", ErrInvalidKey, err)
		}
	})

	t.Run("insecure key", func(t *testing.T) {
		t.Parallel()

		_, err := New(now, id, insecureKey)
		if !errors.Is(err, ErrCannotSign) {
			t.Errorf("expected %q, got %q", ErrCannotSign, err)
		}
	})

	t.Run("valid key", func(t *testing.T) {
		t.Parallel()

		got, err := New(now, id, validKey)
		if err != nil {
			t.Errorf("want no error, got %q", err)
		}

		if got != wantSignature {
			t.Errorf("want %q, got %q", wantSignature, got)
		}
	})
}

func TestEncode(t *testing.T) {
	t.Parallel()

	want := "aGVsbG8"
	got := encode("hello")

	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
