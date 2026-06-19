package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/johannesgt44/krankenhaus-ki/internal/problem"
)

type OIDCKonfiguration struct {
	IssuerURL          string
	ClientID           string
	ErforderlicheRolle string
}

type Verifizierer interface {
	Verifizieren(ctx context.Context, token string, claims any) error
}

type Autorisierer struct {
	verifizierer       Verifizierer
	clientID           string
	erforderlicheRolle string
}

type KeycloakClaims struct {
	Audience        Zielgruppen              `json:"aud"`
	AuthorizedParty string                   `json:"azp"`
	ResourceAccess  map[string]RollenZugriff `json:"resource_access"`
	RealmAccess     RollenZugriff            `json:"realm_access"`
}

type RollenZugriff struct {
	Roles []string `json:"roles"`
}

type Zielgruppen []string

func NeuerOIDCAutorisierer(ctx context.Context, konfig OIDCKonfiguration) (*Autorisierer, error) {
	if strings.TrimSpace(konfig.IssuerURL) == "" {
		return nil, errors.New("OIDC issuer URL fehlt")
	}
	if strings.TrimSpace(konfig.ClientID) == "" {
		return nil, errors.New("OIDC client ID fehlt")
	}
	if strings.TrimSpace(konfig.ErforderlicheRolle) == "" {
		return nil, errors.New("OIDC erforderliche Rolle fehlt")
	}

	provider, err := oidc.NewProvider(ctx, konfig.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("OIDC provider konnte nicht geladen werden: %w", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: konfig.ClientID, SkipClientIDCheck: true})
	return NeuerAutorisierer(oidcVerifizierer{verifier: verifier}, konfig.ClientID, konfig.ErforderlicheRolle), nil
}

func NeuerAutorisierer(verifizierer Verifizierer, clientID string, erforderlicheRolle string) *Autorisierer {
	return &Autorisierer{
		verifizierer:       verifizierer,
		clientID:           clientID,
		erforderlicheRolle: erforderlicheRolle,
	}
}

func (a *Autorisierer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			a.schreibeUnauthorized(w)
			return
		}

		var claims KeycloakClaims
		if a.verifizierer == nil || a.verifizierer.Verifizieren(r.Context(), token, &claims) != nil {
			a.schreibeUnauthorized(w)
			return
		}
		if !claims.HatClient(a.clientID) {
			a.schreibeUnauthorized(w)
			return
		}
		if !claims.HatRolle(a.clientID, a.erforderlicheRolle) {
			problem.Schreiben(w, http.StatusForbidden, "Der Token enthaelt nicht die erforderliche Rolle.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c KeycloakClaims) HatClient(clientID string) bool {
	if c.AuthorizedParty == clientID {
		return true
	}
	for _, audience := range c.Audience {
		if audience == clientID {
			return true
		}
	}
	return false
}

func (c KeycloakClaims) HatRolle(clientID string, rolle string) bool {
	if zugriff, ok := c.ResourceAccess[clientID]; ok && zugriff.enthaelt(rolle) {
		return true
	}
	return c.RealmAccess.enthaelt(rolle)
}

func (z RollenZugriff) enthaelt(rolle string) bool {
	for _, vorhanden := range z.Roles {
		if vorhanden == rolle {
			return true
		}
	}
	return false
}

func (z *Zielgruppen) UnmarshalJSON(data []byte) error {
	var mehrere []string
	if err := json.Unmarshal(data, &mehrere); err == nil {
		*z = mehrere
		return nil
	}

	var einzelne string
	if err := json.Unmarshal(data, &einzelne); err != nil {
		return err
	}
	*z = []string{einzelne}
	return nil
}

func bearerToken(header string) (string, bool) {
	felder := strings.Fields(header)
	if len(felder) != 2 || !strings.EqualFold(felder[0], "Bearer") || strings.TrimSpace(felder[1]) == "" {
		return "", false
	}
	return felder[1], true
}

func (a *Autorisierer) schreibeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="python"`)
	problem.Schreiben(w, http.StatusUnauthorized, "Ein gueltiger Bearer-Token ist erforderlich.")
}

type oidcVerifizierer struct {
	verifier *oidc.IDTokenVerifier
}

func (v oidcVerifizierer) Verifizieren(ctx context.Context, token string, claims any) error {
	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return err
	}
	return idToken.Claims(claims)
}
