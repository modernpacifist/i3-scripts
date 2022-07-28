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
Create a shortcut in your __~/.config/i3/config__ to call the script directly  
`bindsym --release $mod+Mod1+m exec __mark_container`  
User prompt will be unitary if the user inputs a single number and the container will be marked with it.  
This unitary prompt behavior is shown in the gif below, first the script is executed via i3 shortcut, then key '1' is pressed as a sample.  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-1.gif)  
As a result the focused container is marked with `1` mark.   
If the user decides to mark the container with function key identifier, then the first input must be symbol 'f', after which the prompt will appear again, asking for the number [0-9]. 
This function key mode is shown in the gif below, first the script is executed via i3 shortcut, then key 'f' is pressed, which calls another prompt for the number of the function key [0-9].  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-f1.gif)  
As a result the focused container is marked with `f1` mark.  
I suggest you create a mode or a series of shortcuts to access each marked container directly in your __~/.config/i3/config__, e.g:  
`bindsym --release $mod+1 [con_mark="^1$"] focus`  
Or, to focus the container with a function key:  
`bindsym --release $mod+F1 [con_mark="^f1$"] focus`  
Overall usage can be simplified into the algorithm:  
1. execute the script via i3 shortcut
2. input number from 0 to 9 or 'f' symbol
3. if 'f' symbol was pressed, input number from 1 to 9

### \_\_get\_current\_marks

#### Description
This script returns all current mark identifiers of marked containers  

#### Usage
Append through pipe to `status_command` in your __~/.config/i3/config__  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_marks`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-get-current-marks-demonstration.png)  
If a container was marked, its mark will be displayed in the end of your i3status, each mark separated by '|' symbol  

### \_\_get\_current\_track

#### Description
This script returns the artist and a track title of the currently playing composition in i3status. 

#### Usage
Append through pipe to `status_command` in your __~/.config/i3/config__  
`status_command i3status --config $HOME/.config/i3/i3status.conf | __get_current_track`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-get-current-track-demonstration.png)  

### \_\_rename\_i3wm\_workspace

#### Description
This script appends title to the current focused workspace keeping its current index with user prompt. 

#### Usage
Append through pipe to `status_command` in your __~/.config/i3/config__  
`bindsym --release $mod+Shift+r exec __rename_i3wm_workspace`  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-get-current-track-demonstration.png)  

