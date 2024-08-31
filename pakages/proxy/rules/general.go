package rules

// Define the OnOff struct with methods
type OnOff struct {
	On  []string `json:"on"`
	Off []string `json:"off"`
}

// Define the rules struct
type Rules struct {
	Domain OnOff `json:"domain"`
	IP     OnOff `json:"ip"`
}

// Customs + Defaults
var CombinedRules Rules

// Custom rules
var CustomRules Rules

// Default rules
var DefaultRules = Rules{
	Domain: OnOff{
		On: []string{
			"t.me",
			"api.myip.com",
			"play.google.com",
			"facebook.com",
			"twitter.com",
			"x.com",
			"youtube.com",
			"blogspot.com",
			"wordpress.com",
			"pinterest.com",
			"tumblr.com",
			"flickr.com",
			"telegram.org",
			"whatsapp.com",
			"line.me",
			"signal.org",
			"snapchat.com",
			"viber.com",
			"wechat.com",
			"discord.com",
			"reddit.com",
			"twitch.tv",
			"spotify.com",
			"soundcloud.com",
			"vimeo.com",
			"dailymotion.com",
			"hulu.com",
			"netflix.com",
			"bbc.com",
			"cnn.com",
			"foxnews.com",
			"nytimes.com",
			"washingtonpost.com",
			"theguardian.com",
			"news.google.com",
			"yahoo.com",
			"bing.com",
			"duckduckgo.com",
			"baidu.com",
			"linkedin.com",
			"microsoft.com",
			"dropbox.com",
			"onedrive.com",
			"paypal.com",
			"wikileaks.org",
			"amazon.com",
			"ebay.com",
			"alibaba.com",
			"booking.com",
			"airbnb.com",
			"cloudflare.com",
			"dash.cloudflare.com",
			"figma.com",
		},
		Off: []string{
			"master.jedimifarmaeid.online",
			"detectportal.firefox.com",
			"trustseal.enamad.ir",
			"enamad.ir",
			"github.com",
			"alive.github.com",
			"collector.github.com",
			"api.github.com",
			"meet.google.com",
			"linkedin.com",
			"static-asm.secure.skypeassets.com",
			"ocsp.pki.goog",
			"skype.com",
			"google.com",
			"digikala.com",
			"divar.ir",
			"hamshahrionline.ir",
			"bazar.ir",
			"irna.ir",
			"tejaratnews.com",
			"khabaronline.ir",
			"varzesh3.com",
			"iribnews.ir",
			"persianblog.ir",
			"khabarfarsi.com",
			"isna.ir",
			"yjc.ir",
			"mehrnews.com",
			"shahrekhabar.com",
			"irib.ir",
			"tasnimnews.com",
			"hamshahri.org",
			"zoomit.ir",
			"farsnews.ir",
			"tehran.ir",
			"bankmellat.ir",
			"ario.ir",
			"bama.ir",
			"khabarvarzeshi.com",
			"mashreghnews.ir",
			"egov.go.ir",
		},
	},
	IP: OnOff{
		On:  []string{},
		Off: []string{},
	},
}
