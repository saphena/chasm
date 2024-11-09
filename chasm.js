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
    if (nm != "") {
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

function reloadClaimslog() {
  let frm = document.getElementById("claimslogfrm");

  let url = "/claims?x=x";
  let inps = frm.getElementsByTagName("select");
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm != "") {
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
    if (nm != "") {
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

const re = '\s*[A-Za-z]*(\d+)[A-Za-z]*\s*\,?\s*([a-zA-Z0-9\-]+)\s*\,?\s*(\d+)?\.*\d*\s*\,?\s*(\d\d?[.:]*\d\d)?\s*(.*)'
}