
Variable | Type | Default | Description
--- | --- | --- | ---
UseEBC | bool | true | false = new claims entered manually; true = new claims received via EBC
StartOption | int | 0 | 0 = Rally starts with check-out; 1 = Rally starts with first bonus claim
AutoFinisher | bool | false | false = finish at check-in; true = finish with final claim
ShowExcludedClaims | bool | false | false = don't show excluded claims on scorecard; true = do show
SuppressExclusion | bool | true | false = offer 'Exclude' when processing new claims; false = don't offer
PenaltyMilesDNF | int | 99999 | Rally miles/kms exceeding this triggers DNF
PenaltyMilesMax | int | 99999 | Rally miles/kms exceeding this triggers distance penalties
PenaltyMilesMethod | int | 0 | 0 = Fixed number of points; 1 = points per mile/km; 2 = Number of multipliers
PenaltyMilesPoints | int | 0 | The value used for the mileage penalty
RallyMinMiles | int | 0 | Rally miles/kms below this triggers DNF
AutoLateDNF | bool | true | Automatically apply DNF for late finishing
RallyMinPoints | int | -99999 | Scores lower than this value are DNF
RallyUseQA | bool | false | Offer question/answer bonuses
RallyQAPoints | int | 50 | Value of validly claimed QA bonuses
RallyUsePctPen | bool | false | Offer claim acceptance with minor penalty
RallyPctPenVal | int | 10 | Points percentage charged as penalty
RallyRankEfficiency | bool | false | false = rank finishers by points only; true = rank by points per mile/km
RallySplitTies | bool | true | false = tied finishers receive same rank; true = tied finisher with lower miles/km ranked higher
RallyTeamMethod | int | 3 | 0 = team members ranked individually; 1 = rank as highest member; 2 = rank as lowest member; 3 = members cloned
DowngradedClaimDecision | int | 3 | Decision value when claim is downgraded for balancing purposes
RallyTitle | text | * | The title of the rally for screens and certificates
RallyStarttime | datetime | * | The date and time of the earliest possible rider start
RallyFinishtime | datetime | * | The date and time of the latest possible rider finish
RallyMaxHours | int | 12 | The maximum number of hours a until a rider's rally finish
RallyUnitKms | bool | false | false = unit of distance is miles; true = kilometres
RallyTimezone | timezone | Europe/London | The registered name of the relevant timezone
RallyPointIscomma | bool | false | false = decimal point is '.'; true = decimal point is ','
UnitMilesLit | text | miles | Text used to describe distances in miles
UnitKmsLit | text | km | Text used to describe distances in kilometres
DebugRules | bool | false | Show debugging info when applying complex rules
CurrentLeg | int | 1 | Number of current leg during multi-leg rallies

- bool : true or false
- int : an integer with/without a sign
- text : a quoted string of characters
- datetime : a quoted string in the format 'yyyy-mm-ddThh:mm'
- timezone : a quoted string from the range of registered timezones
