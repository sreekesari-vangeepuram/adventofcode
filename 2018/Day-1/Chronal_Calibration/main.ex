# Read File
{:ok, file} = File.read("input.txt")

# Format data
frequencies = file
              |> String.split("\n", trim: true)
              |> Enum.map(&String.to_integer/1)

# Part 1
frequencies
|> Enum.reduce(&(&1 + &2))
|> (&IO.puts("Final resulting frequency: #{&1}")).()

# Part 2
Stream.cycle(frequencies)
|> Enum.reduce_while({0, MapSet.new([0])}, fn change, {current, duplicate} ->
  resulting = current + change
  cond do
    MapSet.member?(duplicate, resulting) -> {:halt, resulting}
    true -> {:cont, {resulting, MapSet.put(duplicate, resulting)}}
  end
end)
|> (&IO.puts("First duplicate frequency: #{&1}")).()
