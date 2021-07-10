
function regionChanged() {
    console.log("regionChanged!");
    let reg = document.getElementById('Region');
    if (typeof(reg) == 'undefined') return;

    let xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			console.log('{'+this.responseText+'}');
            let obj = JSON.parse(this.responseText);
            Object.keys(obj).forEach((key)=>{try{document.getElementById(key).value=obj[key]}catch{}});
            let rdpage = document.getElementById('regiondetail');
            rdpage.style.display = 'inherit';
		}
	};
	xhttp.open("GET", encodeURI("/ajax?c=region&key="+reg.value), true);
	xhttp.send();

}

function dateFromIso(date) {
    return new Date(Number(date.substring(0, 4)), Number(date.substring(5, 7))-1, 
                Number(date.substring(8, 10)), Number(date.substring(11, 13)), 
                Number(date.substring(14, 16)), Number('00'));

}

function rallyTimesChanged() {
    let startdt = document.getElementById('Startdate').value + ' ' + document.getElementById('Starttime').value;
    let finishdt = document.getElementById('Finishdate').value + ' ' + document.getElementById('Finishtime').value;
    if (finishdt < startdt) {
        alert(finishdt+' < '+startdt);
        return;
    }
    let sdt = dateFromIso(startdt);
    let fdt = dateFromIso(finishdt);
    let hrs = Math.floor((fdt - sdt) / (1000*60*60));
    let max = document.getElementById('Maxhours');
    max.value='0';
    max.setAttribute('max',hrs);
    max.value = hrs;
}

function showPrompt(obj) {

    let pnl = document.getElementById('promptpanel');
    if (typeof(pnl) == 'undefined') return;
    let txt = obj.getAttribute('title')
    if (txt) {
        pnl.innerText = txt + " ;;; ";
        return;
    }
    txt = obj.parentElement.getAttribute('title');
    pnl.innerText = txt;

}