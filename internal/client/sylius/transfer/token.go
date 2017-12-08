package transfer

// Token is a Sylius API token object structure.
type Token struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int     `json:"expires_in"`
	TokenType    string  `json:"token_type"`
	Scope        *string `json:"scope"`
	RefreshToken string  `json:"refresh_token"`
}
