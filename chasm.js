// @ts-checks

"use strict";

const myStackItem = "odoStack";

const CAT_OrdinaryScoringRule = 0;
const CAT_DNF_Unless_Triggered = 1;
const CAT_DNF_If_Triggered = 2;
const CAT_PlaceholderRule = 3;
const CAT_OrdinaryScoringSequence = 4;

// This is the maximum number of axes or sets of categories
const NumCategoryAxes = 9;

const ordered_list_icon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-list-ol" viewBox="0 0 16 16">
  <path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5"/>
  <path d="M1.713 11.865v-.474H2c.217 0 .363-.137.363-.317 0-.185-.158-.31-.361-.31-.223 0-.367.152-.373.31h-.59c.016-.467.373-.787.986-.787.588-.002.954.291.957.703a.595.595 0 0 1-.492.594v.033a.615.615 0 0 1 .569.631c.003.533-.502.8-1.051.8-.656 0-1-.37-1.008-.794h.582c.008.178.186.306.422.309.254 0 .424-.145.422-.35-.002-.195-.155-.348-.414-.348h-.3zm-.004-4.699h-.604v-.035c0-.408.295-.844.958-.844.583 0 .96.326.96.756 0 .389-.257.617-.476.848l-.537.572v.03h1.054V9H1.143v-.395l.957-.99c.138-.142.293-.304.293-.508 0-.18-.147-.32-.342-.32a.33.33 0 0 0-.342.338zM2.564 5h-.635V2.924h-.031l-.598.42v-.567l.629-.443h.635z"/>
</svg>`

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
  let sav = document.getElementById("updatedb");
  if (sav && ok) sav.value = "Confirm DELETE";
  if (sav && !ok) sav.value = "";
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

function showBonus(b) {
  let bonus = b.toUpperCase();
  if (bonus == "") return;
  let url = "/bonus?b=" + encodeURIComponent(bonus);
  url += "&back=bonuses";
  window.location.href = url;
}

function showCombo(c) {
  let bonus = c;
  if (bonus == "") return;
  let url = "/combo?c=" + encodeURIComponent(bonus);
  url += "&back=combos";
  window.location.href = url;
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

  switch (obj.getAttribute("data-save")) {
    case "saveOdo":
      obj.timer = setTimeout(saveOdo, 3000, obj);
      break;
    case "saveBonus":
      obj.timer = setTimeout(saveBonus, 3000, obj);
      break;
    default:
      obj.timer = setTimeout(saveSetupConfig, 3000, obj);
  }
  /*
  if (obj.getAttribute("data-save") == "saveOdo")
    obj.timer = setTimeout(saveOdo, 3000, obj);
  else obj.timer = setTimeout(saveSetupConfig, 3000, obj);
  */
  console.log("oi complete " + JSON.stringify(obj));
}

function addBonus(obj) {
  let b = obj.value.toUpperCase();
  let bd = document.getElementById("BriefDesc");
  console.log('addBonus called with "' + b + '"');
  if (b == "") {
    bd.value = "Blank code!";
    return;
  }
  let url = "/x?f=addb&b=" + encodeURIComponent(b);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        // Handle HTTP errors
        bd.value = `HTTP error! Status: ${response.status}`;
        return;
      }
      return response.json();
    })
    .then((data) => {
      if (data.err) {
        // Handle JSON error field
        console.error(`Error: ${data.msg}`);
        bd.value = `Error: ${data.msg}`;
        return;
      } else if (data.ok) {
        // Process the data if no error
        window.location.href = "/bonus?b=" + encodeURIComponent(b);
      } else {
        bd.value = `Error: ${data.msg}`;
        bd.setAttribute("title", bd.value);
      }
    })
    .catch((error) => {
      // Handle network or other errors
      console.error("Fetch error:", error);
      bd.value = "Fetch error";
      return;
    });
}

function addCombo(obj) {
  let b = obj.value;
  let bd = document.getElementById("BriefDesc");
  console.log('addCombo called with "' + b + '"');
  if (b == "") {
    bd.value = "Blank code!";
    return;
  }
  let url = "/x?f=addco&b=" + encodeURIComponent(b);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        // Handle HTTP errors
        bd.value = `HTTP error! Status: ${response.status}`;
        return;
      }
      return response.json();
    })
    .then((data) => {
      if (data.err) {
        // Handle JSON error field
        console.error(`Error: ${data.msg}`);
        bd.value = `Error: ${data.msg}`;
        return;
      } else if (data.ok) {
        // Process the data if no error
        window.location.href = "/combo?c=" + encodeURIComponent(b);
      } else {
        bd.value = `Error: ${data.msg}`;
        bd.setAttribute("title", bd.value);
      }
    })
    .catch((error) => {
      // Handle network or other errors
      console.error("Fetch error:", error);
      bd.value = "Fetch error";
      return;
    });
}

function saveBonus(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let b = obj.getAttribute("data-b");
  let url = "/x?f=saveb&b=" + b;
  let nm = obj.name;
  let ov = obj.value;
  if (nm == "ScoringFlag") {
    let fs = obj.parentElement.parentElement;
    let flgs = fs.querySelectorAll("input[name=ScoringFlag]");
    let fx = "";
    for (let f = 0; f < flgs.length; f++) {
      if (flgs[f].checked) fx = fx + flgs[f].value;
    }
    nm = "Flags";
    ov = fx;
  }
  url += "&ff=" + nm + "&" + nm + "=" + encodeURIComponent(ov);
  console.log("saveBonus: " + url);
  stackTransaction(url, obj.id);
  sendTransactions();
}

function saveCombo(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let b = obj.getAttribute("data-c");
  let url = "/x?f=saveco&c=" + b;
  let nm = obj.name;
  let ov = obj.value;
  url += "&ff=" + nm + "&" + nm + "=" + encodeURIComponent(ov);
  console.log("saveCombo: " + url);
  stackTransaction(url, obj.id);
  sendTransactions();
  if (nm == "MinimumTicks") {
    extractComboPointsArray();
    updateComboPointsList();
  }
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
  let inp = spn.getElementsByTagName("input")[0];
  inp.checked = !inp.checked;
  spn.classList.toggle("selected");
  saveBonus(inp);
}

function updateBonusDB(obj) {
  let del = document.getElementById("enableDelete");
  let bonus = document.getElementById("BonusID");
  if (!del || !del.checked || !bonus || bonus.value == "") {
    obj.disabled = true;
    return;
  }
  let bd = document.getElementById("BriefDesc");
  let url = "/x?f=delb&b=" + encodeURIComponent(bonus.value);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        // Handle HTTP errors
        bd.value = `HTTP error! Status: ${response.status}`;
        return;
      }
      return response.json();
    })
    .then((data) => {
      if (data.err) {
        // Handle JSON error field
        console.error(`Error: ${data.msg}`);
        bd.value = `Error: ${data.msg}`;
        return;
      } else if (data.ok) {
        // Process the data if no error
        window.location.href = "/bonuses";
      } else {
        bd.value = `Error: ${data.msg}`;
        bd.setAttribute("title", bd.value);
      }
    })
    .catch((error) => {
      // Handle network or other errors
      console.error("Fetch error:", error);
      bd.value = "Fetch error";
      return;
    });
}

function updateComboDB(obj) {
  let del = document.getElementById("enableDelete");
  let bonus = document.getElementById("ComboID");
  if (!del || !del.checked || !bonus || bonus.value == "") {
    obj.disabled = true;
    return;
  }
  let bd = document.getElementById("BriefDesc");

  let url = "/x?f=delco&c=" + encodeURIComponent(bonus.value);
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        // Handle HTTP errors
        bd.value = `HTTP error! Status: ${response.status}`;
        return;
      }
      return response.json();
    })
    .then((data) => {
      if (data.err) {
        // Handle JSON error field
        console.error(`Error: ${data.msg}`);
        bd.value = `Error: ${data.msg}`;
        return;
      } else if (data.ok) {
        // Process the data if no error
        window.location.href = "/combos";
      } else {
        bd.value = `Error: ${data.msg}`;
        bd.setAttribute("title", bd.value);
      }
    })
    .catch((error) => {
      // Handle network or other errors
      console.error("Fetch error:", error);
      bd.value = "Fetch error";
      return;
    });
}

// extractComboPointsArray takes the value of the PointsList comma-separated value string
// and creates the corresponding array of number fields. I should be called from the
// comma-separated text field.
function extractComboPointsArray() {
  console.log("extractComboPointsArray called");
  let BL = document.getElementById("BonusList");
  let PL = document.getElementById("PointsList");
  let hdrs = document.getElementById("PointsListArrayHdrs");
  let vals = document.getElementById("PointsListArrayVals");
  let mint = parseInt(document.getElementById("MinimumTicks").value);
  let xb = BL.value.split(",");
  let xv = PL.value.split(",");

  console.log(xv);

  // Zap the space
  hdrs.textContent = "";
  vals.textContent = "";

  let v = xv[0]; // starting value
  let maxt = xb.length;
  let maxv = xv.length;

  if (mint < 1 || mint == maxt) {
    // Special case, all bonuses compulsory so only one score value
    let n = document.createElement("input");
    n.type = "number";
    n.classList.add("Points");
    n.value = v;
    n.onchange = function () {
      updateComboPointsList();
    };
    vals.appendChild(n);
    return;
  }
  for (let i = 0; i < maxt - mint + 1; i++) {
    if (i < xv.length) v = xv[i];
    let n = document.createElement("input");
    n.type = "number";
    n.classList.add("Points");
    n.value = v;
    n.onchange = function () {
      updateComboPointsList();
    };
    vals.appendChild(n);
    let h = document.createElement("span");
    h.classList.add("Points");
    h.innerHTML = i + mint + "/" + maxt;
    hdrs.appendChild(h);
  }
}

function updateComboPointsList() {
  let PL = document.getElementById("PointsList");
  let vals = document.getElementById("PointsListArrayVals");
  let inps = vals.querySelectorAll("input");
  let x = "";
  for (let i = 0; i < inps.length; i++) {
    if (x != "") x += ",";
    x += "" + inps[i].value;
  }
  PL.value = x;
  saveCombo(PL);
}

// Category/set handling

function saveCatCat(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let s = obj.getAttribute("data-set");
  let c = obj.getAttribute("data-cat");
  let url = "/x?f=savecat&s=" + s;
  let nm = "setname";
  let ov = obj.value;
  url += "&c=" + c;
  url += "&ff=" + nm + "&" + nm + "=" + encodeURIComponent(ov);
  console.log("saveCatCat: " + url);
  stackTransaction(url, obj.id);
  sendTransactions();
}

function saveCatSet(obj) {
  if (obj.timer) clearTimeout(obj.timer);
  let b = obj.getAttribute("data-set");
  let url = "/x?f=saveset&s=" + b;
  let nm = "setname";
  let ov = obj.value;
  url += "&ff=" + nm + "&" + nm + "=" + encodeURIComponent(ov);
  console.log("saveCatSet: " + url);
  stackTransaction(url, obj.id);
  sendTransactions();
}

function showCatSet(obj) {
  let art = document.getElementById("setcats");
  let set = obj.getAttribute("data-set");
  let url = "/x?f=fetchcats&s=" + set;
  fetch(url)
    .then((response) => {
      if (!response.ok) {
        // Handle HTTP errors
        bd.value = `HTTP error! Status: ${response.status}`;
        return;
      }
      return response.json();
    })
    .then((data) => {
      if (data.err) {
        // Handle JSON error field
        console.error(`Error: ${data.msg}`);
        bd.value = `Error: ${data.msg}`;
        return;
      } else if (data.ok) {
        // Process the data if no error
        art.innerText = "";
        console.log(data);
        let div = document.createElement("div");
        div.innerText = `Categories for set ${data.Set} ${data.SetName}`;
        art.appendChild(div);
        let btn = document.createElement("button");
        btn.classList.add("plus");
        btn.setAttribute("data-set", `${data.Set}`);
        btn.onclick = function () {
          addCatCat(this);
        };
        btn.innerText = "+";
        art.appendChild(btn);
        for (let i = 0; i < data.Cats.length; i++) {
          let fs = document.createElement("fieldset");
          fs.classList.add("setcat");
          let lbl = document.createElement("label");
          let newid = `s${data.Set}c${data.Cats[i].Cat}`;
          lbl.setAttribute("for", newid);
          lbl.innerText = `${data.Cats[i].Cat}`;
          fs.appendChild(lbl);
          let inp = document.createElement("input");
          inp.setAttribute("id",newid)
          inp.classList.add("setcat");
          inp.setAttribute("value", `${data.Cats[i].CatDesc}`);
          inp.setAttribute("data-set", `${data.Set}`);

          inp.setAttribute("data-cat", `${data.Cats[i].Cat}`);
          inp.onchange = function () {
            saveCatCat(this);
          };
          fs.appendChild(inp);
          btn = document.createElement("button");
          btn.classList.add("minus");
          btn.setAttribute("data-set", `${data.Set}`);
          btn.setAttribute("data-cat", `${data.Cats[i].Cat}`);
          btn.onclick = function () {
            delCatCat(this);
          };
          btn.innerText = "-";
          fs.appendChild(btn);

          art.appendChild(fs);
        }
      } else {
        console.log(`Error: ${data.msg}`);
      }
    })
    .catch((error) => {
      // Handle network or other errors
      console.error("Fetch error:", error);
      return;
    });
}



// addCatCat needs to add the record in order to get a new number
function addCatCat(obj) {
  let dad = obj.parentElement;
  let sets = dad.querySelectorAll("input");
  let lastIx = sets.length;
  if (lastIx >= 0) {
    if (sets[lastIx - 1].value == "") return;
  }
  console.log(sets, lastIx);
  if (lastIx >= NumCategoryAxes) {
    obj.disabled;
    return;
  }
  let fs = document.createElement("fieldset");
  fs.classList.add("sethdr");
  let lbl = document.createElement("label");
  lastIx++; // Make it 1 relative
  let newid = "SetHdr" + lastIx;
  lbl.setAttribute("for", newid);
  lbl.innerText = "Set " + lastIx + " is";
  fs.appendChild(lbl);
  let inp = document.createElement("input");
  inp.setAttribute("id", newid);
  inp.setAttribute("name", newid);
  inp.setAttribute("data-set", lastIx);
  inp.onchange = function () {
    saveCatSet(this);
  };
  inp.onclick = function () {
    showCatSet(this);
  };
  fs.appendChild(inp);
  dad.appendChild(fs);
}

function addCatSet(obj) {
  let dad = obj.parentElement;
  let sets = dad.querySelectorAll("input");
  let lastIx = sets.length;
  if (lastIx >= 0) {
    if (sets[lastIx - 1].value == "") return;
  }
  console.log(sets, lastIx);
  if (lastIx >= NumCategoryAxes) {
    obj.disabled;
    return;
  }
  let fs = document.createElement("fieldset");
  fs.classList.add("sethdr");
  let lbl = document.createElement("label");
  lastIx++; // Make it 1 relative
  let newid = "SetHdr" + lastIx;
  lbl.setAttribute("for", newid);
  lbl.innerText = "Set " + lastIx + " is";
  fs.appendChild(lbl);
  let inp = document.createElement("input");
  inp.setAttribute("id", newid);
  inp.setAttribute("name", newid);
  inp.setAttribute("data-set", lastIx);
  inp.onchange = function () {
    saveCatSet(this);
  };
  inp.onclick = function () {
    showCatSet(this);
  };
  fs.appendChild(inp);
  dad.appendChild(fs);
}
