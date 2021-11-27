package utils

 func AddEllipsis(str string, maxSize int) string{
	  if len(str) > maxSize {
		      return str[:maxSize]+"..."
	  }
	  return str
}
