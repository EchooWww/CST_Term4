defmodule Lec4 do
  def create(n\\5) do
    Enum.map(1..n, fn _ -> parent=self(); # before we spawn, we need to capture the parent process, which is the shell, and then pass it to the loop function
        spawn(fn -> loop(parent) end)
    end)
  end

  def trap_exit(pid) do
    send(pid, :trap_exit)
  end

  def no_trap_exit(pid) do
    send(pid, :no_trap_exit)
  end

  def link(pid1, pid2) do
    send(pid1, {:link, pid2})
  end

  def info(pid) when is_pid(pid) do
    if Process.alive?(pid) do
      {pid, {:alive, true}, Process.info(pid, :trap_exit),
      Process.info(pid, :links)}
    else
      {pid, {:alive, false}}
    end
  end

  def info(pids) do
    Enum.map(pids, fn pid -> info(pid) end)
  end

  def loop(parent) do
    receive do
      :trap_exit ->
        Process.flag(:trap_exit, true)
      :no_trap_exit ->
        Process.flag(:trap_exit, false)
      {:link, other} ->
        Process.link(other)
      msg ->
        send(parent, {self(), msg})
    end
    loop(parent)
  end
end
