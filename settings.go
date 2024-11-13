package main

type emailSettings struct {
	SMTP struct {
		Host          string
		Port          string
		UseInboxCreds bool
		Userid        string
		Password      string
		CertName      string // May need to override the certificate name used for TLS
	}
	IMAP struct {
		HostPort string
		Userid   string
		Password string
	}
}
type chasmSettings struct {
	ShowExcludedClaims  bool // If a claim is marked 'excluded' and is not superseded, show it on the scoresheet
	CurrentLeg          int
	UseCheckinForOdo    bool // If true, OdoRallyFinish updated only by check-in, not by individual claims
	RallyUnitKms        bool // Report in Kms(true) or Miles(false)
	UnitMilesLit        string
	UnitKmsLit          string
	PenaltyMilesDNF     int
	RallyMinMiles       int
	DebugRules          bool
	AutoLateDNF         bool
	RallyTitle          string
	RallyMinPoints      int
	RallyTimezone       string
	RallyUseQA          bool
	FlagTeamTitle       string
	FlagAlertTitle      string
	FlagBikeTitle       string
	FlagDaylightTitle   string
	FlagFaceTitle       string
	FlagNightTitle      string
	FlagRestrictedTitle string
	FlagReceiptTitle    string
	CloseEBCUndecided   string
	CloseEBC            []string
	ImgBonusFolder      string // Holds rally book bonus photos
	ImgEbcFolder        string // Holds images captured from emails
	RallyBookImgTitle   string
	EBCImgTitle         string
	EBCImgSwapTitle     string
	Email               emailSettings
}

var CS chasmSettings

const defaultCS = `{
	"ShowExcludedClaims": 	true,
	"CurrentLeg": 			5,
	"UseCheckInForOdo": 	false,
	"RallyUnitKms": 		false,
	"UnitMilesLit":			"miles",
	"UnitKmsLit":			"km",
	"PenaltyMilesDNF":		99999,
	"RallyTitle":			"Brit Butt Rally 2025",
	"RallyMinMiles":		0,
	"DebugRules":			true,
	"AutoLateDNF": 			true,
	"RallyMinPoints":		-99999,
	"RallyTimezone":		"Europe/London",
	"RallyUseQA":			false,
	"FlagTeamTitle":       	"Team rules",
	"FlagAlertTitle":      	"Read the notes!",
	"FlagBikeTitle":       	"Bike in photo",
	"FlagDaylightTitle":   	"Daylight only",
	"FlagFaceTitle":       	"Face in photo",
	"FlagNightTitle":      	"Night only",
	"FlagRestrictedTitle": 	"Restricted access",
	"FlagReceiptTitle":		"Need a receipt (ticket)",
	"CloseEBCUndecided":	"Leave undecided",
	"CloseEBC":				["Accept good claim","No photo","Wrong/unclear photo","Out of hours/disallowed","Face not in photo","Bike not in photo","Flag not in photo","Missing rider/pillion","Missing receipt","Claim excluded" ],
	"ImgBonusFolder":		"images/bonuses/",
	"ImgEbcFolder":			"images/ebcimg/",
	"RallyBookImgTitle":	"Rally book photo",
	"EBCImgTitle":			"Entrant's image - click to resize",
	"EBCImgSwapTitle":		"Click to view this image",
	"Rally":				{"A1":"AAAAAAAAAAAAAA","A2":"22222222222222"}
}`

const secondDefault = `{
"CloseEBCUndecided":		"Not on your nellie",
"RallyUseQA":				false
}`
