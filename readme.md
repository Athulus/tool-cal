# tool- cal
an event calendar API to manage the use of tools at HackRva

[![Build Status](https://travis-ci.org/Athulus/tool-cal.svg?branch=master)](https://travis-ci.org/Athulus/tool-cal)

## building and running
1.  `git clone` tool-cal repo into your gopath/src
2. `cd tool-cal`
3. `make` will run install dependencies, run test, and compile a static binary
4. `make deploy` will build a docker image and run the image and a redis server with docker-compose

## code organization
all code is in the `main` package

[main.go](main.go) has all of the http server code, inluding the  http handler functions

[cal.go](cal.go) has the code that deals with data access of the calendar events from redis

[middleware.go](middleware.go) (currently empty) should have all of the http middleware functions if they are needed

[user.go](user.go) (currently empty) will hold the code to deal with whatever users will do?

## notes for tool scheduling app 
- calendar    
    - calendar for each tool    
    - user can add time to a tool
- users    
    - users have a certian amount of time they can use per month
- tools
    - tool use time estimation    
        - send file to be analyzed

## endpoints
 - calendar
   - {tool}
     - events
        - GET: return a list of events currently scheduled for the tool
        - POST: add an event to a tools calendar
        - DELETE: delete an calendars event
    