# ScoreMaster v4 - Complex arithmetic

When applying complex rules to individual bonus scores (as opposed to a group score) the following is possible, starting with some definitions:-



**BV** is the points value of current bonus

**RV** is the *results in* value of current rule

**N** is the number of bonuses scored within the category

**M** is **N** - 1 

**SV** is the resulting score for the bonus


## Formulas

If **RV** is 0, **SV** = **BV** * **N**  simple multiplication.

If **RV** is set to "multipliers", **SV** = **BV** * **RV** * **N**  simple multiplication.

If **RV** is set to "points", **SV** = **BV** * **RV** ^ **M** exponential score.


## Examples

So, with **BV** = 5 for all bonuses claimed, **RV** = 2, points. Successive claims give:-

1. **SV** = 5 * 2 ^ 0 = 5
2. **SV** = 5 * 2 ^ 1 = 10
3. **SV** = 5 * 2 ^ 2 = 20
4. **SV** = 5 * 2 ^ 3 = 40

With **BV** = 5, **RV** = 2, multipliers

1. **SV** = 5 * 2 * 1 = 10
2. **SV** = 5 * 2 * 2 = 20
3. **SV** = 5 * 2 * 3 = 30
4. **SV** = 5 * 2 * 4 = 40

With **BV** = 5, **RV** = 3, points

1. **SV** = 5 * 3 ^ 0 = 5
2. **SV** = 5 * 3 ^ 1 = 15
3. **SV** = 5 * 3 ^ 2 = 45
4. **SV** = 5 * 3 ^ 3 = 135

With **BV** = 5, **RV** = 3, multipliers

1. **SV** = 5 * 3 * 1 = 15
2. **SV** = 5 * 3 * 2 = 30
3. **SV** = 5 * 3 * 3 = 45
4. **SV** = 5 * 3 * 4 = 60

