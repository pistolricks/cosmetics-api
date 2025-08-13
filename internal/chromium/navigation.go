package chromium

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func (chrome ChromeClient) ChromeLogin(loginUrl string, rimanStoreName string, username string, password string) (*rod.Page, *rod.Browser, []*proto.NetworkCookie) {

	// --allow-third-party-cookies

	officeUrl := fmt.Sprintf("https://myoffice.riman.com")
	homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	client := chrome.Client

	browser := client.Browser
	page := browser.MustPage().MustWindowNormal()

	client.Page = page

	wait := client.Page.MustWaitNavigation()
	client.Page.MustNavigate(officeUrl)
	wait()

	url := client.Page.MustInfo().URL

	switch {
	case strings.Contains(url, "https://myoffice-1.riman.com/login"):
		fmt.Println("office login")

		client.Page.MustElement("#username").MustSelectAllText().MustInput(username)
		client.Page.MustElement("#password").MustSelectAllText().MustInput(password)
		client.Page.MustElement("#loginBtn").MustClick()

	case strings.Contains(url, officeUrl):
		fmt.Println("office logged in")
		client.Page.MustNavigate(homeUrl)
	}

	cookies := browser.MustGetCookies()

	return page, browser, cookies
}

func (chrome ChromeClient) ChromeHomePage(rimanStoreName string) []*proto.NetworkCookie {
	homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	// client.Page.MustSetCookies(networkCookie...)

	wait := chrome.Client.Page.MustWaitNavigation()
	chrome.Client.Page.MustNavigate(homeUrl)
	wait()

	cookies := chrome.Client.Browser.MustGetCookies()

	return cookies
}

func (chrome ChromeClient) SubmitForm() *rod.Element {

	el := chrome.Client.Page.Timeout(2 * time.Second).MustElement("[type=submit]")
	el.MustType(input.Enter)

	return el
}
