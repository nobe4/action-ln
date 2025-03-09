package github

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

const (
	//nolint:gosec // This is not a secret.
	appTokenAPIPath = "/app/installations/123/access_tokens"
	jwtToken        = "jwt"

	appID        = "app-id"
	appInstallID = "123"

	// Taken from `jwt_test.go`.
	//nolint:gosec // Used only for testing.
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
-----END RSA PRIVATE KEY-----`
	wantSignature = "eyJ0eXAiOiJKV1QiLCAiYWxnIjoiUlMyNTYifQ." +
		"eyJpYXQiOjk0MCwiZXhwIjoxNjAwLCJpc3MiOiJpZCJ9." +
		"B91XuqtyYTpSXCPxgh0lTRl2BFLEylW8UPlCNh0IEth5jM5L-3QIHLpEGqkR-dRD2waCL-67_7uX8ZDqqHFcLVydEDO-zr3Kq6TelpBujKW0e6RO3X3CM3_y40d0MtmadB01wUoAvIbCo5NDvcSoBXFMz9mQHUFYP0WsXB0QBmDfZiaHkZsuSOB6ioHVotvDPoa6uM4CQJKLfV-dVLNOh2leahd0IS72ME7hb2l4hTX1Wd3jVXeQwVHnZ9p8Os4S3_YN_iH0mRDxXbdb3XbKSYH0ixB2GEpnx4IXvEdRXhliOiJUcHVL8_zyzSBK7vPBkZT-0K6XMas1Q08traZIEw"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	t.Run("succeeds with a token", func(t *testing.T) {
		t.Parallel()

		g := GitHub{}

		if err := g.Auth(t.Context(), token, "", "", ""); err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		if g.Token != token {
			t.Fatalf("expected token to be %q but got %q", token, g.Token)
		}
	})

	t.Run("fails to get the app JWT", func(t *testing.T) {
		t.Parallel()

		g := GitHub{}

		err := g.Auth(t.Context(), "", appID, "invalid private key", appInstallID)
		if !errors.Is(err, errGetJWT) {
			t.Fatalf("expected error %q, got %q", errGetJWT, err)
		}
	})

	t.Run("fails to get the app token", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		err := g.Auth(t.Context(), "", appID, validKey, appInstallID)
		if !errors.Is(err, errGetAppToken) {
			t.Fatalf("expected error %q, got %q", errGetAppToken, err)
		}
	})

	t.Run("succeeds with an app", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"token": "%s"}`, token)
		})

		err := g.Auth(t.Context(), "", appID, validKey, appInstallID)
		if err != nil {
			t.Fatalf("expected no error, got %q", err)
		}

		if g.Token != token {
			t.Fatalf("expected token %q, got %q", token, g.Token)
		}
	})
}

func TestGetAppToken(t *testing.T) {
	t.Parallel()

	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodPost, appTokenAPIPath, nil)

		if auth := r.Header.Get("Authorization"); auth != "Bearer "+jwtToken {
			t.Fatal("invalid jwt", auth)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"token": "%s"}`, token)
	})

	got, err := g.GetAppToken(t.Context(), appInstallID, jwtToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != token {
		t.Fatalf("expected token to be '%s' but got '%s'", token, got)
	}
}
