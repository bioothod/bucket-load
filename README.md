# bucket-load
Backrunner (http://github.com/bioothod/backrunner) access log parser and load generator

It reads provided log, searches for access_log entries and repeats them in realtime to another specified backrunner proxy.
POST entries are transformed into GET requests.
