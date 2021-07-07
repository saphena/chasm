package main

const chasmSQL = `BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "config" (
	"dbversion"	INTEGER NOT NULL DEFAULT 1,
	"mileskms"	INTEGER NOT NULL DEFAULT 0,
	"langcode"	TEXT NOT NULL DEFAULT 'en',
	"eventname"	TEXT NOT NULL,
	"dbinitialised" TEXT NOT NULL DEFAULT 0,
	"defaultstart" TEXT NOT NULL DEFAULT '/',
	"region" TEXT NOT NULL DEFAULT 'United Kingdom',
	"localtz" TEXT NOT NULL DEFAULT 'Europe/London',
	"decimalcomma" INTEGER NOT NULL DEFAULT 0,
	"hostcountry" TEXT NOT NULL DEFAULT 'UK',
	"locale" TEXT NOT NULL DEFAULT 'en-GB',
	"virtualrally" INTEGER NOT NULL DEFAULT 0
);
INSERT INTO config (EVENTNAME) VALUES('IBA Rally');
CREATE TABLE IF NOT EXISTS "regions" (
	"region" TEXT NOT NULL,
	"localtz" TEXT NOT NULL,
	"hostcountry" TEXT NOT NULL,
	"locale" TEXT NOT NULL,
	"mileskms" INTEGER NOT NULL DEFAULT 0,
	"decimalcomma" INTEGER NOT NULL DEFAULT 0
);
INSERT INTO regions (region,localtz,hostcountry,locale,mileskms,decimalcomma) VALUES('United Kingdom','Europe/London','UK','en-GB',0,0);
INSERT INTO regions (region,localtz,hostcountry,locale,mileskms,decimalcomma) VALUES('Republic of Ireland','Europe/Dublin','Eire','en-EI',1,0);
INSERT INTO regions (region,localtz,hostcountry,locale,mileskms,decimalcomma) VALUES('Western Europe','Europe/Berlin','DE','de-DE',1,1);
INSERT INTO regions (region,localtz,hostcountry,locale,mileskms,decimalcomma) VALUES('Eastern Europe','Europe/Helsinki','Finland','fi-FI',1,1);
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
