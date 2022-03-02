// this NEEDS
//	· https://github.com/owseiwastaken/gutil
//	· https://github.com/owseiwastaken/termin


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

func log ( add string ) () {
	adda := strings.Split(add, "\n")
	if len(adda) > 1{
		LOG = append(LOG, "")
	}
	for i:=0;i<len(adda);i++ {
		LOG = append(LOG, adda[i])
	}
}

var (
	line string // line cont
	file = []string{} // all lines

	x int = 0// cursor pos in line
	y int = 0// cursor pos in file

	yl = func()(int){return len(spf("%v", y+1))} // len of line number

	k string // key
	running bool = true // end loop

	LOG []string// log

	tbuff string // temporary buffer (for case enter)
	at int
)

func redraw () () {
	clear()
	for i:=0;i<len(file);i++{
		wprint(Win, i, 0,
			spf(
				"%s%d %s",
				strings.Repeat(" ", 3-(len(spf("%v", i+1)))), i+1, file[i],
			),
		)
	}
	wmove(Win, y, x+4)
}

func Exec (c string) (string) {
	switch (c) {
		case "quit", "q":
			running = false
		default:
			return spf("%sNot an editor command: %s%s", RGB(255, 0, 0), c, RGB(255,255,255))
	}
	return ""
}

include "control"
func main(){
	// use termin
	TerminInit()

	CursorMode("I-beam")
	var prtln = func()(){
		wprint(Win, y, 0,
			spf("%s%d %s "/*the last char is here to clean deleted chars*/,
			strings.Repeat(" ", 3-yl()), y+1, line),
		)
		wmove(Win, y, x+4)
	}

	file = append(file, line)
	for running{
		prtln()
		k = gtk(Win)
		switch (k){
			case "f9":
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
			case "up":
				if (y!=0){
					y--
					line = file[y]
				}
				if len(line) < x{
					x = len(line)
				}
			case "enter":
				// TODO borked
				// TODO redesign idea
				//	file IS a string (...\n...\n...\n)
				//	x = file[@]
				//	y = count(file[:x], "\n")
				//	cx = x-findlast@(file[:x], "\n") // cursor x
				if (y!=termy){
					y++
					if len(line)==0 || y == len(file){
						file = append(file, "")
					} else {
						file[y] = ""
					}
					redraw()
				}
				line = file[y]
				if len(line) < x{
					x = len(line)
				}
			case "down":
				if (y!=termy){
					y++
					if (len(file)<=y) {
						y--
					}
					line = file[y]
				}
				if len(line) < x{
					x = len(line)
				}
			case "space":
				line+=" "
				x++
			case "delete":
				if (x!=len(line)) {
					line = line[:x] + line[1+x:]
				}
			case "NULL":
			default:
				line = line[:x] + k + line[x:]
				x++
		}
		file[y] = line
	}
	clear()
	PS(line)
	WriteFile("log", strings.Join(LOG, "\n"))

	TerminEnd()
	CursorMode("block")
	exit(0)
}


