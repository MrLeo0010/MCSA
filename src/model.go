package src

// ParseMOTD безопасно извлекает текстовый MOTD
func (s *StatusResponse) ParseMOTD() string {
	if s.Description == nil {
		return ""
	}
	switch v := s.Description.(type) {
	case string:
		return CleanMOTD(v)
	case map[string]interface{}:
		if text, ok := v["text"].(string); ok {
			return CleanMOTD(text)
		}
	}
	return ""
}
