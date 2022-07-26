# i3-scripts

## Dependencies
```
jq  
playerctl  
```  
## Installation  
Run as root:  
```# ./install.sh```  
## Scripts breakdown
### \_\_get\_current\_marks
This script shows current created marks in i3status, append through pipe to `status_command` in your ~/.config/i3/config  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_marks`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-marks-demonstration.png)  
If a container was marked, its mark will appear in the end of your i3status, each mark separated by '|' symbol  

### \_\_get\_current\_track
This script shows the artist and a track title of currently playing composition  
