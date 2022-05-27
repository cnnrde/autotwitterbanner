# autotwitterbanner

## Why?
I was bored and thought it'd be cool

## How can I use this?
1. Create twitter developer app
2. Install golang
3. Build with `go build`
4. Make .env file following the .env.example template thingy and place your twitter tokens from developer dashboard in the file
5. Set up crontab to run every 15 minutes or so
```
# should look something like this
*/15 * * * * ~/autotwitterbanner/autobannerupdater > ~/autotwitterbanner/log.txt 2>&1
```
6. Hope it works

## why code bad
i wrote this in like a day when i was ill

## footnote
thank you for viewing my repo have a good day