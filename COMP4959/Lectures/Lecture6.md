## Lecture 6: Phoenix Framework and Web Dev with Elixir

Create a new Phoenix project:

```bash
mix phx.new my_app
```

Change directory to the new project and start the server:

```bash
iex -S mix phx.server
```

The `config` directory contains configuration files for the project. The `router.ex` file contains the routes for the project. The `lib` directory contains the main application code. The `web` directory contains the web application code.

`priv/static` contains static files. It also has an `assets` directory which contains js and css files.

`plug` is a middleware for Phoenix. It is a way to modify the request and response. `%Plug.Conn{}` is a struct that represents the connection.

Every function in a controller takes a connection as an argument and returns a connection.

In the controller, we can either use a template or in the function.

To add a simple route, add the following to the `router.ex` file:

```elixir
get "/hello", HelloController, :hello
```

Then create a new controller:

```bash
defmodule HelloWeb.HelloController do
  use HelloWeb, :controller

  def hello(conn, _params) do
    text(conn, "Hello, World!")
  end
end
```

Then, add a `hello_html.ex` file

```elixir
defmodule HelloWeb.HelloHtml do
  use HelloWeb, :html
  embed_templates "hello_html/*"
end
```

Then create a `hello.html.eex` file in the `hello_html` directory:

```html
<h1>Hello, World!</h1>
```
