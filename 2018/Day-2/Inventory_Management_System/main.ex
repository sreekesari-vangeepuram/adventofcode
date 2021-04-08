defmodule ID do
  defp analyze(id) do
    xs = String.graphemes(id)
    Map.new(Enum.map(xs, fn x ->
      {x, Enum.count(xs, fn l -> l == x end)}
    end))
  end

  def parse([id | ids], map) do
    parse(ids, Map.put(map, id, analyze(id)))
  end

  def parse([], map), do: map

  def len(id), do: String.length(id)
  
  def keys_of_count(n, maps) do
    Enum.filter(for map <- maps do
      for {k, v} when v == n <- map, do: k
      end, fn x -> x != [] end)
  end

  def diff(id1, id2) do
    if len(id1) == len(id2) do
      Enum.zip(Enum.map([id1, id2], fn l -> String.graphemes(l) end))
      |> Enum.count(fn {l1, l2} -> l1 != l2 end)
    end
  end
end

# Read File
{:ok, file} = File.read("input.txt")

# Parse input data into required format
ids = String.split(file, "\n", trim: true)
parsed_data = ID.parse(ids, %{})

# Part 1
exactly_twice  = length ID.keys_of_count(2, Map.values(parsed_data))
exactly_thrice = length ID.keys_of_count(3, Map.values(parsed_data))

IO.puts "Checksum: #{exactly_twice} * #{exactly_thrice} = #{exactly_twice * exactly_thrice}"

# Part 2
for id1 <- ids do
  for id2 <- ids do
    if ID.len(id1) == ID.len(id2) do
      if ID.diff(id1, id2) == 1 do
        Enum.zip(Enum.map([id1, id2], fn l -> String.graphemes(l) end))
        |> Enum.filter(fn {l1, l2} -> l1 == l2 end)
        |> Enum.map(fn {l1, _} -> l1 end)
        |> Enum.join()
        |> (&IO.puts("Diff of IDs: #{&1}")).()
        exit(:shutdown)
      end
    end
  end
end
