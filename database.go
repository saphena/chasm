package main

const chasmSQL = `BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "config" (
	"DBVersion"	INTEGER NOT NULL DEFAULT 1,
	"MilesKms"	INTEGER NOT NULL DEFAULT 0,
	"Langcode"	TEXT NOT NULL DEFAULT 'en',
	"Eventname"	TEXT NOT NULL,
	"DBInitialised" TEXT NOT NULL DEFAULT 0,
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
	"entrantid"	INTEGER NOT NULL,
	"riderid"	INTEGER NOT NULL,
	"pillionid"	INTEGER,
	"nokid"	INTEGER,
	"status"	INTEGER NOT NULL DEFAULT 0,
	"class"	INTEGER NOT NULL DEFAULT 1,
	PRIMARY KEY("entrantid")
);
CREATE TABLE IF NOT EXISTS "classes" (
	"class"	INTEGER NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"certificate"	TEXT,
	PRIMARY KEY("class")
);
CREATE TABLE IF NOT EXISTS "specials" (
	"specialid"	TEXT NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"method"	INTEGER NOT NULL DEFAULT 0,
	"points"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("specialid")
);
CREATE TABLE IF NOT EXISTS "combos" (
	"comboid"	TEXT NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"minmatch"	INTEGER NOT NULL DEFAULT 0,
	"method"	INTEGER NOT NULL DEFAULT 0,
	"bonuses"	TEXT NOT NULL,
	"pointvalues"	TEXT NOT NULL,
	PRIMARY KEY("comboid")
);
CREATE TABLE IF NOT EXISTS "reasons" (
	"code"	INTEGER NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"action"	INTEGER NOT NULL DEFAULT 0,
	"param"	TEXT,
	PRIMARY KEY("code")
);
CREATE TABLE IF NOT EXISTS "people" (
	"personid"	INTEGER NOT NULL,
	"firstname"	TEXT,
	"lastname"	TEXT NOT NULL,
	"memberid"	TEXT,
	"phone"	TEXT,
	"email"	TEXT,
	"address1"	TEXT,
	"address2"	TEXT,
	"towncity"	TEXT,
	"county"	TEXT,
	"postcode"	TEXT,
	"country"	TEXT,
	PRIMARY KEY("personid")
);
CREATE TABLE IF NOT EXISTS "catcats" (
	"axis"	INTEGER NOT NULL,
	"cat"	INTEGER NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	PRIMARY KEY("axis","cat")
);
CREATE TABLE IF NOT EXISTS "cataxes" (
	"axis"	INTEGER NOT NULL,
	"briefdesc"	TEXT,
	PRIMARY KEY("axis")
);
CREATE TABLE IF NOT EXISTS "bonuses" (
	"bonusid"	TEXT NOT NULL,
	"briefdesc"	TEXT NOT NULL,
	"points"	INTEGER NOT NULL DEFAULT 0,
	"cat1"	INTEGER NOT NULL DEFAULT 0,
	"cat2"	INTEGER NOT NULL DEFAULT 0,
	"cat3"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("bonusid")
);
COMMIT;
`
