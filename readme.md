# tool- cal
an event calendar API to manage the use of tools at HackRva

[![Build Status](https://travis-ci.org/Athulus/tool-cal.svg?branch=master)](https://travis-ci.org/Athulus/tool-cal)

## building and running
1.  `git clone` tool-cal repo into your gopath/src
2. `cd tool-cal`
3. `make` will run install dependencies, run test, and compile a static binary
4. `make deploy` will build a docker image and run the image and a redis server with docker-compose

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
    