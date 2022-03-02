package main

include "gutil"
include "termin"

func Compress( x []byte ) ( string ) {
	var buff = ""
	for i:=0;i!=6;i++{
		if (x[i] == 0) { break }
		buff+=spf("%.3d,", x[i])
	}
	return buff
}


func gtk ( w Window ) ( string ) {
	x:=read(w)
	e, ok := Control[Compress(x)]
	if (!ok) {
		e = "NULL"
	}
	return e
}

var (
	line string // line cont
	x int = 0// cursor pos
	y int = 0
	yl = func()(int){return len(spf("%v", y))}
	k string // key
	running bool = true
	tw string
)

include "control"
func main(){
	// use termin
	TerminInit()

	CursorMode("I-beam")
	var prtln = func()(){
		tw = spf("%s%d %s", strings.Repeat(" ", 3-yl()), y+1, line)
		wprint(Win, y, 0, tw+" "+spos(y,x+4))
	}

	for running{
		prtln()
		k = gtk(Win)
		switch (k){
			case "enter":
				running = false
			case "backspace":
				if (x!=0) {
					line = line[:x-1] + line[x:]
					x--
				}
			case "left":
				if (x!=0) {
					x--
				}
			case "right":
				if (x!=len(line)) {
					x++
				}
			case "space":
				line+=" "
				x++
			case "NULL":
			default:
				line = line[:x] + k + line[x:]
				x++
		}
	}
	clear()
	PS(line)

	TerminEnd()
	CursorMode("block")
	exit(0)
}


