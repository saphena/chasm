## Table config

This contains a single record holding values reflecting/controlling the
configuration of this event.

Field       | Datatype  | Notes
---         | ---       | ---
dbversion   | integer   | Readonly. The database schema version used
mileskms    | integer   | 0 = Distance unit is miles; 1 is kilometres
langcode    | text      | Language code, default is 'en'
eventname   | text      | Title of the rally
dbinitialised | integer | 0 = uninitialised; 1 = initialised
defaultstart | text     | / = Home/main menu; sc = Scorecards; cl = Claimslog


## Table reasons

This holds details of reasons for bonus claim rejection/modification. There should
only be a small, say 9, number of entries as it is likely to be presented as a
quickpick feature.

Field       | Datatype  | Notes
---         | ---       | ---
code        | integer   | Unique identifier
briefdesc   | text      | Displayed reason text
action      | integer   | 0 = reject claim; 1 = discount percentage
param       | text      | Action dependent. 1 = %

## Table bonuses

This holds records for each ordinary bonus.

Field       | Datatype  | Notes
---         | ---       | ---
bonusid     | text      | Unique identifier
briefdesc   | text      | Bonus description
points      | integer   | Basic points value of this bonus
cat1        | integer   | Analysis category 1
cat2        | integer   | Analysis category 2
cat3        | integer   | Analysis category 3


## Table claims

This holds records for each bonus claim processed.

Field       | Datatype  | Notes
---         | ---       | ---
claimid     | integer   | Unique identifier
entrantid   | integer   | Entrant identifier
leg         | integer   | Number of rally leg
bonusid     | text      | Bonus identifier
odo         | integer   | Odometer reading
claimtime   | text      | Time of claim
decision    | integer   | -1 = undecided, 0 = good claim, 1.. reason for full/partial rejection
points      | integer   | Points value of decided claim

## Table combos

This holds records for each simple combination bonus. 

Field       | Datatype  | Notes
---         | ---       | ---
comboid     | text      | Unique identifier
briefdesc   | text      | Bonus description
minmatch    | integer   | Minimum number of underlying bonuses matched to score; 0 = all
method      | integer   | Scoring method: 0 = any *minmatch* bonuses; 1 = first *minmatch*
bonuses     | text      | Comma separated list of bonus/special/combo IDs
pointvalues | text      | Comma separated list of values

## Table entrants

This holds the main summary of each entrant.

Field       | Datatype  | Notes
---         | ---       | ---
entrantid   | integer   | Unique identifier aka Flag
bike        | text      | Make & model
bikereg     | text      | Registration plate
odokms      | text      | M=miles, K=kilometres
ridername   | text      | Full name of rider (deprecated)
riderfirst  | text      | Rider's first name
riderlast   | text      | Rider's last name
rideriba    | text      | Rider's IBA membership number
pillionname | text      | Full name of pillion (deprecated)
pillionfirst| text      | Pillion's first name
pillionlast | text      | Pillion's last name
pillioniba  | text      | Pillion's IBA membership number
teamid      | integer   | 0 = not team member
country     | text      | Country of origin
odocheckstart | integer | Odo reading during Check-out
odocheckfinish | integer| Odo reading during Check-in
odochecktrip| float     | Deprecated
odoscalefactor | float  | 1.0
correctedmiles | integer| Distance ridden during rally 
starttime   | timestamp | Rally start, check-out or first claim time
finishtime  | timestamp | Check-in or final claim time
bonusesvisited | text   | Deprecated
combosticked   | text   | Deprecated
rejectedclaims | text   | Deprecated
totalpoints | integer   | Final points score
finishposition | integer| Finisher rank
scoringnow  | integer   | Deprecated
scoredby    | text      | Deprecated
extradata   | text      | Key=value pairs
class       | integer   | Certificate class
scorex      | text      | Detailed scorecard
phone       | text      | Rider's phone
email       | text      | Comma-separated list of authorised email addresses
nokname     | text      | Name of emergency contact
nokphone    | text      | Emergency contact's phone
nokrelation | text      | How is NoK related to rider
bcmethod    | integer   | 1 = emailed bonus claims
restminutes | integer   | Number of minutes rest accrued
confirmed   | integer   | 0 - unknown usage
cohort      | integer   | 0 - future use
lastreviewed| timestamp | Time of last rally team review
reviewstatus| integer   | 0 = not reviewed; 1 = team happy; 2 = team NOT happy; 3 = entrant happy





## Datatypes

### Class

A class represents a subset of entrants. Class is held as an integer, default value 1.

Class can be assigned manually where, for example, entrants choose a particular class such as
a particular route or vehicle type. Class can also be assigned automatically based on score components.

### Entrant status

Each entrant record has a code representing the current status of the entrant. Status is held as an integer, default value 1.

Code    | Status        | Notes
---     | ---           | ---
1       | ok            | Neutral, uncommitted, value
2       | DNS           | Did Not Start. Only used once Rally has commenced
3       | DNF           | Did Not Finish
4       | DSQ           | Disqualified
5       | Finisher      | Qualified finisher

