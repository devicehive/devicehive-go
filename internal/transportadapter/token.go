package transportadapter

type token struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
