// this NEEDS
//	· https://github.com/owseiwastaken/gutil
//	· https://github.com/owseiwastaken/termin


package main

include "gutil"
include "termin"

func Compress( x []byte ) ( string ) {
	buff := ""
	for i:=0;i!=6;i++{
		if (x[i] == 0) { break }
		buff+=spf("%.3d,", x[i])
		lk = append(lk, x[i])
	}
	return buff
}


func gtk ( w Window ) ( string ) {
	x:=read(w)
	lk = []byte{}
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


const (
	MI_INSERT = iota
	MI_NORMAL = iota
	MI_REPLACE = iota
)

var (
	prtinfocl = color(200,200,200,17,17,17)
	filename = "out.txt"

	curtypes = map[int]string{
		MI_INSERT : "I-beam",
		MI_NORMAL : "block",
		MI_REPLACE : "underline",
	}
	modnames = map[int]string{
		MI_INSERT : "INSERT",
		MI_NORMAL : "NORMAL",
		MI_REPLACE : "REPLACE",
	}
	modcolor = map[int]string{
		MI_INSERT : color(255,255,255, 110,185,185),
		MI_NORMAL : color(255,255,255, 170,170,170),
		MI_REPLACE : color(255,164,0 , 170,170,170),
	}

	line string // line cont
	file = []string{} // all lines

	x int = 0 // cursor pos in line
	y int = 0 // cursor pos in file
	winoff Ordenate // win (view) offset

	yl = func()(int){return len(spf("%v", y+1))} // len of line number

	k string // key
	running bool = true // end loop
	mode = MI_NORMAL

	lk []byte// last key typed (for logging)
	LOG []string// log

	tbuf1 string // temporary buffer (for case enter;backspace)
	tbuf2 string // temporary buffer (for case enter;backspace)
	at int
)

func redraw () () {
	clear()
	prtinfo()
	for i:=0;i<len(file);i++{
		wprint(Win, i, 0,
			spf(
				"%s%d %s",
				strings.Repeat(" ", 3-(len(spf("%v", i+1)))), i+1, file[i],
			),
		)
	}
	//prtln()
	wmove(Win, y, x+4)
}

func prtinfo()(){
	// "enable" grey bkground
	wuprint(Win, termy-2, 0, prtinfocl)

	wDrawLine(Win, termy-2, " ")
	wprint(Win, termy-2, 0,
		spf(
			"%s %s %s %s",
			modcolor[mode], modnames[mode], prtinfocl, filename,
		),
	)
	ShowCursor()

	// "disable" grey bkground
	wuprint(Win, termy-2, termx, "\x1b[0m")

	wprint(Win, termy-1, 0,
		spf(
			"y:%d, x:%d",
			y, x,
		),
	)
}

func prtln()(){
	wprint(Win, y, 0,
		spf("%s%d %s "/*the last char is here to clean deleted chars*/,
		strings.Repeat(" ", 3-yl()), y+1, line),
	)
	prtinfo()
	wmove(Win, y, x+4)
}
include "control"

func M_insert (k string) () {
	switch (k){
	case "esc":
		mode = MI_NORMAL
		ReCur()
	case "f1":
	case "f2":
	case "f3":
	case "f4":
	case "f5":
	case "f6":
	case "f7":
	case "f8":
	case "f9":
		running = false
	case "backspace":
		if (x!=0) {
			line = line[:x-1] + line[x:]
			x--
		} else {
			if (y!=0) {
				tbuf1 = file[y-1] // prev
				tbuf2 = file[y] // now
				x = len(file[y-1])
				file[y-1] = tbuf1+tbuf2
				file = append(file[:y], file[y+1:]...)
				y--
				line = file[y]
				redraw()
			}
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
		if (y!=(termy-3)){
			if x == len(line) {
				// dumb but it works
				file = append(file[:y], append([]string{file[y]}, file[y:]...)...)
				file[y+1] = ""
			} else {
				// >hel$lo
				// hel, lo
				// >file..., hel, $lo, ...file
				tbuf1 = line[:x] // prev
				tbuf2 = line[x:] // next
				if len(file) > 1 {
					file = append(file[:y], append([]string{tbuf1, tbuf2}, file[y+1:]...)...)
				} else { // len(file) == 1
					file = []string{tbuf1, tbuf2}
				}
				x = 0
			}
			y++
			redraw()
			line = file[y]
			if len(line) < x{
				x = len(line)
			}
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
		} else {
			if (y+1!=len(file)) {
				tbuf1 = file[y] // prev
				tbuf2 = file[y+1] // now
				//x = len(file[y+1])
				file[y] = tbuf1+tbuf2
				file = append(file[:y+1], file[y+2:]...)
				//y--
				line = file[y]
				redraw()
			}
		}
	case "NULL":
		// key (KeyCode) [(KeyHint)] mappend to NULL
		log(spf("key %v [%s] mapped to NULL", lk, string(lk)))
	default:
		line = line[:x] + k + line[x:]
		x++
	}
}

func ExecCmd ( ca []string ) () {
	c := ca[0]
	ca = ca[1:]
	switch (c) {
	case "quit", "q":
		running = false
	case "write", "w":
		if len(ca) != 0{
			switch (len(ca)) {
			case 1:
				WriteFile(ca[0], strings.Join(file, "\n"))
			default:
				WriteFile(ca[0], strings.Join(ca[1:], " "))
			}
		} else {
			WriteFile(filename, strings.Join(file, "\n"))
		}
	case "wq":
		WriteFile(filename, strings.Join(file, "\n"))
		running = false
	}
}

func GetCmd () ([]string) {
	c := "" // command
	k := "" // key
	nx := 0  // pos
	wDrawLine(Win, termy-1, " ")
	wDrawLine(Win, termy-2, " ")
	CursorMode(curtypes[MI_INSERT])
	for {

		wprint(Win, termy-1, 0,
			spf(":%s ", c),
		)
		wmove(Win, termy-1, nx+1)
		Win.stream.Flush()
		k = gtk(Win)

		switch (k) {
		case "backspace":
			if (nx!=0) {
				c = c[:nx-1] + c[nx:]
				nx--
			}
		case "space":
			c+=" "
			nx++
		case "left":
			if (nx!=0) {
				nx--
			}
		case "right":
			if (nx!=len(c)) {
				nx++
			}
		case "enter":
			CursorMode(curtypes[MI_NORMAL])
			return strings.Split(c, " ")
		default:
			if len(k) == 1{
				c = c[:nx] + k + c[nx:]
				nx++
			}
		}
	}
}

func M_normal (k string) () {
	switch (k) {
	case ":":
		ExecCmd(GetCmd())
	case "i":
		mode = MI_INSERT
		ReCur()
		x--
	case "a":
		mode = MI_INSERT
		ReCur()
	}
}

func ReCur () () {
	CursorMode(curtypes[mode])
}

func main(){
	if argc != 0 {
		switch (argc) {
		case 1:
			filename = argv[0]
		}
	}
	// use termin
	TerminInit()

	ReCur()

	file = append(file, line)
	for running{
		prtln()
		k = gtk(Win)
		if mode == MI_INSERT {
			M_insert(k)
		} else if mode == MI_NORMAL {
			M_normal(k)
		}
		file[y] = line
	}
	clear()
	if len(LOG) > 0 {
		WriteFile("log", strings.Join(LOG, "\n"))
	}

	TerminEnd()
	CursorMode("block")
	exit(0)
}


