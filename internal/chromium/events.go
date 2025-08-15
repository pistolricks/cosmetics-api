package chromium

import (
	"fmt"

	"github.com/go-rod/rod/lib/proto"
)

func (chrome ChromeClient) GetEventDetails() *proto.NetworkResponse {

	page := chrome.Client.Browser.MustPage().MustWindowNormal()

	e := proto.NetworkResponseReceived{}
	wait := page.WaitEvent(&e)
	page.MustNavigate("https://mall.riman.com/checkout/billing?orderNum=0x02000000A9A4881D3585EEABD6EEF0ADB8468A4260B92A9F8F4585A59402818255ADE5E05C3906EF79EEDA499508576C3F30D92E")
	wait()

	fmt.Println(e.Response)

	return e.Response
}
