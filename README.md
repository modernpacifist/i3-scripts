# i3-scripts

## Dependencies
```jq```  
## Installation  
Run as root:  
```# ./install.sh```  
## Scripts breakdown
### \_\_get\_current\_marks
This script shows current created marks in i3status, append through pipe to status\_command in your ~/.config/i3/config  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_marks`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-marks-demonstration.png)  
