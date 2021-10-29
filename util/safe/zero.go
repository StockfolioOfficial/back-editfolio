package safe

func StringOrZero(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}