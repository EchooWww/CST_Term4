defmodule Arithmetic.Worker do
  use GenServer
  def start() do
    GenServer.start(__MODULE__, nil)
  end

  def square(pid, x) do
    GenServer.call(pid, {:square, x})
  end

  def sqrt(pid, x) do
    GenServer.call(pid, {:sqrt, x})
  end

  @impl true
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_call({:square, x}, _from, state) do
    {:reply,{self(), x * x}, state}
  end

  @impl true
  def handle_call({:sqrt, x}, _from, state) do
    Process.sleep(4000)
    reply = if x >=0, do: :math.sqrt(x), else: :error
    {:reply, {self(),reply}, state}
  end
end

defmodule Arithmetic.Server do
  use GenServer
  def start(n\\1)
    when is_integer(n) and n > 0 do
      GenServer.start(__MODULE__, n, name: __MODULE__)
    end

  def select_worker() do
    GenServer.call(__MODULE__, {:select_worker})
  end

  def square(x) do
    worker = select_worker()
    Arithmetic.Worker.square(worker, x)
  end

  def sqrt(x) do
    worker = select_worker()
    Arithmetic.Worker.sqrt(worker, x)
  end

  @impl true
  def init(num_workers) do
    workers =
      for _ <- 1..num_workers do
        {:ok, pid} = Arithmetic.Worker.start()
        IO.puts("Worker started with PID: #{inspect(pid)}")
        pid
      end
    {:ok, %{workers: workers, next: 0}}
  end

  @impl true
  def handle_call({:select_worker}, _from, state) do
    worker = Enum.at(state.workers, state.next)
    next = rem(state.next + 1, length(state.workers))
    {:reply, worker, %{state | next: next}}
  end
end
