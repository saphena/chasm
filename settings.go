package main

type chasmSettings struct {
	ShowExcludedClaims bool // If a claim is marked 'excluded' and is not superseded, show it on the scoresheet
	CurrentLeg         int
	UseCheckinForOdo   bool // If true, OdoRallyFinish updated only by check-in, not by individual claims
	RallyUnitKms       bool // Report in Kms(true) or Miles(false)
}

var CS chasmSettings

const defaultCS = `{
	"ShowExcludedClaims": 	true,
	"CurrentLeg": 			5,
	"UseCheckInForOdo": 	false,
	"RallyUnitKms": 		false
}`