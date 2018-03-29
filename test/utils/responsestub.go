package utils

var ResponseStub = &responseStub{}

type responseStub struct {}

func (s *responseStub) Authenticate(reqId string) map[string]string {
	return map[string]string{
		"action": "authenticate",
		"requestId": reqId,
		"status": "success",
	}
}

func (s *responseStub) Token(reqId, accessToken, refreshToken string) map[string]string {
	return map[string]string{
		"action": "token",
		"requestId": reqId,
		"status": "success",
		"accessToken": accessToken,
		"refreshToken": refreshToken,
	}
}

func (s *responseStub) TokenRefresh(reqId, accessToken string) map[string]string {
	return map[string]string{
		"action": "token/refresh",
		"requestId": reqId,
		"status": "success",
		"accessToken": accessToken,
	}
}