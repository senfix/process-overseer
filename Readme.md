# Process overseer

## Build

`dep ensure`
`make build`

## Setup 
### Mac OSX

Create new service `custom.process.overseer.plist`

`/Library/LaunchDaemons/` for one instance per pc
`~/Library/LauchAgents/` for user only

````
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>custom.process.overseer</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>ProgramArguments</key>
    <array>
        <string>./process_overseer</string>
    </array>
    <key>WorkingDirectory</key>
    <string>/path/to/compiled/process_overseer</string>
    <key>StandardErrorPath</key>
    <string>/path/to/log/stderr.log</string>
    <key>StandardOutPath</key>
    <string>/path/to/log/stdout.log</string>
  </dict>
</plist>
````

to load configuration into lauchctl
````
sudo launchctl load custom.process.overseer.plist
````

### Linux
