// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package floor

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/kasworld/goguelike/game/fieldobject"

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/weblib"
)

func (f *Floor) GetAllActiveObj() []gamei.ActiveObjectI {
	aos := make([]gamei.ActiveObjectI, 0)
	for _, v := range f.GetActiveObjPosMan().GetAllList() {
		aos = append(aos, v.(gamei.ActiveObjectI))
	}
	return aos
}

func (f *Floor) Web_FloorInfo(w http.ResponseWriter, r *http.Request) {
	cmd := weblib.GetStringByName("cmd", "", w, r)
	switch cmd {
	case "Ageing":
		f.processAgeing()
	}

	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>floor {{.}}</title>
	<script>
	function getCursorPos(img, e) {
		var a, x = 0, y = 0;
		e = e || window.event;
		/*get the x and y positions of the image:*/
		a = img.getBoundingClientRect();
		/*calculate the cursor's x and y coordinates, relative to the image:*/
		x = e.pageX - a.left;
		y = e.pageY - a.top;
		/*consider any page scrolling:*/
		x = x - window.pageXOffset;
		y = y - window.pageYOffset;
		return {x : x, y : y};
	}
	function showTile(obj,e) {
		// img = document.getElementById("terimg");
		img = obj;
		// console.log(e, img);
		var pos, x, y;
		/*prevent any other actions that may occur when moving over the image*/
		e.preventDefault();
		/*get the cursor's x and y positions:*/
		pos = getCursorPos(img, e);
		x = pos.x;
		y = pos.y;
		getTileInfo("/floortile?floorid={{.GetUUID}}&x="+x+"&y="+y);
	}
	function getTileInfo(url) {
		var xhttp;
		xhttp=new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
			ajaxCallback(this);
			}
		};
		xhttp.open("GET", url, true);
		xhttp.send();
	}
	function ajaxCallback(xhttp) {
		document.getElementById("tileInfo").innerHTML =	xhttp.responseText;
	}
	</script>
	</head>
	<body>
	{{.}} {{.GetEnvBias}}
	<br/>
	<a href= "/terrain?floorid={{$.GetUUID}}" >
		[Goto Terrain {{.GetName}}]
	</a>
	<br/>
    <a href='/floor?floorid={{.GetUUID}}&move=Before'>[Before]</a>
	<a href='/floor?floorid={{.GetUUID}}&move=Next'>[Next]</a>
    <a href='/floor?floorid={{.GetUUID}}&cmd=Ageing'>[ProcessAgeing]</a>
	<br/>
	<span id="tileInfo"></span>
	<br/>
	<img src="/floorimageautozoom?floorid={{.GetUUID}}" id="floorimg" onmousemove="showTile(this,event)" >

	<table border=1 style="border-collapse:collapse;">` +
		gamei.ActiveObjectI_HTML_tableheader + `
	{{range $i,$v := .GetAllActiveObj}}
		{{if $v}}` +
		gamei.ActiveObjectI_HTML_row + `
		{{end}}
	{{end}}` +
		gamei.ActiveObjectI_HTML_tableheader + `
	</table>
	</body> </html> `)
	if err != nil {
		f.log.Error("%v %v", f, err)
	}
	if err := tplIndex.Execute(w, f); err != nil {
		f.log.Error("%v", err)
	}
}

func (f *Floor) MakeImage(zoom int) *image.RGBA {
	ta := f.GetTerrain().GetTiles()
	img := image.NewRGBA(image.Rect(0, 0, f.w*zoom, f.h*zoom))

	for srcY := 0; srcY < f.h; srcY++ {
		for srcX := 0; srcX < f.w; srcX++ {
			co := color.RGBA{0x00, 0x00, 0x00, 0xff}
			if ao := f.aoPosMan.Get1stObjAt(srcX, srcY); ao != nil {
				r, g, b := ao.(gamei.ActiveObjectI).GetBias().ToRGB()
				co = color.RGBA{r, g, b, 0xff}
			} else if po := f.poPosMan.Get1stObjAt(srcX, srcY); po != nil {
				switch po.(type) {
				default:
					f.log.Error("unknown carryingobj %v", po)
				case gamei.EquipObjI:
					r, g, b := po.(gamei.EquipObjI).GetBias().ToRGB()
					co = color.RGBA{r, g, b, 0xff}
				case gamei.PotionI:
					rgb := po.(gamei.PotionI).GetPotionType().Color24().UInt8Array()
					co = color.RGBA{rgb[0], rgb[1], rgb[2], 0xff}
				case gamei.MoneyI:
					co = color.RGBA{0xff, 0xd7, 0x00, 0xff} // gold color
				}
			} else if iao := f.foPosMan.Get1stObjAt(srcX, srcY); iao != nil {
				ww, ok := iao.(*fieldobject.FieldObject)
				if !ok {
					f.log.Fatal("not *fieldobject.FieldObject %v", iao)
					continue
				}
				oco := ww.GetActType().Color24()
				co = color.RGBA{oco.R(), oco.G(), oco.B(), 0xff}
			} else {
				r, g, b := ta[srcX][srcY].ToRGB()
				co = color.RGBA{r, g, b, 0xff}
			}
			for y := srcY * zoom; y < srcY*zoom+zoom; y++ {
				for x := srcX * zoom; x < srcX*zoom+zoom; x++ {
					img.Set(x, y, co)
				}
			}
		}
	}
	return img
}

func (f *Floor) calcZoom() int {
	zoom := 1024 / f.w
	if zoom > 1024/f.h {
		zoom = 1024 / f.h
	}
	if zoom < 2 {
		zoom = 2
	}
	return zoom
}

func (f *Floor) Web_FloorImageAutoZoom(w http.ResponseWriter, r *http.Request) {
	zoom := f.calcZoom()
	img := f.MakeImage(zoom)
	err := png.Encode(w, img)
	if err != nil {
		f.log.Error("%v", err)
	}
}

func (f *Floor) Web_FloorImageZoom(w http.ResponseWriter, r *http.Request) {
	zoom := weblib.GetIntByName("zoom", 1, w, r)
	img := f.MakeImage(zoom)
	err := png.Encode(w, img)
	if err != nil {
		f.log.Error("%v", err)
	}
}

func (f *Floor) Web_TileInfo(w http.ResponseWriter, r *http.Request) {
	x := weblib.GetIntByName("x", -1, w, r)
	y := weblib.GetIntByName("y", -1, w, r)
	if x == -1 || y == -1 {
		f.log.Error("invalid pos %v", r.RequestURI)
		return
	}
	zoom := f.calcZoom()
	x /= zoom
	y /= zoom
	tl := f.terrain.GetTiles()[x][y]
	rtl := f.terrain.GetRcsTiles().GetXY(x, y)
	fmt.Fprintf(w, "[%v %v] %v %v", x, y, tl, rtl)
	// aolist := f.aoPosMan.GetObjListAt(x, y)
	// polist := f.aoPosMan.GetObjListAt(x, y)
	// for _, v := range aolist {
	// 	fmt.Fprintf(w, "%v</br>", v)
	// }
	// for _, v := range polist {
	// 	fmt.Fprintf(w, "%v</br>", v)
	// }
}
