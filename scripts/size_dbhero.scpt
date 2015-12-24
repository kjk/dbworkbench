-- http://stackoverflow.com/questions/6164164/resizing-windows-of-unscriptable-applications-in-applescript
-- a script to set a known size of dbHero window on mac so that we take screenshots of predictable size
-- do 'open scripts/size_dbhero.scpt' to open in Script Editor, press "play" button, make sure to allow
-- Script Editor to use accessibility apis in system preferences (Security & Privacy / Privacy)
tell application "System Events"
	tell process "dbHero"
		set size of window 1 to {1024, 640}
	end tell
end tell
