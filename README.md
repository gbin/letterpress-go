letterpress-go
==============

This is a simple min-max algorithm for the game of letterpress (or wordstrike under android).

It is an experiment I ran to see how well a computer can play at the game.

I have only one requirement for the people who wants to use it :
Don't be an a** and always tell your opponent that you are using this as it seriously spoils the fun.

How to compile it ?
-------------------

You need a golang installation (see http://www.golang.org) with $GOPATH correctly set.

```bash
# clone it to your local go src folder
git clone https://github.com/gbin/letterpress-go.git $GOPATH/src
cd $GOPATH/src/letterpress-go
go build
# an executable called letterpress-go should be there.
```

How to use it ?
---------------

```bash
./letterpress-go [input-file]
```

The input format is a text file with a very rigid format :
1. the first 5 lines are the letters of the game
2. one blank line
3. the next five lines are the color mask : [space] for white, r for red, b for blue
4. one blank line
5. the list of words already played

For example (you can copy paste it as a base):
```itcla
nkfln
edkyu
geeez
hotss

rrrrr
rrrrb
rbb b
b    
b   b

unsticked
gesellschaften
knuckeheads
calflike
```
