package security_test

import (
	"testing"

	"github.com/johannesgt44/krankenhaus-ki/internal/security"
)

func TestKeycloakClaimsHatClientRolle(t *testing.T) {
	claims := security.KeycloakClaims{
		ResourceAccess: map[string]security.RollenZugriff{
			"javascript-client": {Roles: []string{"admin"}},
		},
	}

	if !claims.HatRolle("javascript-client", "admin") {
		t.Fatal("admin Rolle wurde nicht gefunden")
	}
}

func TestKeycloakClaimsLehntFehlendeClientRolleAb(t *testing.T) {
	claims := security.KeycloakClaims{
		ResourceAccess: map[string]security.RollenZugriff{
			"javascript-client": {Roles: []string{"user"}},
		},
	}

	if claims.HatRolle("javascript-client", "admin") {
		t.Fatal("admin Rolle wurde faelschlich gefunden")
	}
}

func TestKeycloakClaimsAkzeptiertRealmRolleAlsFallback(t *testing.T) {
	claims := security.KeycloakClaims{
		RealmAccess: security.RollenZugriff{Roles: []string{"admin"}},
	}

	if !claims.HatRolle("javascript-client", "admin") {
		t.Fatal("realm admin Rolle wurde nicht als Fallback akzeptiert")
	}
}
