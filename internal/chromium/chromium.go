package chromium

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

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

func ChromeBrowser() *rod.Browser {
	path, _ := launcher.LookPath()

	u := launcher.
		NewUserMode().
		UserDataDir("path").
		Headless(false).
		NoSandbox(true).
		Bin(path).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect().NoDefaultDevice()

	return browser

}
