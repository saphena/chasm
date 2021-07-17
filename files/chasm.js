
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
	xhttp.open("GET", encodeURI("/ajax?c=getregion&key="+reg.value), true);
	xhttp.send();

}

function dateFromIso(date) {
    return new Date(Number(date.substring(0, 4)), Number(date.substring(5, 7))-1, 
                Number(date.substring(8, 10)), Number(date.substring(11, 13)), 
                Number(date.substring(14, 16)), Number('00'));

}

function enableSaveButton(fld) {
    let tr = fld.parentElement.parentElement; // tr
    let flds = tr.cells;
    for (i = 0; i < flds.length; i++) {
        //console.log(flds[i].firstChild);
        for (j = 0; j < flds[i].children.length; j++) 
            if (flds[i].children[j].name=='SaveButton') {
                flds[i].children[j].disabled = false;
                return;
            }
    }
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

function addnewReason() {
    let tab = document.getElementById('reasonstable').getElementsByTagName('tbody')[0];
    let flds = tab.getElementsByTagName('input');
    let lastcode = 0;
    let code = 0;
    for (i = 0; i < flds.length; i++) {
        if (flds[i].name == 'Code') {
            code = parseInt(flds[i].value);
            if (code > lastcode)
                lastcode = code;
        }
    }
    code++;
    let nr = tab.insertRow();
    nr.innerHTML = document.getElementById('newTR').innerHTML;
    flds = nr.getElementsByTagName('input');
    for (i = 0; i < flds.length; i++) {
        if (flds[i].name == 'Code') {
            flds[0].value = code;
        }
    }
}



function deleteReason(delButton) {

    let tr = delButton.parentElement.parentElement; // tr
    let flds = tr.cells;
    let code = '0';
    for (i = 0; i < flds.length; i++) {
        console.log(flds[i].firstChild);
        for (j = 0; j < flds[i].children.length; j++) 
            if (flds[i].children[j].name=='Code') {
                code = flds[i].children[j].value;
            }
    }

    if (code == '') return;

    let xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
            console.log('{'+this.responseText+'}');
            delButton.disabled = true;
            let tr = delButton.parentElement.parentElement;
            tr.parentElement.removeChild(tr);
		}
	};
    let x = "/ajax?c=delreason&Code="+code;
    console.log(x);
	xhttp.open("POST", encodeURI(x), true);
	xhttp.send();

}
function saveReason(obj) {
    let sb, db, cf, code, brief, action, param


    let tr = obj.parentElement.parentElement; // tr
    
    let flds = tr.cells;
    for (i = 0; i < flds.length; i++) {
        console.log(flds[i].firstChild);
        for (j = 0; j < flds[i].children.length; j++) 

        
        switch(flds[i].children[j].name) {
            case 'DeleteButton':
                db = flds[i].children[j];
                break;
            case 'SaveButton':
                sb = flds[i].children[j];
                break;
            case 'Code':
                code = flds[i].children[j].value;
                cf = flds[i].children[j];
                break;
            case 'Briefdesc':
                brief = flds[i].children[j].value;
                break;
            case 'Action':
                action = flds[i].children[j].value;
                break;
            case 'Param':
                param = flds[i].children[j].value;
                break;

    }


    }
    let xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
            console.log('{'+this.responseText+'}');
            sb.disabled = true;
            db.disabled = false;
            cf.readOnly = true;
		}
	};
    let x = "/ajax?c=putreason&Code="+code+'&Briefdesc='+brief+'&Action='+action+'&Param='+param;
    console.log(x);
	xhttp.open("POST", encodeURI(x), true);
	xhttp.send();


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