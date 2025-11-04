# ScoreMaster v4

## Demonstration &amp; Training Guide

The purpose of this release is to provide a demonstration of the software, to act as a training facility for rally team members tasked with administering rallies or basic claim judging and, last but by no means least, to expose the software to more widespread testing.

## Explore!

This demo is disposable, you can't do any harm by experimenting with it. You can add and delete entrants and bonuses; you can vary the rally parameters; you can set up complex category structures and scoring methods. Under the [Here be dragons] option, the [Reset Rally] facility provides options to start from scratch or reload the demo. Please report to Bob anything that doesn't work, anything you don't like, anything that could be better. If you do like it and it does work, you could mention that also.

So what works?

## With or without EBC

Most rally teams will make use of Electronic Bonus Claiming using email or some other method but ScoreMaster works equally well with claims being processed in the "old-fashioned" way, entered manually at the end of the rally.

## Scoring options

The short version is that all scoring methods (bonuses, combos, etc, etc) currently in use are handled automatically. Full details are [available here](/guide/smethods). If you can't find the method you're looking for, describe it to Bob and it'll be included in the next release.

## Bonus claims

This setup does not cater for processing live EBC claims but several such claims are available in the demo database. These can be used to examine normal claims judging processes. Manual claims can also be raised via [Show all claims] in the normal way or, if the option <em>Use EBC</em> is set to "Manually enter claims", via the "Judge incoming claims" button.

Claims can also be amended or deleted via [Show all claims]. All such activity results in scorecard recalculation for the entrant.

Claims can be decided (judged) in one of four ways:-

- Accept good claim
- Reject claim
- Exclude claim
- Leave undecided

As it suggests, *Accept good claim* awards the score for this claim.
Rejecting a claim, further refined into one of eight reasons, denies the score for this claim.
Excluding a claim, different from rejecting a claim, means that the claim does not affect scoring and will not appear on scorecards.
Leave undecided returns the claim to the bottom of the queue. "I'll look again later" or "I'll let Suzy decide this one".

Normally, claims will be either accepted or rejected. Excluding a claim should be a rare event, in response to special circumstances.


## Certificate classes

In most rallies only one certificate class, Finisher, is used but it is also possible to use multiple classes with different certificates.

Such classes can be assigned manually or automatically. Manual classes might be used for example to distinguish between entrants riding BMWs and those riding Hondas. Automatic classes could be used for example to provide separate certificates for podium finishers.

## Teams

A team consists of at least two bikes. A bike with a rider and passenger is not a team, a team means multiple bikes.

Normally, team members all receive the same score, cloned automatically by the system.

## Entrant status

Each entrant record includes a status, one of DNS, ok, Finisher and DNF. The initial status is DNS (Did Not Start). This changes, depending on the settings 'Rally Start option', during Check-out or at when the first bonus claim is submitted.

## Rally time

Regardless of the size of the territory, each rally is associated with a single timezone.

The 'Rally starts' and 'Rally finishes' variables set the outer limits of rally time. These set the overall window within which everything must happen. The variable 'Max Rideable hours', normally simply the difference between those two variables, can be used for rallies which allow some flexibility in start times.

## Resetting the demo rally

An option under "Here be dragons" lets you  restore the demo database back to its initial state. Part of that process involves updating all the dates to the current date so that you have a realistic idea of how the system behaves with respect to timing.


## FEATURES NOT IMPLEMENTED YET

Some V3 features have not been implemented yet but might be in the future. [See here](/guide/niy)
