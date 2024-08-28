package rules

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

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

// Has checks if a value exists in the On or Off slice
func (oo OnOff) Has(slice string, value string) bool {
	if slice == "on" {
		for _, v := range oo.On {
			if v == value {
				return true
			}
		}
	} else if slice == "off" {
		for _, v := range oo.Off {
			if v == value {
				return true
			}
		}
	}
	return false
}

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

// Customs + Defaults
var combinedRules Rules

// get custom rules on init
func init() {
	// Get the directory of the executable
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Construct the full path to the JSON file
	configPath := filepath.Join(filepath.Dir(path), "/golang-server/.configs/custom-rules.json")

	// Open the JSON file
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error opening file '%s': %v", configPath, err)
	}
	defer file.Close()

	// Decode JSON data into rules struct
	var rules Rules
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rules); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Cache the rules data
	CustomRules = rules

	// Combine the Domain and IP rules from both default and custom rules
	combinedRules = Rules{
		Domain: OnOff{
			On:  append(CustomRules.Domain.On, DefaultRules.Domain.On...),
			Off: append(CustomRules.Domain.Off, DefaultRules.Domain.Off...),
		},
		IP: OnOff{
			On:  append(CustomRules.IP.On, DefaultRules.IP.On...),
			Off: append(CustomRules.IP.Off, DefaultRules.IP.Off...),
		},
	}
}
