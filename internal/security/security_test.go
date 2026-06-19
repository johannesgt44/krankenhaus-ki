package security_test

import (
	"testing"

	"github.com/johannesgt44/krankenhaus-ki/internal/security"
)

func TestKeycloakClaimsHatClientRolle(t *testing.T) {
	claims := security.KeycloakClaims{
		AuthorizedParty: "python-client",
		ResourceAccess: map[string]security.RollenZugriff{
			"python-client": {Roles: []string{"admin"}},
		},
	}

	if !claims.HatRolle("python-client", "admin") {
		t.Fatal("admin Rolle wurde nicht gefunden")
	}
}

func TestKeycloakClaimsAkzeptiertAuthorizedPartyAlsClient(t *testing.T) {
	claims := security.KeycloakClaims{AuthorizedParty: "python-client"}

	if !claims.HatClient("python-client") {
		t.Fatal("azp wurde nicht als Client akzeptiert")
	}
}

func TestKeycloakClaimsAkzeptiertAudienceAlsClient(t *testing.T) {
	claims := security.KeycloakClaims{Audience: security.Zielgruppen{"account", "python-client"}}

	if !claims.HatClient("python-client") {
		t.Fatal("aud wurde nicht als Client akzeptiert")
	}
}

func TestKeycloakClaimsLehntFalschenClientAb(t *testing.T) {
	claims := security.KeycloakClaims{
		AuthorizedParty: "anderer-client",
		Audience:        security.Zielgruppen{"account"},
	}

	if claims.HatClient("python-client") {
		t.Fatal("falscher Client wurde akzeptiert")
	}
}

func TestKeycloakClaimsLehntFehlendeClientRolleAb(t *testing.T) {
	claims := security.KeycloakClaims{
		ResourceAccess: map[string]security.RollenZugriff{
			"python-client": {Roles: []string{"user"}},
		},
	}

	if claims.HatRolle("python-client", "admin") {
		t.Fatal("admin Rolle wurde faelschlich gefunden")
	}
}

func TestKeycloakClaimsAkzeptiertRealmRolleAlsFallback(t *testing.T) {
	claims := security.KeycloakClaims{
		RealmAccess: security.RollenZugriff{Roles: []string{"admin"}},
	}

	if !claims.HatRolle("python-client", "admin") {
		t.Fatal("realm admin Rolle wurde nicht als Fallback akzeptiert")
	}
}
