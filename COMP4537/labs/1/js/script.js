import userMessages from "../lang/messages/en/user.js";

const addButton = document.querySelector("#addButton");
const backButton = document.querySelector("#backButton");
const noteContainer = document.querySelector("#noteContainer");
const bodyId = document.querySelector("body").id;
let notes = JSON.parse(localStorage.getItem("notes")) || [];

document.addEventListener("DOMContentLoaded", () => {
  if (bodyId === "home-body") {
    renderIndex();
  } else if (bodyId === "writer-body") {
    renderWriter();
  } else if (bodyId == "reader-body") {
    renderReader();
  }
});

function renderIndex() {
  document.querySelector("#mainHeading").innerHTML =
    userMessages.home.mainHeading;
  document.querySelector("#subHeading").innerHTML =
    userMessages.home.subHeading;
  document.querySelector("#reader").innerHTML =
    userMessages.home.readerLinkText;
  document.querySelector("#writer").innerHTML =
    userMessages.home.writerLinkText;
}

function renderWriter() {
  document.querySelector("#timestamp").innerHTML = `${
    userMessages.writer.savedText
  }: ${new Date().toLocaleTimeString()}`;
  renderNotes();
  addButton.innerHTML = userMessages.writer.addButton;
  backButton.innerHTML = userMessages.writer.backButton;
  addButton.addEventListener("click", () => {
    notes.push({ content: "" });
    renderNotes();
    saveNotes();
  });
  backButton.addEventListener("click", () => {
    window.location.href = "index.html";
  });
}

function renderReader() {
  notes = JSON.parse(localStorage.getItem("notes")) || [];
  renderNotes();
  updateTimestamp();
  setInterval(() => {
    notes = JSON.parse(localStorage.getItem("notes")) || [];
    renderNotes();
    updateTimestamp();
  }, 2000);
  backButton.innerHTML = userMessages.reader.backButton;
  backButton.addEventListener("click", () => {
    window.location.href = "index.html";
  });
}

function updateTimestamp() {
  const now = new Date();
  document.querySelector("#timestamp").textContent = `${
    userMessages.reader.updateText
  }: ${now.toLocaleTimeString()}`;
}

function saveNotes() {
  localStorage.setItem("notes", JSON.stringify(notes));
  const now = new Date();
  timestamp.textContent = `${
    userMessages.writer.savedText
  }: ${now.toLocaleTimeString()}`;
}

function renderNotes() {
  noteContainer.innerHTML = "";
  notes.forEach((note, index) => {
    const noteDiv = document.createElement("div");
    noteDiv.classList.add("note");

    const textarea = document.createElement("textarea");
    textarea.value = note.content;
    textarea.addEventListener("input", (e) => {
      notes[index].content = e.target.value;
      saveNotes();
    });

    const removeButton = document.createElement("button");
    removeButton.textContent = "Remove";
    removeButton.addEventListener("click", () => {
      notes.splice(index, 1);
      saveNotes();
      renderNotes();
    });

    noteDiv.appendChild(textarea);
    if (bodyId == "writer-body") {
      noteDiv.appendChild(removeButton);
    }
    noteContainer.appendChild(noteDiv);
  });
}
