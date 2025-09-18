# ScoreMaster v4

## Demonstration &amp; Training Guide

This is a brand new version of ScoreMaster, written from the ground up. It currently implements almost all of the v3 features including all those features actually in use in IBAUK/Benelux rallies.

The purpose of this release is to provide a demonstration of the software, to act as a training facility for rally team members tasked with administering rallies or basic claim judging and, last but by no means least, to expose the software to more widespread testing.

## Explore!

This demo is disposable, you can't do any harm by experimenting with it. You can't fiddle with the photos, they are what they are, but you can add and delete entrants and bonuses; you can vary the rally parameters; you can set up complex category structures and scoring methods. Under the [Here be dragons] option, the [Reset Rally] facility provides options to start from scratch or reload the demo. Please report to Bob anything that doesn't work, anything you don't like, anything that could be better. If you do like it and it does work, you could mention that also.

So what works?

## Bonus claims

This setup does not cater for processing live emailed bonus claims but several such claims are available in the demo database. These can be used to examine normal claims judging processes. Manual claims can also be raised via [Show all claims] in the normal way.

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

## Scoring possibilities

ScoreMaster caters for some very complex methods of building scores. The simplest rallies use only simple, fixed point, bonuses and simple combos but the software also allows question/answer extras, score multipliers, bonus points modification based on groups of bonuses and additional points based on groups of bonuses.

1. Scores are built in the first instance by processing individual bonus claims. The score for the bonus itself, usually a fixed number of points, can be adjusted in several ways: some bonuses have variable points, manually input by the judge; some bonuses might represent a multiple of the last bonus scored.
2. A second optional layer examines whether an uninterrupted sequence of bonus claims within a particular category of bonuses exists.
3. Bonus claim counts by category may update the individual bonus values.
4. Bonus claim counts by category may update provide an additional layer of scoring.
5. The whole score can be multiplied by a factor generated from certain bonuses and/or combos.

## Combos

Combos provide an additional layer of scoring, claimed automatically by the software when underlying bonuses (or alphabetically lower combos) are scored. The value of a combo can vary depending on the number of underlying bonuses scored.

## DNF opportunities

Entrants can be rendered DNF in several interesting ways:-

- any bonus or combo can be made compulsory
- entrants not finishing on time are DNF
- limits can be set for minimum/maximum distance ridden
- a minimum points value can be set
- DNF can be triggered by failure to score enough (or too many) categories
- DNF can be triggered by the ratio of scores between categories

## Lesser penalties

- any bonus or combo can be used as a penalty by setting a negative points value - "odd socks" maybe?
- time penalties can be set as a fixed or per minute value
- distance penalties can be set as a fixed or per mile(km) value

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

## FEATURES NOT IMPLEMENTED YET

Several features of v3 ScoreMaster are not implemented in this version. Such features might appear in future versions as and when they're needed. They're described here in anticipation of questions beginning "What happened to feature X?"

### Multiple Legs

It is possible to establish a rally with multiple distinct legs rather than the common arrangement of a single leg. A multi-leg rally has start and finish times for each leg, might be scored separately, and can have bonuses and/or combos restricted to a single leg.

### Entrant emailing

This feature is used to send customised emails to entrants, with or without small attachments, selected by entrant status and/or flag numbers. Typically used to distribute GPX files and/or Flag numbers at the rally start.

### Starting cohorts

Probably no longer needed. Was useful during Covid but starts now tend to be either a single mass start or with first claim starts.

### Odo check facility

No longer used cos, let's face it, a 20 mile check route really doesn't give anything like accuracy, even when the team remembers to include tenths in the trip reading.

### AI

A fully AI version of ScoreMaster may well appear by the 3rd quarter of 2026. This is an exciting prospect: ScoreMaster AI will be distinguished from regular ScoreMaster by the addition of "AI" after "ScoreMaster" on all screens and printouts. In all other respects, ScoreMaster AI will be exactly the same as ScoreMaster.


