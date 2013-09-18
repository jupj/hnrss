package main

import (
	"fmt"
	"os"
	"testing"
)

type testItem struct {
	url,
	title,
	domain,
	comments,
	guid string
}

var tests = []testItem{
	testItem{url: "http://www.nytimes.com/2013/09/06/us/nsa-foils-much-internet-encryption.html",
		title:    "N.S.A. Foils Much Internet Encryption",
		domain:   " (nytimes.com) ",
		comments: "https://news.ycombinator.com/item?id=6336178",
		guid:     "6336178"},

	testItem{url: "http://www.BeeLineReader.com",
		title:    "Show HN: Ditch Black Text to Read Faster, Easier",
		domain:   " (BeeLineReader.com) ",
		comments: "https://news.ycombinator.com/item?id=6335784",
		guid:     "6335784"},

	testItem{url: "http://rigsomelight.com/2013/09/09/frameless-geodesic-dome.html",
		title:    "The Frameless Geodesic Dome I currently live in",
		domain:   " (rigsomelight.com) ",
		comments: "https://news.ycombinator.com/item?id=6355488",
		guid:     "6355488"},

	testItem{url: "https://www.usenix.org/blog/my-daughters-high-school-programming-teacher",
		title:    "To my daughter's high school programming teacher",
		domain:   " (usenix.org) ",
		comments: "https://news.ycombinator.com/item?id=6357317",
		guid:     "6357317"},

	testItem{url: "http://blog.eat24hours.com/how-to-advertise-on-a-porn-website/",
		title:    "How to Advertise on a Porn Website",
		domain:   " (eat24hours.com) ",
		comments: "https://news.ycombinator.com/item?id=6359555",
		guid:     "6359555"},

	testItem{url: "http://jacoboneal.com/car-engine/",
		title:    "How a Car Engine Works",
		domain:   " (jacoboneal.com) ",
		comments: "https://news.ycombinator.com/item?id=6332385",
		guid:     "6332385"},

	testItem{
		title:    "0 A.D., an Open-Source Strategy Game",
		domain:   " (indiegogo.com) ",
		url:      "http://www.indiegogo.com/projects/support-0-a-d-an-open-source-strategy-game/",
		comments: "https://news.ycombinator.com/item?id=6339917",
		guid:     "6339917"},

	testItem{
		title:    "PayPal Freezes Mailpile Campaign Funds",
		domain:   " (mailpile.is) ",
		url:      "http://www.mailpile.is/blog/2013-09-05_PayPal_Freezes_Campaign_Funds.html",
		comments: "https://news.ycombinator.com/item?id=6333203",
		guid:     "6333203"},

	testItem{
		title:    "Android is for startups",
		domain:   " (audobox.com) ",
		url:      "http://blog.audobox.com/android-is-for-startups/",
		comments: "https://news.ycombinator.com/item?id=6329490",
		guid:     "6329490"},

	testItem{
		title:    " The US government has betrayed the internet. We need to take it back",
		domain:   " (theguardian.com) ",
		url:      "http://www.theguardian.com/commentisfree/2013/sep/05/government-betrayed-internet-nsa-spying",
		comments: "https://news.ycombinator.com/item?id=6336373",
		guid:     "6336373"},

	testItem{
		title:    "Left with Nothing",
		domain:   " (washingtonpost.com) ",
		url:      "http://www.washingtonpost.com/sf/investigative/2013/09/08/left-with-nothing/",
		comments: "https://news.ycombinator.com/item?id=6349476",
		guid:     "6349476"},

	testItem{
		title:    "YC Will Now Fund Nonprofits Too",
		domain:   " (ycombinator.com) ",
		url:      "http://ycombinator.com/np.html",
		comments: "https://news.ycombinator.com/item?id=6341568",
		guid:     "6341568"},

	testItem{
		title:    "NSA Codebreaking: I Am The Other",
		domain:   " (popehat.com) ",
		url:      "http://www.popehat.com/2013/09/06/nsa-codebreaking-i-am-the-other/",
		comments: "https://news.ycombinator.com/item?id=6341570",
		guid:     "6341570"},

	testItem{
		title:    "Unix Commands I Wish I’d Discovered Years Earlier",
		domain:   " (atomicobject.com) ",
		url:      "http://spin.atomicobject.com/2013/09/09/5-unix-commands/",
		comments: "https://news.ycombinator.com/item?id=6360320",
		guid:     "6360320"},

	testItem{
		title:    "Postgresql 9.3 Released",
		domain:   " (postgresql.org) ",
		url:      "http://www.postgresql.org/about/news/1481/",
		comments: "https://news.ycombinator.com/item?id=6353140",
		guid:     "6353140"},

	testItem{
		title:    "School is a prison and damaging our kids",
		domain:   " (salon.com) ",
		url:      "http://www.salon.com/2013/08/26/school_is_a_prison_and_damaging_our_kids/",
		comments: "https://news.ycombinator.com/item?id=6329191",
		guid:     "6329191"},

	testItem{
		title:    "Don't trust me: I might be a spook",
		domain:   " (daemonology.net) ",
		url:      "http://www.daemonology.net/blog/2013-09-10-I-might-be-a-spook.html",
		comments: "https://news.ycombinator.com/item?id=6359719",
		guid:     "6359719"},

	testItem{
		title:    "3-Sweep: Extracting Editable Objects from a Single Photo [video]",
		domain:   " (youtube.com) ",
		url:      "http://www.youtube.com/watch?v=Oie1ZXWceqM&hd=1",
		comments: "https://news.ycombinator.com/item?id=6363672",
		guid:     "6363672"},

	testItem{
		title:    "NSA Spying Documents to be Released As Result of EFF Lawsuit",
		domain:   " (eff.org) ",
		url:      "https://www.eff.org/deeplinks/2013/09/hundreds-pages-nsa-spying-documents-be-released-result-eff-lawsuit",
		comments: "https://news.ycombinator.com/item?id=6332657",
		guid:     "6332657"},

	testItem{
		title:    "Ten Years of Bootstrapping: Lessons Learned",
		domain:   " (davegooden.com) ",
		url:      "http://davegooden.com/2013/08/our-story-10-years-of-bootstrapping/",
		comments: "https://news.ycombinator.com/item?id=6334806",
		guid:     "6334806"},

	testItem{
		title:    "Statement of Condemnation of U.S. Mass-Surveillance Programs [pdf]",
		domain:   " (ucdavis.edu) ",
		url:      "http://www.cs.ucdavis.edu/~rogaway/politics/surveillance.pdf",
		comments: "https://news.ycombinator.com/item?id=6356514",
		guid:     "6356514"},

	testItem{
		title:    "Hermit: a font for programmers, by a programmer",
		domain:   " (pcaro.es) ",
		url:      "http://pcaro.es/p/hermit",
		comments: "https://news.ycombinator.com/item?id=6354396",
		guid:     "6354396"},

	testItem{
		title:    "Why you should not pirate Google’s geo APIs",
		domain:   " (petewarden.com) ",
		url:      "http://petewarden.com/2013/09/09/why-you-should-stop-pirating-googles-geo-apis/",
		comments: "https://news.ycombinator.com/item?id=6356399",
		guid:     "6356399"},

	testItem{
		title:    "Quark: A secure Web Browser with a Formally Verified Kernel",
		domain:   " (ucsd.edu) ",
		url:      "http://goto.ucsd.edu/quark/",
		comments: "https://news.ycombinator.com/item?id=6348532",
		guid:     "6348532"},

	testItem{
		title:    "Speculation on \"BULLRUN\"",
		domain:   " (mail-archive.com) ",
		url:      "http://www.mail-archive.com/cryptography@metzdowd.com/msg12325.html",
		comments: "https://news.ycombinator.com/item?id=6346531",
		guid:     "6346531"},

	testItem{
		title:    "Why you should not trust emails sent from Google",
		domain:   " (vagosec.org) ",
		url:      "http://vagosec.org/2013/09/google-scholar-email-html-injection/",
		comments: "https://news.ycombinator.com/item?id=6364481",
		guid:     "6364481"},

	testItem{
		title:    "Obama administration had restrictions on NSA reversed in 2011",
		domain:   " (washingtonpost.com) ",
		url:      "http://www.washingtonpost.com/world/national-security/obama-administration-had-restrictions-on-nsa-reversed-in-2011/2013/09/07/c26ef658-0fe5-11e3-85b6-d27422650fd5_story.html",
		comments: "https://news.ycombinator.com/item?id=6347790",
		guid:     "6347790"},

	testItem{
		title:    "Apple Unveils The iPhone 5S",
		domain:   " (techcrunch.com) ",
		url:      "http://techcrunch.com/2013/09/10/apple-unveils-the-iphone-5s/",
		comments: "https://news.ycombinator.com/item?id=6361558",
		guid:     "6361558"},

	testItem{
		title:    "A new breed of Chrome Apps",
		domain:   " (chrome.blogspot.com) ",
		url:      "http://chrome.blogspot.com/2013/09/a-new-breed-of-chrome-apps.html",
		comments: "https://news.ycombinator.com/item?id=6335016",
		guid:     "6335016"},

	testItem{
		title:    "How to remain secure against NSA surveillance",
		domain:   " (theguardian.com) ",
		url:      "http://www.theguardian.com/world/2013/sep/05/nsa-how-to-remain-secure-surveillance",
		comments: "https://news.ycombinator.com/item?id=6336523",
		guid:     "6336523"}}

func TestHnParsing(t *testing.T) {
	htmlReader, err := os.Open("hn_test.html")
	if err != nil {
		t.Errorf("Couldn't open 'hn_test.html', error: %s", err)
	}
	defer htmlReader.Close()

	rss, err := parseHnHtmlToRss(htmlReader)
	if err != nil {
		t.Error(err)
	}

	if (rss.Version != "2.0") ||
		(rss.Title != "Hacker News Top Links") ||
		(rss.Link != "https://news.ycombinator.com/best") ||
		(rss.Description != "Links for the intellectually curious, ranked by readers.") {
		t.Errorf("Wrong header")
	}

	if len(tests) != len(rss.Items) {
		t.Errorf("Wrong number of test cases, expected %d but found %d.", len(tests), len(rss.Items))
	}

	for i, ri := range rss.Items {
		ti := tests[i]

		if fmt.Sprintf("%s%s", ti.title, ti.domain) != ri.Title {
			t.Errorf("Expected title '%s%s', found '%s'", ti.title, ti.domain, ri.Title)
		}

		if ti.url != ri.Link {
			t.Errorf("Expected url '%s', found '%s'", ti.url, ri.Link)
		}

		if ti.comments != ri.Comments {
			t.Errorf("Expected comment-url '%s', found '%s'", ti.comments, ri.Comments)
		}

		if ti.guid != ri.Guid {
			t.Errorf("Expected guid '%s', found '%s'", ti.guid, ri.Guid)
		}
	}
}
