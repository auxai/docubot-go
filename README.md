Docubot
======================

What is Docubot?
----
Docubot™ is artificial intelligence, designed specifically for legal websites. This handy, document-generating plug-in, is looking to tap into the estimated $45 billion dollars of unspent revenue by people not currently participating in the legal market. It is an augmented legal service that gives consumers access to legal services without having to sit down with a lawyer, they simply log in, select the document to be generated, and Docubot™ guides them through the process. If the consumer needs help, at any time during the document-generating process, they can type ‘help’ and Docubot’s On-Demand video or chat will connect them to a legal professional in your office for assistance.

Why Docubot™ For Legal Professionals?
----
* Allows legal professionals to efficiently and effectively help more individuals
* Creates opportunity to build a personal connection with consumers for future 		legal help
* Fully customizable to offer the services most requested by legal professional
* Expands legal professional’s reach by being available from anywhere
* Reduces appointments for smaller services, allowing more time for larger cases
* Taps into an estimated $45 billion of unspent revenue by people not currently participating in the legal market.

Why Docubot™ For Legal Consumers?
----
* Receive legal services at a reduced rate or even free.
* Perfect for those lacking a personal connection to a lawyer
* Simple, step-by-step guidance reduces legal system anxieties
* Easy access for consumers hampered with geographical constraints
* Good fit for consumer that fall into the justice gap—have too much for pro bono 		work, but unable to afford traditional private legal services

How To Use
----
To use this Go library for Docubot™ simply:

`go get github.com/auxai/docubot-go`

Then use in your program as follows:

```go
import "github.com/auxai/docubot-go"
docubotURL := "https://docubotapi.1law.com"
key := "" // Your API Key
secret := "" // Your API Secret
docubot := docubotlib.NewClient(
    docubotURL,
    key,
    secret,
)
thread := "" // The thread that this message is being sent from
sender := "" // The id of the sender that this message is being sent from
message := "" // The message text to send docubot

response, err := docubot.SendMessage(message, thread, sender)
if err != nil {
	// Handle Error
}
if response.Data.Complete {
    urlResponse, err := docubot.GetDocubotDocURL(
		thread,
		sender,
		24*time.Hour,
	)
	if err != nil {
		// Handle Error
	}
    // Deal with Document URL {urlResponse.Data.URL}
}
// Handle Docubot reponse messages {response.Data.Messages}
```

Reporting bugs
----
We try to fix as many bugs we can. If you find an issue, [let us know here](https://github.com/auxai/docubot-go/issues/new).

Contributions
-------------
Anyone is welcome to contribute to Docubot. Just issue a pull request.

There are various ways you can contribute:

* [Raise an issue](https://github.com/auxai/docubot-go/issues) on GitHub.
* Send us a Pull Request with your bug fixes and/or new features.
* Provide feedback and [suggestions on enhancements](https://github.com/auxia/docubot-go/issues?direction=desc&labels=Enhancement&page=1&sort=created&state=open).
* Provide us with a new document for Docubot to generate.
