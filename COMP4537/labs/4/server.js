const http = require("http");
const url = require("url");
const messages = require("./lang/en/en");

class Definition {
  constructor(word, definition) {
    this.word = word;
    this.definition = definition;
  }
}
class Server {
  constructor(port) {
    this.port = port || 80;
    this.dictionary = [];
    this.requestCount = 0;
  }

  handleRequest(req, res) {
    const parsedUrl = url.parse(req.url, true);
    let query = parsedUrl.query;
    if (parsedUrl.pathname === "/api/definitions") {
      if (req.method === "GET") {
        if (Object.keys(query).length === 0) {
          return this.handleUsage(res);
        }
        this.searchHandler(query, res);
      } else if (req.method === "POST") {
        let body = "";
        req.on("data", (chunk) => {
          body += chunk;
        });
        req.on("end", () => {
          try {
            const data = JSON.parse(body);
            this.storeHandler(data, res);
          } catch (error) {
            return this.badRequestHandler(res, "Invalid JSON format.");
          }
        });
      } else if (req.method === "OPTIONS") {
        res.writeHead(204, {
          "Access-Control-Allow-Origin": "*",
          "Access-Control-Allow-Methods": "POST, GET, OPTIONS",
          "Access-Control-Allow-Headers": "Content-Type",
        });
        return res.end();
      }
    } else {
      this.notFoundHandler(res);
    }
  }

  searchHandler(query, res) {
    this.requestCount++;
    const word = query.word || "";
    const definition = this.dictionary.find((d) => d.word === word);
    if (!word) {
      return this.badRequestHandler(res, "Invalid word input.");
    }
    const response = {
      requestCount: this.requestCount,
      definition: definition,
    };
    if (definition) {
      this.handleSuccess(res, response);
    } else {
      this.handleWarning(res, messages.wordNotFound.replace("%1", word));
    }
  }

  storeHandler(data, res) {
    this.requestCount++;
    const word = data.word || "";
    const definition = data.definition || "";
    if (!word || !definition) {
      return this.badRequestHandler(res, messages.badRequest);
    }
    if (this.dictionary.find((entry) => entry.word === word)) {
      return this.handleWarning(
        res,
        messages.wordAlreadyExists.replace("%1", word)
      );
    } else {
      const newItem = new Definition(word, definition);
      this.dictionary.push(newItem);
      this.handleSuccess(res, {
        ...newItem,
        wordCount: this.dictionary.length,
      });
    }
  }

  notFoundHandler(res, message = messages.notFound) {
    this.handleError(res, 404, message);
  }

  badRequestHandler(res, message) {
    this.handleError(res, 400, message);
  }

  handleUsage(res) {
    res.writeHead(200, {
      "Content-Type": "text/html",
      "access-control-allow-origin": "*",
    });
    res.end(messages.usage);
  }

  handleWarning(res, message) {
    const response = JSON.stringify({
      type: "warning",
      requestCount: this.requestCount,
      message: message,
    });
    res.writeHead(200, {
      "Content-Type": "application/json",
      "access-control-allow-origin": "*",
    });
    res.end(response);
  }

  handleError(res, status, message) {
    const response = JSON.stringify({
      type: "error",
      requestCount: this.requestCount,
      message: message,
    });
    res.writeHead(status, {
      "Content-Type": "application/json",
      "access-control-allow-origin": "*",
    });
    res.end(response);
  }

  handleSuccess(res, data) {
    const response = JSON.stringify({
      type: "success",
      requestCount: this.requestCount,
      ...data,
    });
    res.writeHead(200, {
      "Content-Type": "application/json",
      "access-control-allow-origin": "*",
    });
    res.end(response);
  }

  start() {
    http
      .createServer((req, res) => this.handleRequest(req, res))
      .listen(this.port, "0.0.0.0", () => {
        console.log(`Server running on port ${this.port}`);
      });
  }
}

const server = new Server(process.env.PORT);
server.start();
