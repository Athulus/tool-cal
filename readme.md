
# tool scheduling app 
- calendar    
    - calendar for each tool    
    - user can add time to a tool
- users    
    - users have a certian amount of time they can use per month
- tools
    - tool use time estimation    
        - send file to be analyzed

# endpoints
 - calendar
   - {tool}
     - events
        - GET: return a list of events currently scheduled for the tool
        - POST: add an event to a tools calendar
   - POST: add a new tool calendar
   - GET: get all current tool calendars
   - DELTE: delete a tool calendar
    