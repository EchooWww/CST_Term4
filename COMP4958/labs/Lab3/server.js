const http = require("http");

http
  .createServer((_, res) => {
    console.log("The server received a request ");
    res.writeHead(200, {
      "Content-Type": "text/html",
      "access-control-allow-origin": "*",
    });
    res.end("Hello <b>World!</b>");
  })
  .listen(8888);
