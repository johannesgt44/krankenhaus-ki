package config

import "testing"

func TestOIDCDefaultsSindAktiv(t *testing.T) {
	t.Setenv("OIDC_ENABLED", "")
	t.Setenv("OIDC_ISSUER_URL", "")
	t.Setenv("OIDC_CLIENT_ID", "")
	t.Setenv("OIDC_REQUIRED_ROLE", "")

	konfig := Laden()

	if !konfig.OIDC.Aktiv {
		t.Fatal("OIDC sollte standardmaessig aktiv sein")
	}
	if konfig.OIDC.IssuerURL != "http://localhost:8880/realms/python" {
		t.Fatalf("IssuerURL = %q", konfig.OIDC.IssuerURL)
	}
	if konfig.OIDC.ClientID != "python-client" {
		t.Fatalf("ClientID = %q", konfig.OIDC.ClientID)
	}
	if konfig.OIDC.ErforderlicheRolle != "admin" {
		t.Fatalf("ErforderlicheRolle = %q", konfig.OIDC.ErforderlicheRolle)
	}
}
