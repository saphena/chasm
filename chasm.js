// @ts-checks

"use strict";

const myStackItem = "odoStack";

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

function enableDelete(ok) {
  console.log("enableDelete ", ok);
  let del = document.getElementById("enableDelete");
  if (!del) return;
  del.checked = !del.checked;
  enableSave(ok);
}
function enableSave(ok) {
  console.log("enableSave", ok);
  let sav = document.getElementById("updatedb");
  if (!sav) return;
  sav.disabled = !ok;
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
        let pts = document.getElementById("Points");
        pts.value = data.points;
        pts.setAttribute("data-pm", data.pointsaremults);
        let qan = document.getElementById("QuestionAnsweredN");
        if (qan) qan.checked = true;
        let qa = document.getElementById("CorrectAnswer");
        if (qa) qa.innerHTML = data.answer;

        let bp = document.getElementById("bonusPhoto");
        if (bp) {
          bp.setAttribute("src", bp.getAttribute("data-folder") + data.img);
          bp.setAttribute("alt", bp.getAttribute("data-folder") + data.img);
        }

        let pct = document.getElementById("PercentPenalty");
        if (pct && pct.checked) {
          console.log("applying 10%");
          applyPercentPenalty(true);
        }
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
  if (typeof window.killReload === "function") killReload();
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

/** cycleClaimImgSize handles photo displays during claim editing */
function cycleClaimImgSize(obj) {
  let img = obj.id;
  let sz = obj.style.width;
  console.log("sz == " + sz);
  let otherimg = "";
  if (img == "claimPhoto") {
    otherimg = "bonusPhoto";
  } else {
    otherimg = "claimPhoto";
  }
  let other = document.getElementById(otherimg);
  if (sz == "512px" || sz == "") {
    other.style.width = "100px";
    obj.style.width = "99%";
  } else {
    other.style.width = "";
    obj.style.width = "512px";
  }
}

/** cycleImgSize handles photo displays during EBC judging */
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

function oi(obj) {
  obj.classList.remove("oc");
  obj.classList.add("oi");

  console.log("oi called == " + obj.getAttribute("data-save"));
  // autosave handler
  if (obj.timer) {
    clearTimeout(obj.timer);
  }
  //obj.timer = setTimeout(obj.getAttribute('data-save'), 3000, obj);

  if (obj.getAttribute("data-save") == "saveOdo")
    obj.timer = setTimeout(saveOdo, 3000, obj);
  else obj.timer = setTimeout(saveSetupConfig, 3000, obj);
  console.log("oi complete " + JSON.stringify(obj));
}

function saveRS(obj) {
  let e = obj.getAttribute("data-e");
  let url = "/x?f=savers&e=" + e;
  url += "&rs=" + obj.value;
  console.log("saveRS: " + url);
  stackTransaction(url, obj.id);
  sendTransactions();
}
function saveSetupConfig(obj) {
  console.log("saveSetupConfig called");
  if (obj.timer) {
    clearTimeout(obj.timer);
  }
  if (obj.getAttribute("data-chg") == 0) return;
  obj.setAttribute("data-chg", 0);
  let url = "/x?f=putcfg";
  url += "&ff=" + obj.name + "&v=" + encodeURIComponent(obj.value);
  stackTransaction(url, obj.id);

  sendTransactions();
}

function saveSetupFinish(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let url = "/x?f=putcfg";
  let dt = document.getElementById("RallyFinishDate");
  let tm = document.getElementById("RallyFinishTime");
  url += "&ff=RallyFinish&v=" + encodeURIComponent(dt.value + "T" + tm.value);
  stackTransaction(url, obj.id);
  sendTransactions();
}
function saveSetupStart(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let url = "/x?f=putcfg";
  let dt = document.getElementById("RallyStartDate");
  let tm = document.getElementById("RallyStartTime");
  url += "&ff=RallyStart&v=" + encodeURIComponent(dt.value + "T" + tm.value);
  stackTransaction(url, obj.id);
  sendTransactions();
}
function saveOdo(obj) {
  console.log("saveOdo called");
  if (obj.timer) {
    clearTimeout(obj.timer);
  }

  let timeDisplay = document.querySelector("#timenow");

  let ent = obj.getAttribute("data-e");
  let url =
    "/x?f=putodo&e=" +
    ent +
    "&st=" +
    obj.getAttribute("data-st") +
    "&ff=" +
    obj.name +
    "&v=" +
    obj.value +
    "&t=" +
    timeDisplay.getAttribute("data-time");

  stackTransaction(url, obj.id);
}

function sendTransactions() {
  let stackx = sessionStorage.getItem(myStackItem);
  if (stackx == null) return;

  let stack = JSON.parse(stackx);

  console.log(stack);

  let errlog = document.getElementById("errlog");

  while (stack.length > 0) {
    let itm = stack[0];
    stack.splice(0, 1);
    sessionStorage.setItem(myStackItem, JSON.stringify(stack));
    console.log("Sending: " + itm.url);

    fetch(itm.url)
      .then((response) => {
        if (!response.ok) {
          // Handle HTTP errors
          stackTransaction(itm.url, itm.objid);
          //if (errlog){errlog.innerHTML=`HTTP error! Status: ${response.status}`}

          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        if (data.err) {
          // Handle JSON error field
          console.error(`Error: ${data.msg}`);
        } else {
          // Process the data if no error
          //if (errlog){errlog.innerHTML="Hello sailor: "+JSON.stringify(data)}
          console.log("Data:", data);
          document.getElementById(itm.objid).classList.replace("oi", "ok");
          try {
            reloadok = true;
          } catch {}
        }
      })
      .catch((error) => {
        // Handle network or other errors
        //if (errlog) {errlog.innerHTML="ERROR CAUGHT"}
        stackTransaction(itm.url, itm.objid);
        console.error("Fetch error:", error);
        return;
      });
  }
}

function stackTransaction(url, objid) {
  console.log(url);
  let newTrans = {};
  newTrans.url = url;
  newTrans.objid = objid;
  newTrans.sent = false;

  const stackx = sessionStorage.getItem(myStackItem);
  let stack = [];
  if (stackx != null) {
    stack = JSON.parse(stackx);
  }
  stack.push(newTrans);
  sessionStorage.setItem(myStackItem, JSON.stringify(stack));
  /*
  obj.classList.remove("oi");
  obj.classList.add("oc");
  */
}

function swapconfig(obj) {
  let art = obj.parentElement.parentElement;
  let arts = document.querySelectorAll("article");
  console.log("arts == ", arts.length);
  for (let i = 0; i < arts.length; i++) {
    let flds = arts[i].querySelectorAll("fieldset");
    if (arts[i] === art) {
      console.log("showing ", flds.length);
      for (let j = 1; j < flds.length; j++) {
        flds[j].classList.remove("hide");
      }
    } else {
      console.log("hifdng ", flds.length);
      for (let j = 1; j < flds.length; j++) {
        flds[j].classList.add("hide");
      }
    }
  }
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
  dec.value = parseInt(obj.getAttribute("data-result"));

  console.log("Closing " + obj.name + " ==" + dec.value);
  if (dec.value == 0) {
    applyCorrectAnswerBonus(true);
    if (obj.getAttribute("id") == "PercentPenalty") {
      applyPercentPenalty(true);
    }
  }

  let url = "/x?f=saveebc";
  let inps = frm.getElementsByTagName("input");
  for (let i = 0; i < inps.length; i++) {
    let nm = inps[i].getAttribute("name");
    if (nm && nm != "") {
      if (inps[i].getAttribute("type") != "radio" || inps[i].checked) {
        url += "&" + nm + "=" + encodeURIComponent(inps[i].value);
      }
    }
  }
  if (obj.getAttribute("id") == "PercentPenalty") {
    url += "&PercentPenalty=1";
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
        window.location.href = "/ebclist";
      }
    })
    .catch((error) => {
      console.error("Fetch error:", error);
    });

  //frm.submit();
}

function fixClaimTimeISO() {
  let iso = document.getElementById("ClaimTimeISO");
  let dt = document.getElementById("ClaimDate");
  let tm = document.getElementById("ClaimTime");

  iso.value = dt.value + "T" + tm.value + iso.value.substring(16);
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

function showQHotChanged(obj) {
  let hot = "";
  if (obj.checked) {
    hot = "hot";
  }
  reloadRankings("hot", hot);
}
function showQOkChanged(obj) {
  let ok = "";
  if (obj.checked) {
    ok = "ok";
  }
  reloadRankings("ok", ok);
}
function showQSpeedChanged(obj) {
  let speed = "";
  if (obj.checked) {
    speed = "speed";
  }
  reloadRankings("speed", speed);
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
  frm = document.getElementById("optionsfrm");
  inps = frm.getElementsByTagName("input");
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
  fetchEntrantDetails(e);

  let b = document.getElementById("BonusID");
  b.value = x[2];
  fetchBonusDetails(b);

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
      if (inps[i].getAttribute("type") != "radio" || inps[i].checked) {
        if (inps[i].getAttribute("type") == "checkbox" && !inps[i].checked) {
          url +=
            "&" +
            nm +
            "=" +
            encodeURIComponent(inps[i].getAttribute("data-unchecked"));
        } else {
          url += "&" + nm + "=" + encodeURIComponent(inps[i].value);
        }
      }
    }
  }

  console.log(url);
  //alert(url)
  //return
  fetch(url, { method: "POST" })
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
  let dec = document.getElementById("chosenDecision");

  dec.value = obj.options[obj.selectedIndex].value;
}

function applyPercentPenalty(apply) {
  let pts = document.getElementById("Points");
  if (!pts) return;
  if (pts.getAttribute("data-pm") == "m") return; // Can't discount mulitipliers
  let pv = parseInt(pts.value);
  let qv = document.getElementById("valPercentPenalty");
  if (!qv) return;
  let qvv = parseInt(qv.value);
  let points2deduct = Math.floor((qvv / 100) * pv);
  let points2return = pv - Math.ceil(((100 - qvv) / 100) * pv);

  if (apply) pv -= points2deduct;
  else pv += points2return;
  pts.value = pv;
}

function applyCorrectAnswerBonus(apply) {
  let qa = document.getElementById("QuestionAnswered");
  let qpts = parseInt(qa.getAttribute("data-pts"));
  let ptsinp = document.getElementById("Points");
  let pts = parseInt(ptsinp.value);
  if (apply) {
    pts += qpts;
  } else {
    pts -= qpts;
  }
  ptsinp.value = pts;
}

// this function is called while processing an EBC claim.
// its effect is to leave the claim undecided and move it
// to the end of the queue.
function leaveUndecided() {
  let lu = document.getElementById("leavebutton");
  closeEBC(lu);
}

function loadPage(pg) {
  console.log("loadPage called with '" + pg + "'");
  window.location.href = pg;
}

// span includes img and input
function toggleButton(obj) {
  let spn = obj.parentElement;
  let inp = spn.getElementsByTagName("input");
  inp.checked = !inp.checked;
  spn.classList.toggle("selected");
}
