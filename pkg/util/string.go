package util

func AnalysisImage(image string) (string, string) {
	// image format: lang:version
	// example: python:3.8
	lang := ""
	version := ""
	for i := 0; i < len(image); i++ {
		if image[i] == ':' {
			lang = image[:i]
			version = image[i+1:]
			break
		}
	}
	return lang, version
}
