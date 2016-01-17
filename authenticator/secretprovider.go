package authenticator

// noSecrets is a auth.SecretProvider that always fails authentication.
func noSecrets(user, realm string) string {
	// NOTE "Returning an empty string means failing the authentication."
	// (from godoc.org/github.com/abbot/go-http-auth)
	return ""

}
