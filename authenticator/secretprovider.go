package authenticator

// noSecrets is a auth.SecretProvider that always fails authentication.
func noSecrets(user, realm string) string {
	// NOTE "Returning an empty string means failing the authentication."
	// (from godoc.org/github.com/abbot/go-http-auth)
	return ""

}

// aladdinsSecret is a auth.SecretProvider that contains the secret for Aladdin.
// This is from https://en.wikipedia.org/wiki/Basic_access_authentication.
func aladdinsSecret(user, realm string) string {
	if user == "Aladdin" {
		// OpenSesame as MD5Crypt
		return "$1$f61ef9e3$eZbgYlPxbNPJsF2Yb6.8G."
	}
	return noSecrets(user, realm)
}
