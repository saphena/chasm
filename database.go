package main

import (
	"fmt"
	"reflect"
)

type REASON struct {
	Code      int
	Briefdesc string
	Action    int
	Param     string
}

type RALLYCONFIG struct {
	DBVersion         int
	MilesKms          int
	Langcode          string
	Eventname         string
	DBInitialised     int
	DefaultStart      string
	Region            string
	Localtz           string
	Decimalcomma      int
	Hostcountry       string
	Locale            string
	VirtualRally      int
	Startdate         string
	Starttime         string
	Finishdate        string
	Finishtime        string
	Maxhours          int
	MaxmilesDNF       int
	MaxmilesPenalties int
	MinmilesDNF       int
	ExcessPPM         int
	Tankrange         int
	MinsPerStop       int
}

func fetchReasons() *[]REASON {

	var cfg []REASON

	cfg = make([]REASON, 0)
	rows, err := DBH.Query("SELECT * FROM reasons ORDER BY Code")
	if err != nil {
		fmt.Printf("%v\n", err)
		return &cfg
	}
	defer rows.Close()

	for rows.Next() {
		var r REASON
		s := reflect.ValueOf(&r).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}
		rows.Scan(columns...)
		cfg = append(cfg, r)
	}

	return &cfg

}

func fetchConfig(cfg *RALLYCONFIG) {

	rows, _ := DBH.Query("SELECT * FROM config")

	defer rows.Close()
	rows.Next()

	s := reflect.ValueOf(cfg).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)

	if err != nil {
		fmt.Printf("Scan ! %v\n", err)
	}
}

const chasmSQL = `BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "config" (
	"DBVersion"	INTEGER NOT NULL DEFAULT 1,
	"MilesKms"	INTEGER NOT NULL DEFAULT 0,
	"Langcode"	TEXT NOT NULL DEFAULT 'en',
	"Eventname"	TEXT NOT NULL,
	"DBInitialised" INTEGER NOT NULL DEFAULT 0,
	"DefaultStart" TEXT NOT NULL DEFAULT '/',
	"Region" TEXT NOT NULL DEFAULT 'United Kingdom',
	"Localtz" TEXT NOT NULL DEFAULT 'Europe/London',
	"Decimalcomma" INTEGER NOT NULL DEFAULT 0,
	"Hostcountry" TEXT NOT NULL DEFAULT 'UK',
	"Locale" TEXT NOT NULL DEFAULT 'en-GB',
	"VirtualRally" INTEGER NOT NULL DEFAULT 0,
	"Startdate" TEXT NOT NULL DEFAULT '2021-01-01',
	"Starttime" TEXT NOT NULL DEFAULT '09:00',
	"Finishdate" TEXT NOT NULL DEFAULT '2021-01-01',
	"Finishtime" TEXT NOT NULL DEFAULT '17:00',
	"Maxhours" INTEGER NOT NULL DEFAULT 8,
	"MaxmilesDNF" INTEGER NOT NULL DEFAULT 0,
	"MaxmilesPenalties" INTEGER NOT NULL DEFAULT 0,
	"MinmilesDNF" INTEGER NOT NULL DEFAULT 0,
	"ExcessPPM" INTEGER NOT NULL DEFAULT 0,
	"Tankrange" INTEGER NOT NULL DEFAULT 200,
	"MinsPerStop" INTEGER NOT NULL DEFAULT 10
);
INSERT INTO config (EVENTNAME) VALUES('IBA Rally');
CREATE TABLE IF NOT EXISTS "regions" (
	"Region" TEXT NOT NULL,
	"Localtz" TEXT NOT NULL,
	"Hostcountry" TEXT NOT NULL,
	"Locale" TEXT NOT NULL,
	"MilesKms" INTEGER NOT NULL DEFAULT 0,
	"Decimalcomma" INTEGER NOT NULL DEFAULT 0
);
INSERT INTO regions (Region,Localtz,Hostcountry,Locale,MilesKms,Decimalcomma) VALUES('United Kingdom','Europe/London','UK','en-GB',0,0);
INSERT INTO regions (Region,Localtz,Hostcountry,Locale,MilesKms,Decimalcomma) VALUES('Republic of Ireland','Europe/Dublin','Eire','en-EI',1,0);
INSERT INTO regions (Region,Localtz,Hostcountry,Locale,MilesKms,Decimalcomma) VALUES('Western Europe','Europe/Berlin','Germany','de-DE',1,1);
INSERT INTO regions (Region,Localtz,Hostcountry,Locale,MilesKms,Decimalcomma) VALUES('Eastern Europe','Europe/Helsinki','Finland','fi-FI',1,1);
CREATE TABLE IF NOT EXISTS "entrants" (
	"Entrantid"	INTEGER NOT NULL,
	"Riderid"	INTEGER NOT NULL,
	"Pillionid"	INTEGER,
	"Nokid"	INTEGER,
	"Status"	INTEGER NOT NULL DEFAULT 0,
	"Class"	INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY("Entrantid")
);
CREATE TABLE IF NOT EXISTS "classes" (
	"Class"	INTEGER NOT NULL,
	"Briefdesc"	TEXT NOT NULL,
	"Certificate"	TEXT,
	PRIMARY KEY("Class")
);
CREATE TABLE IF NOT EXISTS "specials" (
	"specialid"	TEXT NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"method"	INTEGER NOT NULL DEFAULT 0,
	"points"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("specialid")
);
CREATE TABLE IF NOT EXISTS "combos" (
	"Comboid"	TEXT NOT NULL,
	"Briefdesc"	TEXT NOT NULL,
	"Minmatch"	INTEGER NOT NULL DEFAULT 0,
	"Method"	INTEGER NOT NULL DEFAULT 0,
	"Bonuses"	TEXT NOT NULL,
	"Pointvalues"	TEXT NOT NULL,
	PRIMARY KEY("Comboid")
);
CREATE TABLE IF NOT EXISTS "reasons" (
	"Code"	INTEGER NOT NULL,
	"Briefdesc"	TEXT NOT NULL,
	"Action"	INTEGER NOT NULL DEFAULT 0,
	"Param"	TEXT,
	PRIMARY KEY("Code")
);
INSERT INTO reasons (Code,Briefdesc,Action,Param) VALUES(1,'Face not in photo',0,'');
INSERT INTO reasons (Code,Briefdesc,Action,Param) VALUES(2,'Flag not in photo',0,'');
INSERT INTO reasons (Code,Briefdesc,Action,Param) VALUES(3,'Bike not in photo',0,'');
INSERT INTO reasons (Code,Briefdesc,Action,Param) VALUES(4,'Wrong/missing photo',0,'');
INSERT INTO reasons (Code,Briefdesc,Action,Param) VALUES(5,'Out of hours',0,'');
CREATE TABLE IF NOT EXISTS "people" (
	"Personid"	INTEGER NOT NULL,
	"Firstname"	TEXT,
	"Lastname"	TEXT NOT NULL,
	"Memberid"	TEXT,
	"Phone"	TEXT,
	"Email"	TEXT,
	"Address1"	TEXT,
	"Address2"	TEXT,
	"Towncity"	TEXT,
	"County"	TEXT,
	"Postcode"	TEXT,
	"Country"	TEXT,
	PRIMARY KEY("Personid")
);
CREATE TABLE IF NOT EXISTS "catcats" (
	"Axis"	INTEGER NOT NULL,
	"Cat"	INTEGER NOT NULL,
	"Briefdesc"	TEXT NOT NULL,
	PRIMARY KEY("Axis","Cat")
);
CREATE TABLE IF NOT EXISTS "cataxes" (
	"Axis"	INTEGER NOT NULL,
	"Briefdesc"	TEXT,
	PRIMARY KEY("Axis")
);
CREATE TABLE IF NOT EXISTS "bonuses" (
	"Bonusid"	TEXT NOT NULL,
	"Briefdesc"	TEXT NOT NULL,
	"Points"	INTEGER NOT NULL DEFAULT 0,
	"Cat1"	INTEGER NOT NULL DEFAULT 0,
	"Cat2"	INTEGER NOT NULL DEFAULT 0,
	"Cat3"	INTEGER NOT NULL DEFAULT 0,
	"Compulsory" INTEGER NOT NULL DEFAULT 0,
	"RestMinutes" INTEGER NOT NULL DEFAULT 0,
	"AskPoints" INTEGER NOT NULL DEFAULT 0,
	"AskMinutes" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("Bonusid")
);
COMMIT;
`
