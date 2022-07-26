# i3-scripts

## Dependencies
```jq```  
## Installation  
Run as root:  
```# ./install.sh```  
## Scripts breakdown
### \_\_get\_current\_marks
This script shows current created marks in i3status, append through pipe to status\_command in your ~/.config/i3/config  
`status\_command i3status --config $HOME/.config/i3/i3status.conf | \_\_get\_current\_marks`  
