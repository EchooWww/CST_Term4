# Lecture 7: CRUD Application with Phoenix

## Initial Setup

1. Create new Phoenix project with Postgres:

```bash
mix phx.new <app name> --database postgres
```

mix 2. Configure Database

- Edit `config/dev.exs` to set database credentials
- Create database:

```bash
mix ecto.create
```

3. Start Phoenix server:

```bash
mix phx.server
```

## Creating CRUD Application

1. Generate CRUD resources:

```bash
mix phx.gen.html School Student students sid:string:unique firstname:string last:string score:integer
```

- Context: `School`
- Schema: `Student` (singular)
- Table: `students` (plural)
- Fields: `sid`, `firstname`, `last`, `score`

2. Configure Routes

- Add to `router.ex` after the line of `get "/", PageController, :home`

```elixir
resources "/students", StudentController
```

3. Verify and Apply Changes

```bash
mix phx.routes        # Verify routes
mix ecto.migrate      # Create database tables
mix phx.server        # Start server
```

4. We can also go to `school.ex` and add a new function to do other operations like `update`, `delete`, etc. It is the "context", which provides the API to the controller.

5. The schema is in `school/student.ex` and the controller is in `school/student_controller.ex`. There's a `changeset` function in the schema that validates the data before it is inserted into the database. We can try calling it in `iex`:

```elixir
iex -S mix
alias Bcit.School.Student
Student.changeset(%Student{}, %{}) # Returns a changeset struct with errors
```

We can add more rules to the `changeset` function to validate the data.

```elixir
    |> validate_number(:score, greater_than_or_equal_to: 0, less_than_or_equal_to: 100, message: "Score must be between 0 and 100")
```

We can also transform the data before validation in `changeset`:

```elixir
    attrs = for {key, value} <- attrs, into: %{}, do: {key, String.trim(value)}
```

6. We can add authentication to the application by using `phx.gen.auth`:

```bash
mix phx.gen.auth Accounts User users
```

This will generate a new context for authentication, with a schema for `User` and a controller for `SessionController`. We can then add the routes to `router.ex`. Then we can move the resources for `students` to a new scope.

## PostgreSQL Commands

1. Connect to database:

```bash
psql -U postgres
```

2. Useful Meta-commands:
   | Command | Description |
   |---------|-------------|
   | `\l` | List databases |
   | `\c` | Connect to database |
   | `\dt` | List tables |
   | `\d table_name` | Describe table |
   | `\q` | Quit |

## Phoenix LiveView Process Flow

### 1. Initial HTTP Request (Static Mount)

When a user first visits a LiveView page:

- Browser makes an HTTP request
- Server responds with initial HTML
- Initial `mount/3` is called in static mode (`connected?(socket)` is false)
- Initial assigns are set
- First render occurs
- HTML is sent to client with LiveView specific data attributes

```elixir
def mount(_params, _session, socket) do
  if connected?(socket) do
    # Not run during static mount
  end
  {:ok, assign(socket, initial_state: "loading")}
end
```

### 2. WebSocket Connection

After initial HTML load:

- Browser detects LiveView data attributes
- JavaScript initializes LiveView client
- WebSocket connection is established
- Server creates a new LiveView process
- Second `mount/3` call occurs with WebSocket (`connected?(socket)` is true)
- Subscriptions and real-time setup occur

```elixir
def mount(_params, _session, socket) do
  if connected?(socket) do
    Phoenix.PubSub.subscribe(MyApp.PubSub, "room:lobby")
  end
  {:ok, assign(socket, state: "ready")}
end
```

### 3. Render Process

Rendering occurs:

- After both static and live mounts
- Any time state changes via `assign/2`
- Uses the HEEx template engine
- Has access to all assigns via `@` syntax

```elixir
def render(assigns) do
  ~H"""
  <div>
    <h1>Status: <%= @state %></h1>
    <%= for item <- @items do %>
      <div><%= item %></div>
    <% end %>
  </div>
  """
end
```

### 4. Event Handling

LiveView handles various types of events:

#### Client Events (handle_event)

- User interactions (clicks, forms, etc.)
- Triggered by `phx-` bindings
- Updates state and triggers re-render

```elixir
def handle_event("save", params, socket) do
  {:noreply, assign(socket, saved: true)}
end
```

#### Server Events (handle_info)

- PubSub messages
- System messages
- Process monitoring
- External service events

```elixir
def handle_info({:new_message, message}, socket) do
  messages = [message | socket.assigns.messages]
  {:noreply, assign(socket, messages: messages)}
end
```

### 5. State Management

State is managed through:

- Socket assigns
- Temporary assigns (for memory optimization)
- Updates trigger re-renders automatically

```elixir
# Regular assigns
{:ok, assign(socket, count: 0)}

# Temporary assigns
{:ok, socket, temporary_assigns: [messages: []]}
```

### 6. Lifecycle Events

LiveView processes can handle various lifecycle events:

#### Initialization

- `mount/3`: Initial setup and state
- Connection status checking
- Subscription setup

#### Runtime

- `handle_params/3`: URL params changes
- `handle_event/3`: Client-side events
- `handle_info/2`: Server messages

#### Termination

- `terminate/2`: Cleanup on process exit
- Subscription cleanup
- Resource release

### 7. Common Patterns

#### Subscriptions

```elixir
if connected?(socket) do
  Phoenix.PubSub.subscribe(MyApp.PubSub, "topic")
end
```

#### Broadcasting

```elixir
Phoenix.PubSub.broadcast(MyApp.PubSub, "topic", {:event, payload})
```

#### State Updates

```elixir
socket = assign(socket, key: value)
```

#### Error Handling

```elixir
def mount(_params, _session, socket) do
  case dangerous_operation() do
    {:ok, data} -> {:ok, assign(socket, data: data)}
    {:error, _reason} -> {:error, :unauthorized}
  end
end
```

### Best Practices

1. Always check connection status before subscriptions
2. Use temporary assigns for large collections
3. Keep render functions focused on presentation
4. Handle all potential message types in handle_info
5. Clean up resources in terminate
6. Use proper error handling in mount and events
7. Keep state updates atomic and minimal
