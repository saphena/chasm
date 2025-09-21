# chasm
A scoring system for navigational scatter rallies


## Compound rules

### Affecting individual bonus scores
- Creating a modify bonus rule with a power of 0 results in a value of bonusPoints x (number of hits - 1).
- Creating a modify bonus rule with a power of 1 results in a value of bonusPoints.
- Creating a modify bonus rule with a power > 1 results in a value of bonusPoints x power raised to the power of (number of hits - 1).

Examples:- 
| Bonus points | Power | Num hits | Result
| --- | --- | --- | ---
| 1 | 0 | 5 | 4
| 1 | 2 | 5 | 32
| 2 | 2 | 5 | 64

If the result of the rule is multiplier rather than points, the formula is BP x Power x (number of hits -1)


## Compound vs Combo
There is a large amount of crossover in functionality between combos and compound scoring rules.

| Feature | Compound | Combo
| --- | --- | ---
| Affect bonus score | Y | N
| Affect group score | Y | Y
| Affect sequence score | Y | N
| DNF if triggered | Y | N
| DNF unless triggered | Y | Y

## Odo readings

There is an option UseCheckinForOdo which, if true, prevents individual claims from updating final odo readings. This is so that oodo readings can be simple serial numbers rather than actual odo readings. If distance is to be recorded in a rally at all, it comes either from individual claim readings or from the final check-in.

In any event, a reading recorded during check-in is always recorded as the final reading.

## Teams

Claims can be submitted from any email address associated with a team member but MUST always use the same entrant number if score cloning is in force.