## Table config

This contains a single record holding values reflecting/controlling the
configuration of this event.

Field       | Datatype  | Notes
---         | ---       | ---
dbversion   | integer   | Readonly. The database schema version used
mileskms    | integer   | 0 = Distance unit is miles; 1 is kilometres
langcode    | text      | Language code, default is 'en'
eventname   | text      | Title of the rally


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

