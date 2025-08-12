package chromium

import "github.com/go-rod/rod"

type ChromeConfig struct {
	Browser *rod.Browser
	Page    *rod.Page
}

type ChromeConnector struct {
	Chrome ChromeClient
}

func NewChromeConnector(config *ChromeConfig) ChromeConnector {
	return ChromeConnector{
		Chrome: ChromeClient{Client: config},
	}
}
