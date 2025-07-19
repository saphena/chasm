package main

import (
	"fmt"
	"net/http"
	"strings"
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
	FinishPosition       int
	TotalPoints          int
}

func list_entrants(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Entrants")
	mk := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mk = CS.UnitKmsLit
	}

	sqlx := `SELECT EntrantID,ifnull(RiderLast,''),ifnull(RiderFirst,''),ifnull(PillionLast,''),ifnull(PillionFirst,'')
			,ifnull(Bike,''),EntrantStatus,ifnull(FinishPosition,0),ifnull(TotalPoints,0),ifnull(CorrectedMiles,0)
			FROM entrants`

	seqs := []string{"EntrantID", "RiderLast,RiderFirst", "PillionLast,PillionFirst", "Bike", "EntrantStatus", "Finishposition", "TotalPoints", "CorrectedMiles"}
	hdrs := []string{"Flag", "Name", "Pillion", "Bike", "Status", "Rank", "Points", mk}
	ord := r.FormValue("o")
	if ord == "" {
		ord = seqs[0]
	}
	thisorder := ord
	desc := r.FormValue("d") == "d"
	if desc {
		xx := strings.Split(thisorder, ",")
		for i := 0; i < len(xx); i++ {
			xx[i] += " DESC"
		}
		thisorder = strings.Join(xx, ",")
	}

	sqlx += " ORDER BY " + thisorder
	//fmt.Println(sqlx)
	fmt.Fprint(w, `<div class="entrantlist hdr">`)
	fmt.Fprint(w, `<button class="plus" onclick="addNewEntrant()()" title="Add new entrant">+</button>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="entrantlist hdr"><div class="row hdr">`)

	for ix, x := range seqs {
		urlx := "/entrants?o=" + x
		if ord == x && !desc {
			urlx += "&d=d"
		}

		xclass := ""
		if ix+1 >= len(hdrs) {
			xclass = "num"
		}
		fmt.Fprintf(w, `<span class="col hdr %v" title="Click to reorder"><a href="%v">%v</a></span>`, xclass, urlx, hdrs[ix])

	}
	fmt.Fprint(w, `</div></div><hr>`)
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="entrantlist">`)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var e EntrantDBRecord
		err = rows.Scan(&e.EntrantID, &e.RiderLast, &e.RiderFirst, &e.PillionLast, &e.PillionFirst, &e.Bike, &e.EntrantStatus, &e.FinishPosition, &e.TotalPoints, &e.CorrectedMiles)
		checkerr(err)
		fmt.Fprint(w, `<div class="row">`)

		fmt.Fprintf(w, `<span class="col mid">%v</span>`, e.EntrantID)
		fmt.Fprintf(w, `<span class="col"><strong>%v</strong>, %v</span>`, e.RiderLast, e.RiderFirst)
		pillion := ""
		if e.PillionLast != "" {
			pillion = e.PillionLast + ", " + e.PillionFirst
		}
		fmt.Fprintf(w, `<span class="col">%v</span>`, pillion)
		fmt.Fprintf(w, `<span class="col">%v</span>`, e.Bike)
		fmt.Fprintf(w, `<span class="col">%v</span>`, EntrantStatusLits[e.EntrantStatus])
		rank := ""
		if e.EntrantStatus == EntrantFinisher {
			rank = fmt.Sprintf("%v", e.FinishPosition)
		}
		fmt.Fprintf(w, `<span class="col mid">%v</span>`, rank)
		fmt.Fprintf(w, `<span class="col num">%v</span>`, e.TotalPoints)
		fmt.Fprintf(w, `<span class="col num">%v</span>`, e.CorrectedMiles)

		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)
}
