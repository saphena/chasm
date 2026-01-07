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

Where check-out and check-in are used, those readings are the authoritative ones used throughout the system, regardless of any individual claim readings.
Individual claim readings are used where check-out/in are not used and are also applied as interim values during the rally. Such readings need not be actual
odo readings, they could be simple sequence numbers or even a constant unless the rally relies on them to report distances ridden.

In any event, a reading recorded during check-in is always recorded as the final reading.

## Teams

Claims can be submitted from any email address associated with a team member but MUST always use the same entrant number if score cloning is in force.