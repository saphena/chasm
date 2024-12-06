// Custom Javascript for use with certificates.
// Not standalone, relies on Jodit editor being available in a separate file.

// @ts-check
"use strict";
const editorid = "#editor"

function enablesavecert() {
    let sb = document.getElementById('savecert')
    if (sb) sb.disabled = false
}
const liveFields = {
  Bike: "Bike",
  CrewName: "Crew name(s)",
  Distance: "Distance ridden",
  EntrantID: "Rider #",
  Place: "Finish place (1)",
  Rank: "Finish rank (1st)",
  MKlit: "'miles' | 'km'",
  RallyTitle: "Rally title",
  TeamName: "Team name",
  Points: "Points scored",
};
const myButtons = [
  "bold",
  "italic",
  "underline",
  "font",
  "fontsize",
  "brush",
  "align",
  "image",
  {
    name: "db",
    iconURL: "images/fields.png",
    list: liveFields,
    exec: (...args) => {
      for (let p in args[2]) {
        console.log(JSON.stringify(p));
      }
      console.log(JSON.stringify(args[2].control));
      if (args[2].control.hasOwnProperty("list")) return;
      let ed = args[0];
      let key = args[2].control.args[0];
      //let val = args[2].control.args[1];
      if (key == "") return;
      ed.selection.insertNode(ed.create.element("span", "{{." + key + "}}"));
    },
  },
];
const editor = Jodit.make(editorid, {
  width: "100%",
  height: "100%",
  statusbar: false,
  beautifyHTML: false,
  useAceEditor: false,

  buttons: myButtons,
  buttonsMD: myButtons,
  buttonsSM: myButtons,
  buttonsXS: myButtons,

  uploader: {
    insertImageAsBase64URI: true,
  },
});

editor.e.on('change', param => enablesavecert());