# Navigating around Chasm

Home
- Setup wizard
- Default
-- Front menu
-- Last page
-- Scorecards
-- Claims log

if !dbinitialised {
    Setup wizard
} else {
    if default == claimslog {
        show claims log
    } else if default == scorecards {
        show scorecards
    } else if authneeded {
        login
        show front menu
    } else {
        show front menu
    }
}