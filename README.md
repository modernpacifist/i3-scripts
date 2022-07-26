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

### \_\_mark\_container
#### Description
This script marks current focused container with user prompt from i3-input  
#### Usage
Create a shortcut in your ~/.config/i3/config to call the script directly  
`bindsym --release $mod+Mod1+m exec __mark_container`  
This script will prompt the user twice, based on the desired mark.  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-1.gif)  
Hehe  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-f1.gif)  

### \_\_get\_current\_marks

#### Description
This script returns all current mark identifiers of marked containers  

#### Usage
Append through pipe to `status_command` in your ~/.config/i3/config  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_marks`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-get-current-marks-demonstration.png)  
If a container was marked, its mark will be displayed in the end of your i3status, each mark separated by '|' symbol  

### \_\_get\_current\_track

#### Description
This script returns the artist and a track title of the currently playing composition in i3status. 

#### Usage
Append through pipe to `status_command` in your ~/.config/i3/config  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_track`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-get-current-track-demonstration.png)  

