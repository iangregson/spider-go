package urls

import "net/url"

func ValidUrl(urlstr string, baseUrl *url.URL) (bool, *url.URL) {
	validUrl, err := url.Parse(urlstr)
	if err != nil {
		return false, nil
	}
	if len(validUrl.Host) == 0 {
		validUrl = baseUrl.ResolveReference(validUrl)
	}
	if len(validUrl.Path) == 0 {
		validUrl.Path = "/"
	}
	return true, validUrl
}

func SameHost(base *url.URL, target *url.URL) bool {
	return base.Host == target.Host
}

func ExternalHost(base *url.URL, target *url.URL) bool {
	return base.Host != target.Host
}
