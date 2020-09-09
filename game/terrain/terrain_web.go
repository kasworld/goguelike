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

package terrain

import (
	"fmt"
	"html/template"
	"image/png"
	"net/http"

	"github.com/kasworld/weblib"
)

func (tr *Terrain) Web_TerrainInfo(w http.ResponseWriter, r *http.Request) {
	cmd := weblib.GetStringByName("cmd", "", w, r)
	var err error
	switch cmd {
	case "Init":
		err = tr.Init()
	case "Ageing":
		err = tr.Ageing()
	case "ResetAgeing":
		err = tr.ResetAgeing()
	case "AgeingNoCheck":
		err = tr.AgeingNoCheck()
	}
	if err != nil {
		errstr := fmt.Sprintf("fail to %v", cmd)
		tr.log.Warn(errstr)
		http.Error(w, errstr, 404)
		return
	}

	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>terrain {{.}}</title>
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
		x = Math.round(pos.x);
		y = Math.round(pos.y);
		getTileInfo("/terraintile?floorname={{.Name}}&x="+x+"&y="+y);
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
	{{.}}
	<br/>
	<a href= "/floor?floorname={{$.Name}}" >
		[Goto Floor {{.GetName}}]
	</a>
	<br/>
    <a href='/terrain?floorname={{.Name}}&move=Before'>[Before]</a>
    <a href='/terrain?floorname={{.Name}}&move=Next'>[Next]</a>
    <a href='/terrain?floorname={{.Name}}&cmd=Init'>[Init]</a>
    <a href='/terrain?floorname={{.Name}}&cmd=Ageing'>[Ageing]</a>
    <a href='/terrain?floorname={{.Name}}&cmd=ResetAgeing'>[ResetAgeing]</a>
    <a href='/terrain?floorname={{.Name}}&cmd=AgeingNoCheck'>[AgeingNoCheck]</a>
	<br/>
	<span id="tileInfo"></span>
	<br/>
	<img src="/terrainimageautozoom?floorname={{.Name}}" id="terimg" onmousemove="showTile(this,event)">
	<br/>

	{{range $i,$v := .GetScript}}
	{{$v}}
	<br/>
	{{end}}
	<hr/>

	Total discover area {{.GetTile2Discover }}
	<br/>
	Starting ResourceTileArea :	{{.GetOriRcsTiles.CalcStat}}
	<br/>
	Current ResourceTileArea : {{.GetRcsTiles.CalcStat}}
	<br/>
	TileArea : {{.GetTiles.CalcStat}}
	<br/>
	Field Object List <br/>
	{{range $i, $v := .GetFieldObjPosMan.GetAllList}}
		{{$v}}
		<br/>
	{{end}}
	<br/>
	ActiveObj count {{.GetActiveObjCount}} start CarryObj count {{.GetCarryObjCount}} 
	<br/>
	<hr/> 
	<table border=1 style="border-collapse:collapse;"> 
	{{range $i, $v := .GetRoomList}}
		{{$v}}
		<br/>
	{{end}}
	</table>
	<hr/>
	</body> </html> `)
	if err != nil {
		tr.log.Error("%v %v", tr, err)
	}
	if err := tplIndex.Execute(w, tr); err != nil {
		tr.log.Error("%v", err)
	}
}

func (tr *Terrain) calcZoom() int {
	zoom := 1024 / tr.Xlen
	if zoom > 1024/tr.Ylen {
		zoom = 1024 / tr.Ylen
	}
	if zoom < 2 {
		zoom = 2
	}
	return zoom
}

func (tr *Terrain) Web_TerrainImageZoom(w http.ResponseWriter, r *http.Request) {
	zoom := weblib.GetIntByName("zoom", 1, w, r)
	img := tr.serviceTileArea.ToImage(zoom)
	err := png.Encode(w, img)
	if err != nil {
		tr.log.Error("%v", err)
	}
}

func (tr *Terrain) Web_TerrainImageAutoZoom(w http.ResponseWriter, r *http.Request) {
	zoom := tr.calcZoom()
	img := tr.serviceTileArea.ToImage(zoom)
	err := png.Encode(w, img)
	if err != nil {
		tr.log.Error("%v", err)
	}
}

func (tr *Terrain) Web_TileInfo(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%v\n", r.RequestURI)
	x := weblib.GetIntByName("x", -1, w, r)
	y := weblib.GetIntByName("y", -1, w, r)
	if x == -1 || y == -1 {
		tr.log.Error("invalid pos %v", r.RequestURI)
		return
	}
	zoom := tr.calcZoom()
	x /= zoom
	y /= zoom
	x, y = tr.WrapXY(x, y)
	tl := tr.serviceTileArea.GetByXY(x, y)
	rtl := tr.resourceTileArea.GetXY(x, y)
	fmt.Fprintf(w, "[%v %v] %v %v", x, y, tl, rtl)
}
