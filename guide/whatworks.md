# ScoreMaster v4 - What works

## With or without EBC

Most rally teams will make use of Electronic Bonus Claiming using email or some other method but ScoreMaster works equally well with claims being processed in the "old-fashioned" way, entered manually at the end of the rally.

## Scoring options

The short version is that all scoring methods (bonuses, combos, etc, etc) currently in use are handled automatically. Full details are [available here](/guide/smethods). If you can't find the method you're looking for, describe it to Bob and it'll be included in the next release.

## Bonus claims

This setup does not cater for processing live EBC claims but several such claims are available in the demo database. These can be used to examine normal claims judging processes. Manual claims can also be raised via [Show all claims] in the normal way or, if the option <em>Use EBC</em> is set to "Manually enter claims", via the "Judge incoming claims" button. [See also](/guide/claims)

## Certificate classes

In most rallies only one certificate class, Finisher, is used but it is also possible to use multiple classes with different certificates. [See also](/guide/classes)

## Teams

A team consists of at least two bikes. A bike with a rider and passenger is not a team, a team means multiple bikes.

Normally, team members all receive the same score, cloned automatically by the system.


## Rally time

Regardless of the size of the territory, each rally is associated with a single timezone.

The 'Rally starts' and 'Rally finishes' variables set the outer limits of rally time. These set the overall window within which everything must happen. The variable 'Max Rideable hours', normally simply the difference between those two variables, can be used for rallies which allow some flexibility in start times.

## Resetting the demo rally

An option under "Here be dragons" lets you  restore the demo database back to its initial state. Part of that process involves updating all the dates to the current date so that you have a realistic idea of how the system behaves with respect to timing.
