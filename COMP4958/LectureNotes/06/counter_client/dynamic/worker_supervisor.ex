defmodule Counter.WorkerSupervisor do
  def start_link(_) do
    DynamicSupervisor.start_link(__MODULE__, nil, name: {:global, __MODULE__})
  end

  def start_worker(name) do
    child_spec = %{
      id: Counter.Worker,
      start: {Counter.Worker, :start_link, [name]}
    }
    DynamicSupervisor.start_child({:global, __MODULE__}, child_spec) # this is equivalent to Supervisor.child_spec({Counter.Worker, name}), which calls Counter.Worker.start_link(name)
  end
end
