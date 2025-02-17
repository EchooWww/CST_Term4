const http = require("http");
const url = require("url");
const fs = require("fs");
const path = require("path");
const utils = require("./modules/utils");
const messages = require("./lang/en/en");

class Server {
  constructor(port) {
    this.port = port || 3000;
  }

  handleRequest(req, res) {
    const parsedUrl = url.parse(req.url, true);
    const query = parsedUrl.query;

    if (parsedUrl.pathname === "/COMP4537/labs/3/getDate/") {
      this.getDateHandler(query, res);
    } else if (parsedUrl.pathname === "/COMP4537/labs/3/writeFile/") {
      this.writeFileHandler(query, res);
    } else if (parsedUrl.pathname.includes("/COMP4537/labs/3/readFile")) {
      this.readFileHandler(parsedUrl, res);
    } else {
      this.notFoundHandler(res);
    }
  }

  getDateHandler(query, res) {
    const name = query.name || "";
    const time = utils.getDate();
    const output = messages.getDate.replace("%1", name).replace("%2", time);
    res.writeHead(200, {
      "Content-Type": "text/html",
      "access-control-allow-origin": "*",
    });
    res.end(`<h1 style="color:0076df">${output}</h1>`);
  }

  writeFileHandler(query, res) {
    const text = query.text || "";
    fs.appendFile(path.join(__dirname, "file.txt"), text, (err) => {
      if (err) {
        this.handleError(res, messages.writeFailed);
      } else {
        this.handleSuccess(res, messages.writeSuccess);
      }
    });
  }

  readFileHandler(parsedUrl, res) {
    const fileName = parsedUrl.pathname.replace(
      "/COMP4537/labs/3/readFile/",
      ""
    );
    const file = path.join(__dirname, fileName === "" ? "file.txt" : fileName);
    fs.readFile(file, "utf8", (err, data) => {
      if (err) {
        this.handleError(res, messages.fileNotFound.replace("%1", fileName));
      } else {
        res.writeHead(200, {
          "Content-Type": "text/html",
          "access-control-allow-origin": "*",
        });
        res.end(`<h1 style="color:0076df">${data}</h1>`);
      }
    });
  }

  notFoundHandler(res) {
    this.handleError(res, messages.notFound);
  }

  handleError(res, message) {
    res.writeHead(404, {
      "Content-Type": "text/html",
      "access-control-allow-origin": "*",
    });
    res.end(`<h1>${message}</h1>`);
  }

  handleSuccess(res, message) {
    res.writeHead(200, {
      "Content-Type": "text/html",
      "access-control-allow-origin": "*",
    });
    res.end(`<h1>${message}</h1>`);
  }

  start() {
    http
      .createServer((req, res) => this.handleRequest(req, res))
      .listen(this.port, () => {
        console.log(`Server running on port ${this.port}`);
      });
  }
}

const server = new Server(process.env.PORT || 3000);
server.start();
