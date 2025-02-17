const http = require("http");
const url = require("url");
const messages = require("./lang/en/en");
const mysql = require("mysql2");

// Disclosure: Used ChatGPT 4o in hosting MySQL DB on Google Cloud and SSL connection
const dbConfig = {
  host: process.env.MYSQL_HOST,
  user: process.env.MYSQL_USER,
  password: process.env.MYSQL_PASSWORD,
  port: process.env.MYSQL_PORT,
  ssl: {
    ca: process.env.MYSQL_SSL_CA,
    key: process.env.MYSQL_SSL_KEY,
    cert: process.env.MYSQL_SSL_CERT,
  },
};

class Database {
  constructor() {
    this.connection = mysql.createConnection(dbConfig);
    this.dbname = "v7_lab5";
    this.initializeDB();
  }

  initializeDB() {
    this.connection.query(
      `CREATE DATABASE IF NOT EXISTS ${this.dbname}`,
      (err) => {
        if (err) throw err;
        this.connection.changeUser({ database: this.dbname }, (err) => {
          if (err) throw err;
          this.createTable();
        });
      }
    );
  }

  createTable() {
    const createTableQuery = `CREATE TABLE IF NOT EXISTS patient (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        dateOfBirth DATETIME NOT NULL
      ) ENGINE=InnoDB;`;
    this.connection.query(createTableQuery, (err) => {
      if (err) throw err;
    });
  }

  executeQuery(query, callback, params = []) {
    // Disclosure: Used ChatGPT 4o for Regex pattern generation, same for all other regex patterns in this file
    if (/^\s*(DROP|DELETE|UPDATE)\s/i.test(query)) {
      callback(messages.INVALID_QUERY, null);
      return;
    }
    this.createTable();
    this.connection.execute(query, params, (err, results) => {
      if (err) callback(err, null);
      else callback(null, results);
    });
  }
}

class Server {
  constructor() {
    this.port = process.env.PORT || 80;
    this.database = new Database();
  }

  handleRequest(req, res) {
    const parsedUrl = url.parse(req.url, true);
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader("Access-Control-Allow-Methods", "GET, POST");
    res.setHeader("Access-Control-Allow-Headers", "Content-Type");
    res.setHeader("Content-Type", "application/json");

    if (req.method === "GET" && parsedUrl.pathname.startsWith("/api/v1/sql")) {
      const queryStr = decodeURIComponent(
        parsedUrl.pathname.split("/api/v1/sql/")[1]
      );
      if (!/^\s*SELECT\s/i.test(queryStr)) {
        return this.handleError(res, 400, messages.INVALID_GET);
      }
      this.database.executeQuery(queryStr, (err, results) => {
        if (err) {
          return this.handleError(res, 500, err.message);
        }
        this.handleSuccess(res, results);
      });
    } else if (
      req.method === "POST" &&
      parsedUrl.pathname.startsWith("/api/v1/sql")
    ) {
      console.log("POST request received");
      let body = "";
      req.on("data", (chunk) => {
        body += chunk;
      });
      req.on("end", () => {
        try {
          const data = JSON.parse(body);
          const queryStr = data.query;
          if (!/^\s*INSERT\s/i.test(queryStr)) {
            return this.handleError(res, 400, messages.INVALID_POST);
          }
          this.database.executeQuery(queryStr, (err, results) => {
            if (err) {
              return this.handleError(res, 500, err.message);
            }
            this.handleSuccess(res, results);
          });
        } catch (error) {
          return this.handleError(res, 400, messages.INVALID_FORMAT);
        }
      });
    } else if (
      req.method === "POST" &&
      parsedUrl.pathname === "/api/v1/insertDefaultPatients"
    ) {
      const queryStr =
        "INSERT INTO patient (name, dateOfBirth) VALUES (?, ?), (?, ?), (?, ?), (?, ?)";
      const params = [
        "Sara Brown",
        "1901-01-01",
        "John Smith",
        "1941-01-01",
        "Jack Ma",
        "1961-01-30",
        "Elon Musk",
        "1999-01-01",
      ];
      this.database.executeQuery(
        queryStr,
        (err, results) => {
          if (err) {
            return this.handleError(res, 500, err.message);
          }
          this.handleSuccess(res, results);
        },
        params
      );
    } else if (req.method === "OPTIONS") {
      res.writeHead(200);
      res.end();
    } else {
      this.handleError(res, 404, "Not found");
    }
  }

  handleError(res, status, message) {
    const response = JSON.stringify({
      error: message,
    });
    res.writeHead(status);
    res.end(response);
  }

  handleSuccess(res, data) {
    const response = JSON.stringify({ success: true, data });
    res.writeHead(200);
    res.end(response);
  }

  start() {
    http
      .createServer((req, res) => this.handleRequest(req, res))
      .listen(this.port, () => {
        console.log(`Server running on port ${this.port}`);
      });
  }
}

const server = new Server();
server.start();
