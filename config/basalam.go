package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("basalam", map[string]any{
		// OAuth client credentials (from the Basalam developer panel).
		// Leave as placeholders for now and fill them in your .env.
		"client_id":     config.Env("BASALAM_CLIENT_ID", ""),
		"client_secret": config.Env("BASALAM_CLIENT_SECRET", ""),

		// The callback URL registered for this client in the Basalam panel.
		// Must exactly match, e.g. http://127.0.0.1:3000/auth/basalam/callback
		"redirect_uri": config.Env("BASALAM_REDIRECT_URI", ""),

		// Space separated list of scopes. For login we only need the user identity.
		"scope": config.Env("BASALAM_SCOPE", "customer.profile.read"),

		// Basalam OAuth endpoints (overridable, sensible defaults).
		"sso_url":    config.Env("BASALAM_SSO_URL", "https://basalam.com/accounts/sso"),
		"token_url":  config.Env("BASALAM_TOKEN_URL", "https://auth.basalam.com/oauth/token"),
		"whoami_url": config.Env("BASALAM_WHOAMI_URL", "https://auth.basalam.com/whoami"),
	})
}
