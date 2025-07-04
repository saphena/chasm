package main

import (
	"fmt"
	"net/http"
)

type EntrantDBRecord struct {
	EntrantID            int
	Bike                 string
	BikeReg              string
	RiderFirst           string
	RiderLast            string
	RiderCountry         string
	RiderIBA             string
	RiderPhone           string
	RiderEmail           string
	PillionFirst         string
	PillionLast          string
	PillionIBA           string
	PillionPhone         string
	PillionEmail         string
	OdoKms               int
	OdoStart             int
	OdoFinish            string
	CorrectedMiles       string
	FinishTime           string
	StartTime            string
	EntrantStatus        int
	NokName              string
	NokPhone             string
	NokRelation          string
	EntryDonation        string
	SquiresCheque        string
	SquiresCash          string
	RBLRAccount          string
	JustGivingAmt        string
	JustGivingURL        string
	Route                string
	RiderRBL             string
	PillionRBL           string
	Tshirt1              string
	Tshirt2              string
	Patches              int
	FreeCamping          string
	CertificateDelivered string
	CertificateAvailable string
}

func list_entrants(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Entrants")
	fmt.Fprint(w, `<div class="intro">`)
	fmt.Fprint(w, `<button class="plus" onclick="addNewEntrant()()" title="Add new entrant">+</button>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="entrantlist hdr">`)
	fmt.Fprint(w, `<span>Entrant</span><span>Rider</span><span>Bike</span><span>Status</span>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)

}
