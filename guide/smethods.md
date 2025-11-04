# ScoreMaster 4 - Scoring possibilities

ScoreMaster caters for a wide range of methods of building scores. The simplest rallies use only simple, fixed point, bonuses and simple combos but the software also allows question/answer extras, score multipliers, bonus points modification based on groups of bonuses and additional points based on groups of bonuses.

## Overview
Entrants simply make individual bonus claims throughout the rally and ScoreMaster takes care of the resultant scoring calculations. It does so by following a sequence of levels as follows:-

1. Does this bonus extend an uninterrupted sequence? If not, apply accrued sequence value
2. Calculate value of the current bonus and add any question/answer value
3. Apply bonus-modifying complex rules
4. Apply combo scores
5. Apply category-based complex rules
6. Apply time penalties
7. Apply distance penalties
8. Apply final multipliers 

***
## Some definitions
**Points**
: A points value is an integer (whole number), fractions are not allowed. In general a points value can be positive or negative.

**Bonus code**
: A unique identifier comprising letters and digits only. Spaces and punctuation are forbidden. All letters must be uppercase. 

As far as the system is concerned, bonus codes are alphanumeric even when they are all numeric. In order to keep them presented in the "correct" sequence they should all contain the same number of digits. Use leading '0's to pad the codes or start from, say, 11 or 111 instead of 1.

It is not necessary (or perhaps even desirable) that bonuses should be coded with some sequence and you certainly don't need to include any kind of serial number. Each bonus is unique as far as the system is concerned and each code stands alone. The only sequence that matters is the way they're presented on scorecards and in maintenance lists. Use your imagination: codes should be short (entrants will need to type them) but can be anything (uppercase) that you choose.

**Combo code**
: A unique identifier as with an ordinary bonus code. Letters may be upper or lowercase but AA23 is the same as aa23, Aa23 and aA23.

---

## Simple bonus points
Each bonus scores a specific number of points.

## Variable bonus points
Each bonus specifies a default value (which might be zero or even negative). At claim time the actual value of the bonus must be entered manually. This might be used for example when a bonus includes a clock face or a varying number of flags. This facility should be used sparingly as it has a significant impact on scoring efficiency.

## Last Bonus multiplier
The points value of this bonus is used as a multiplier applied to the value of the most recently claimed bonus unless that bonus was itself a multiplier in which case the value is zero.

## Simple combination bonus
A combination bonus (combo) specifies a list of underlying bonuses (which might include some* other combos). Each combo scores a specific, fixed, number of points (or multipliers, see below). This is in addition to the scores of the underlying bonuses. Combos are claimed automatically, entrants don't need to claim them separately.

*Any combos included as underlying bonuses must be coded alphanumerically lower than the current combo code. Combo "A12" cannot depend on combo "A13" or "B1".

## Variable combination bonus
A variable combo specifies a minimum number of underlying bonuses (rather than 'all' for a simple combo). Each combo scores a specific number of points or multipliers depending on the number of underlying bonuses scored between the minimum and all. A simple example:-

Combo *CLUBS* has underlying bonuses *MANU*,*CHELSEA*,*ARSENAL*,*RANGERS*,*ROVERS*. The combo will be scored if at least three underlying bonuses are scored. The value of the combo is 100, 200 or 600 points. 

## Multipliers
Combos can be used to build a final score multiplier which, as the names suggest, will be used to multiply the final score. This can act as a major incentive to score certain combos. Multipliers can also be accrued using complex methods below.

## Questions & answers
Ordinary bonuses may have associated questions which allow for the award of additional points when correctly answered. The question/answer will be worth a fixed number of points, the value being set for the rally as a whole rather than for individual bonuses. The extra points are only awarded if the underlying bonus claim succeeds.

---

## Complex methods

Powerful complex scoring methods are available based on sets of categories applied to bonuses and/or combos. [Details here](/guide/complex).


---

## Penalties

A variety of penalties can be specified either limiting points scored or imposing DNF status. [See here](/guide/penalty)