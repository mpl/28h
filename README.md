28h
===

Computes and prints a week with 28 hours days, like in http://xkcd.com/320/

Example:
--------

You want to see what the week would look like if you got up at 17:00 on thursday, because you have kung-fu at 18:00 this very day.

`
% 28h -day thursday -wake 17:00
`

mon	bed 05:00 up 00:00 bed  
tue	bed 09:00 up  
wed	up 04:00 bed 13:00 up  
thu	up 08:00 bed 17:00 up  
fri	up 12:00 bed 21:00 up  
sat	up 16:00 bed  
sun	bed 01:00 up 20:00 bed  

Install:
--------

 1. Install Go: http://golang.org/doc/install
 2. go get github.com/mpl/28h

