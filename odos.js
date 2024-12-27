// @ts-checks

"use strict";

var timertick;
var reloadok = true;


function calcMileage() {
  let mlgdiv = document.getElementById("OdoMileage");
  if (!mlgdiv) return;
  let km = document.getElementById("OdoCountsK").checked;
  let units =
    parseInt(document.getElementById("OdoFinish").value) -
    parseInt(document.getElementById("OdoStart").value);
  if (isNaN(units) || units < 1) {
    mlgdiv.innerHTML = "";
    return;
  }
  if (km) {
    const KmsPerMile = 1.60934;
    units = parseInt(units / KmsPerMile);
  }
  mlgdiv.innerHTML = " " + units + " miles";
}



function clickTime() {
  let timeDisplay = document.querySelector("#timenow");
  console.log("Clicking time");
  clearInterval(timertick);
  if (timeDisplay.getAttribute("data-paused") != 0) {
    timeDisplay.setAttribute("data-paused", 0);
    timertick = setInterval(
      refreshTime,
      timeDisplay.getAttribute("data-refresh")
    );
    timeDisplay.classList.remove("held");
  } else {
    timeDisplay.setAttribute("data-paused", 1);
    timertick = setInterval(clickTime, timeDisplay.getAttribute("data-pause"));
    timeDisplay.classList.add("held");
  }
  console.log("Time clicked");
}

function clickTimeBtn(btn) {
  let timeDisplay = document.querySelector("#timenow");
  clickTime();
  if (btn.innerHTML == btn.getAttribute("data-hold")) {
    btn.innerHTML = btn.getAttribute("data-unhold");
    setTimeout(
      clickTimeBtnRelease,
      timeDisplay.getAttribute("data-pause"),
      btn
    );
  } else {
    btn.innerHTML = btn.getAttribute("data-hold");
  }
}

function clickTimeBtnRelease(btn) {
  btn.innerHTML = btn.getAttribute("data-hold");
}

function endEditEntrant() {
  let mode = document.getElementById("EditMode").value;

  window.location = "signin?mode=" + mode;
}

function fix2(n) {
  if (n < 10) {
    return "0" + n;
  }
  return n;
}

function getRallyTime(dt) {
  let yy = dt.getFullYear();
  let mm = dt.getMonth() + 1;
  let dd = dt.getDate();
  let dateString =
    yy + "-" + fix2(mm) + "-" + fix2(dd) + "T" + dt.toLocaleTimeString("en-GB");
  return dateString.substring(0, 16);
}

function loadPage(x) {
  window.location.href = x;
}


function nextTimeSlot() {
  let timeDisplay = document.querySelector("#timenow");
  if (!timeDisplay) return;
  let dt = new Date();
  let gap = parseInt(timeDisplay.getAttribute("data-gap"));
  let xtra = parseInt(timeDisplay.getAttribute("data-xtra"));
  let newdt = getRallyTime(dt);
  let curdt = timeDisplay.getAttribute("data-time");

  if (xtra > 0 && gap > 0) {
    dt = parseDatetime(curdt);
    dt = new Date(
      dt.getFullYear(),
      dt.getMonth(),
      dt.getDate(),
      dt.getHours(),
      dt.getMinutes() + gap
    );
    newdt = getRallyTime(dt);
    console.log("Choosing next slot " + newdt);
    xtra--;
    timeDisplay.setAttribute("data-xtra", xtra);
  }
  timeDisplay.setAttribute("data-time", newdt);
  let dateString = dt.toLocaleString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
  });
  let formattedString = dateString.replace(", ", " - ");
  timeDisplay.innerHTML = formattedString;
}

function nextButtonSlot() {
  //console.log("nextButtonSlot");
  let btn = document.getElementById("nextSlot");
  if (!btn) return;
  console.log("nBS ok");
  let timeDisplay = document.querySelector("#timenow");
  let dt = new Date();
  let dateString = dt.toLocaleString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
  });
  let formattedString = getRallyTime(dt);

  let gap = parseInt(timeDisplay.getAttribute("data-gap"));
  let xtra = parseInt(timeDisplay.getAttribute("data-xtra"));
  let newdt = getRallyTime(dt);
  let curdt = timeDisplay.getAttribute("data-time");

  if (xtra > 0 && gap > 0) {
    dt = parseDatetime(curdt);
    dt = new Date(
      dt.getFullYear(),
      dt.getMonth(),
      dt.getDate(),
      dt.getHours(),
      dt.getMinutes() + gap
    );
    newdt = getRallyTime(dt);
    if (formattedString >= newdt) {
      btn.classList.add("hide");
    } else {
      btn.classList.remove("hide");
    }
  } else {
    btn.classList.add("hide");
  }
  dateString = dt.toLocaleString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
  });
  formattedString = dateString.replace(", ", " - ");
  btn.innerHTML = formattedString;
}

// Called during Odo capture when input is entered

function oic(obj) {
  reloadok = false;
  // Checkbox handler
  obj.setAttribute("data-chg", "1");
  // autosave handler
  if (obj.timer) {
    clearTimeout(obj.timer);
  }
  obj.timer = setTimeout(saveData, 1000, obj);
}

function oid(obj) {
  obj.classList.remove("oc");
  obj.classList.add("oi");
  obj.setAttribute("data-chg", "1");
  // autosave handler
  if (obj.timer) {
    clearTimeout(obj.timer);
  }
  obj.timer = setTimeout(saveData, 3000, obj);
}

function oidcfg(obj) {
  obj.classList.remove("oc");
  obj.classList.add("oi");
  obj.setAttribute("data-chg", "1");
  // autosave handler
  if (obj.timer) {
    clearTimeout(obj.timer);
  }
  obj.timer = setTimeout(saveConfig, 3000, obj);
}

// Called during Odo capture when input is complete
function oc(obj) {
  saveOdo(obj);
}

function ocd(obj) {
  if (obj.getAttribute("data-chg") == "1") {
    console.log("ocd: " + obj.name);
    saveData(obj);
  }
}

function ocdcfg(obj) {
  if (obj.getAttribute("data-chg") == "1") {
    console.log("ocd: " + obj.name);
    saveConfig(obj);
  }
}

function parseDatetime(dt) {
  let yy = parseInt(dt.substring(0, 4));
  let mm = parseInt(dt.substring(5, 7)) - 1;
  let dd = parseInt(dt.substring(8, 10));
  let hh = parseInt(dt.substring(11, 13));
  let mn = parseInt(dt.substring(14, 16));
  return new Date(yy, mm, dd, hh, mn);
}

function refreshTime() {
  sendTransactions();
  let timeDisplay = document.querySelector("#timenow");
  if (!timeDisplay) return;
  let dt = new Date();
  let gap = parseInt(timeDisplay.getAttribute("data-gap"));
  let xtra = parseInt(timeDisplay.getAttribute("data-xtra"));
  let newdt = getRallyTime(dt);
  let curdt = timeDisplay.getAttribute("data-time");
 // console.log(
 //   "Comparing " + curdt + " > " + newdt + "; xtra=" + xtra + "(" + gap + ")"
 // );
  nextButtonSlot();
  if (curdt >= newdt) {
    return;
  }
  nextTimeSlot();
}

function reloadPage() {
  console.log("reloadPage called");
  if (!reloadok) {
    setTimeout(reloadPage, 1000);
    return;
  }

  let url = window.location.href;

  console.log("Reloading " + url);
  window.location = url;
}

function saveConfig(obj) {
  if (obj.getAttribute("data-static") == "") obj.setAttribute("data-chg", "");
  console.log("saveConfig: " + obj.name);
  if (obj.timer) {
    clearTimeout(obj.timer);
  }

  let val = obj.value;

  let url = encodeURI("config?" + obj.name + "=" + val);
  stackTransaction(url, obj.id);

  sendTransactions();
}

function saveData(obj) {
  if (obj.getAttribute("data-static") == "") obj.setAttribute("data-chg", "");
  console.log("saveData: " + obj.name);
  if (obj.timer) {
    clearTimeout(obj.timer);
  }

  let ent = document.getElementById("EntrantID").value;
  let val = obj.value;
  switch (obj.name) {
    case "RiderPostcode":
    case "PillionPostcode":
    case "BikeReg":
    case "RiderCountry":
    case "PillionCountry":
      val = val.toUpperCase();
      break;

    case "RiderIBA":
    case "PillionIBA":
    case "RiderRBLR":
    case "PillionRBLR":
    case "FreeCamping":
    case "CertificateAvailable":
    case "CertificateDelivered":
      val = "N";
      if (obj.checked) val = "Y";
      break;

    case "EntrantStatus":
      setTimeout(endEditEntrant, 1000);
      break;

    case "OdoStart":
    case "OdoFinish":
    case "OdoKms":
    case "OdoCounts":
      calcMileage();
      break;
  }

  let url = encodeURI(
    "putentrant?EntrantID=" + ent + "&" + obj.name + "=" + val
  );
  stackTransaction(url, obj.id);

  validate_entrant();
}

function saveFinalStatus(obj) {
  const FinisherOK = 8;

  let ent = obj.getAttribute("data-e");
  let val = obj.value;

  let url = "putentrant?EntrantID=" + ent + "&" + obj.name + "=" + val;
  if (val != FinisherOK) {
    url += "&CertificateAvailable=N&CertificateDelivered=N";
  }

  stackTransaction(encodeURI(url), obj.id);
}


function showMoneyAmt() {
  let amt = addMoney();
  let sf = document.getElementById("showmoney");
  if (sf) {
    sf.innerHTML = "£" + amt;
  }
}

function showPillionPresent() {
  let first = document.getElementById("PillionFirst");
  let last = document.getElementById("PillionLast");
  let present = first.value != "" && last.value != "";
  let ps = document.getElementById("showpillion");
  if (ps) {
    ps.innerHTML = "";
    if (present) ps.innerHTML = "&#9745;";
  }
}

function signin(m, e) {
  window.location = "/edit?m=" + m + "&e=" + e;
}


function validate_entrant() {
  let mustFields = [
    "RiderFirst",
    "RiderLast",
    "RiderEmail",
    "RiderPhone",
    "Bike",
    "NokName",
    "NokRelation",
    "NokPhone",
  ];
  let noktabFields = ["NokName", "NokRelation", "NokPhone"];

  mustFields.forEach((fld) => {
    let f = document.getElementById(fld);
    f.setAttribute("placeholder", "must not be blank");
    if (f.value == "") f.classList.add("notblank");
    else f.classList.remove("notblank");
  });
  let nokAlert = false;
  noktabFields.forEach((fld) => {
    let f = document.getElementById(fld);
    nokAlert = nokAlert || f.value == "";
  });
  let noktab = document.getElementById("noktab");
  if (noktab) {
    if (nokAlert) noktab.classList.add("notblank");
    else noktab.classList.remove("notblank");
  }
}
