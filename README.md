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
This script will prompt the user once or twice, based on the user input.  
Containers can be marked with numbers [0-9] or with the 'f' prefix, therefore having f[0-9] mark pattern.
This is done to access the marked container by making shortcuts to the numbers or function keys directly.  
User prompt will be unitary if the user inputs a single number and the container will be marked with it.  
This unitary prompt behavior is shown in the gif below, first the script is called via i3 shortcut, then key '1' is pressed as a sample.  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-1.gif)  
Therefore this container is marked with '1' mark identifier, which can be accessed via i3's `bindsym --release $mod+1 [con_mark="^1$"] focus` for example.  
If the user decides to mark the container with function key identifier, then the first input must be symbol 'f', after which the prompt will appear again, asking for the number [0-9].  
This function key mode is shown in the gif below, first the script is called via i3 shortcut, then key 'f' is pressed, which calls another prompt for the number of the function key [0-9].  
![alt text](https://github.com/modernpacifist/i3-scripts/blob/master/img/i3-mark-container-f1.gif)  
Therefore this container is marked with 'f1' mark identifier.  
I suggest you create a mode or a series of shortcuts to access each marked container directly in your ~/.config/i3/config  

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

