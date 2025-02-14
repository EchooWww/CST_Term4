defmodule Counter.Worker do
  # use GenServer # we can remove this line because we are not implementing the GenServer callbacks in this module, we can still use the GenServer functions
  @name {:global, __MODULE__} # @name is a module attribute that is used to store the name of the GenServer process

  def start_link(n\\0) do
    Genserver.start_link(__MODULE__, n, name: @name) # start_link is a function that starts a process and links it to the current process
    # the first argument is the callback module
  end

  def inc(amt \\1) do
    GenServer.cast(@name, {:inc, amt})
  end

  def dec (amt \\1) do
    GenServer.cast(@name, {:dec, amt})
  end

  def value do
    GenServer.call(@name, :value)
  end
end
