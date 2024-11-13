"use strict";

const CAT_OrdinaryScoringRule = 0;
const CAT_DNF_Unless_Triggered = 1;
const CAT_DNF_If_Triggered = 2;
const CAT_PlaceholderRule = 3;
const CAT_OrdinaryScoringSequence = 4;

function chgRuleType(obj) {
  let div = obj.parentElement.parentElement;

  console.log('RuleType is "' + obj.value + '"');
  let fs = div.querySelectorAll("fieldset");
  for (let i = 0; i < fs.length; i++) {
    if (fs[i].classList.contains("rule" + obj.value)) {
      fs[i].classList.remove("hide");
    } else {
      fs[i].classList.add("hide");
    }
  }
}

function chgAxis(obj) {
  let div = obj.parentElement.parentElement;

  // This only works for singular axes
  let fs = div.querySelectorAll("label");
  for (let i = 0; i < fs.length; i++) {
    if (fs[i].getAttribute("for") == "Cat") {
      fs[i].innerHTML = obj.options[obj.selectedIndex].innerHTML;
      break;
    }
  }

  console.log(div);
  console.log(obj);
  let url = "x?f=axiscats&a=" + obj.value + "&s=0";
  console.log(url);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      if (!data.OK) {
        console.error(`Error! ${data.Msg}`);
      } else {
        let cat = div.querySelector("#Cat");
        console.log(data);
        cat.innerHTML = data.Msg;
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });
}

function fetchBonusDetails(obj) {
  const allflags = "ABDFNRT";

  let b = obj.value;
  let url = "x?f=fetchb&b=" + b.toUpperCase();
  console.log(url);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      if (!data.ok) {
        console.error(`Error! ${data.ok}`);
      } else {
        let bd = document.getElementById("bonusDetails");
        bd.innerHTML = data.name;
        let flags = data.flags;
        for (let i = 0; i < allflags.length; i++) {
          let f = allflags[i];
          let s = document.getElementById("bflag" + f);
          if (flags.includes(f)) {
            s.classList.remove("hide");
          } else {
            s.classList.add("hide");
          }
        }
        let qa = document.getElementById("CorrectAnswer");
        if (qa) qa.innerHTML = data.answer;
        //console.log(data);
        /*
        if (data.team) {
          edf.classList.remove('hide')
        } else {
          edf.classList.add('hide')
        }
          */
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });
}

function fetchEntrantDetails(obj) {
  let e = obj.value;
  let url = "x?f=fetche&e=" + e;
  console.log(url);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      if (!data.ok) {
        console.error(`Error! ${data.ok}`);
      } else {
        let ed = document.getElementById("entrantDetails");
        let edf = document.getElementById("edflag");
        //console.log(data);
        ed.innerHTML = data.name;
        if (data.team) {
          edf.classList.remove("hide");
        } else {
          edf.classList.add("hide");
        }
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });
}

function showEvidence(obj) {
  killReload();
  let ft = document.getElementById("finetune");
  let ov = document.getElementById("claimstats");
  ft.classList.remove("hide");
  ov.classList.add("hide");
}

function showRule(obj) {
  window.location.href = "/rule?r=" + obj.getAttribute("data-rowid");
}

function setupForm() {
  //chgRuleType(document.getElementById("RuleType"));
  //chgAxis(document.getElementById("Axis"));
}

function showEBC(obj) {
  window.location.href = "/ebc?c=" + obj.getAttribute("data-claimid");
}

function showFirstClaim() {
  let rows = document.querySelectorAll("fieldset.row");
  showEBC(rows[1]); // 1 not 0; 0 = hdr
}
function cycleImgSize(obj) {
  let img = obj.id;
  let sz = obj.style.width;
  console.log("sz == " + sz);
  let otherimg = "";
  if (img == "ebcimgdiv") {
    otherimg = "bonusimgdiv";
  } else {
    otherimg = "ebcimgdiv";
  }
  let other = document.getElementById(otherimg);
  if (sz == "50%" || sz == "") {
    other.style.width = "100px";
    obj.style.width = "99%";
  } else {
    other.style.width = "50%";
    obj.style.width = "50%";
  }
}

function goHome(obj) {
  window.location.href = "/";
}
function swapimg(img) {
  let me = img.getAttribute("src");
  let main = document.getElementById("imgdivimg");
  let mainsrc = main.getAttribute("src");
  let inp = document.getElementById("chosenPhoto");
  main.setAttribute("src", me);
  img.setAttribute("src", mainsrc);
  inp.setAttribute("value", me);
}

function closeEBC(obj) {
  let frm = document.getElementById("ebcform");
  let dec = document.getElementById("chosenDecision");
  dec.value = obj.getAttribute("data-result");
  let url = "/x?f=saveebc";
  let inps = frm.getElementsByTagName("input");
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm && nm != "") {
      url += "&" + nm + "=" + encodeURIComponent(inps[i].value);
    }
  }

  console.log(url);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      if (!data.OK) {
        window.location.href = "/listebc";
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });

  //frm.submit();
}


function fixClaimTimeISO() {

  let iso = document.getElementById('ClaimTimeISO')
  let dt = document.getElementById('ClaimDate')
  let tm = document.getElementById('ClaimTime')

  iso.value = dt.value+'T'+tm.value+iso.value.substring(16)
}
function reloadClaimslog() {
  let frm = document.getElementById("claimslogfrm");

  let url = "/claims?x=x";
  let inps = frm.getElementsByTagName("select");
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm && nm != "") {
      url +=
        "&" +
        nm +
        "=" +
        encodeURIComponent(
          inps[i].options[inps[i].selectedIndex].getAttribute("value")
        );
    }
  }

  console.log(url);
  window.location.href = url;
}

function reloadRankings(fld, val) {
  const args = new Map();

  let frm = document.getElementById("rankingsfrm");
  let inps = frm.getElementsByTagName("input");
  let url = "/qlist?x=x";
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm && nm != "") {
      args.set(nm, inps[i].getAttribute("value"));
    }
  }
  if (args.get("seq") == val && fld == "seq") {
    if (args.get("desc") == "") {
      args.set("desc", "desc");
    } else {
      args.set("desc", "");
    }
  } else {
    args.set("desc", "");
  }
  args.set(fld, val);
  args.forEach(function (val, key) {
    url += "&" + key + "=" + encodeURIComponent(val);
  });

  console.log(url);
  window.location.href = url;
}

function showAboutChasm(obj) {
  window.location.href = "/about";
}
function toggleLicenceMIT(obj) {
  let mit = document.getElementById("LicenceMIT");
  if (!mit) return;
  if (mit.classList.contains("hide")) {
    mit.classList.remove("hide");
  } else {
    mit.classList.add("hide");
  }
}

function showHelp(topic) {
  window.open(
    "/help?topic=" + topic,
    "smhelp",
    "location=no,height=800,width=800,scrollbars=yes,status=no"
  );
}

// This wil parse a correctly formatted email Subject line into
// the relevant fields of a blank new claim form
function pasteNewClaim(obj) {
  const re =
    /\s*[A-Za-z]*(\d+)[A-Za-z]*\s*\,?\s*([a-zA-Z0-9\-]+)\s*\,?\s*(\d+)?\.*\d*\s*\,?\s*(\d\d?[.:]*\d\d)?\s*(.*)/;

  let subject = obj.value;

  let x = re.exec(subject);
  console.log(x);

  if (x.length < 5) {
    return;
  }
  let e = document.getElementById("EntrantID");
  e.value = x[1];

  let b = document.getElementById("BonusID");
  b.value = x[2];

  let o = document.getElementById("OdoReading");
  o.value = x[3];

  let t = document.getElementById("ClaimTime");
  if (x[4].length < 5) {
    x[4] = x[4].slice(0, 2) + ":" + x[4].slice(2);
  }
  t.value = x[4];

  if (x.length > 5) {
    let a = document.getElementById("AnswerSupplied");
    a.value = x[5];
  }
}


function saveUpdatedClaim(obj) {
  let frm = document.getElementById("iclaim");
  let url = "/x?f=saveclaim";
  let inps = frm.getElementsByTagName("input");
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm && nm != "") {
      url += "&" + nm + "=" + encodeURIComponent(inps[i].value);
    }
  }

  console.log(url);
  //alert(url)
  //return
  fetch(url,{method: 'POST'})
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      if (!data.OK) {
        window.location.href = "/claims";
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });

  //frm.submit();
}

function updateClaimDecision(obj) {

  let dec = document.getElementById('chosenDecision')
 
  dec.value = obj.options[obj.selectedIndex].value

}