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