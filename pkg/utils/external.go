package utils

import "golang-api/internal/config"

func GetExternalUrl(key string, path string) string {
	var baseUrl string

	switch key {
    case "main":
      baseUrl = config.MainUrl
    case "cdn":
      baseUrl = config.CdnUrl
    default:
      baseUrl = config.MainUrl
  }

	return baseUrl + "/" + path
}