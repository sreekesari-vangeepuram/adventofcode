defmodule Claim do
  @moduledoc """
  `Claim` module: AoC-2018 Day-3 module
  """

  # <Private function>
  # Returns a tuple {id , cells} where the
  # cells are covered by the estimated claim!
  defp parse_to_cells(claim) do
    capture = Regex.named_captures(~r/^#(?<id>\d+) @ (?<x>\d+),(?<y>\d+): (?<w>\d+)x(?<h>\d+)$/, claim)
    [id, x, y, w, h] = Enum.map(
      [
        capture["id"],
        capture["x"],
        capture["y"],
        capture["w"],
        capture["h"],
      ], &String.to_integer/1)
    {id, for(dx <- x..x+w-1, dy <- y..y+h-1, do: {dx, dy})}
  end

  @doc """
  Finds all the overlapped cells and return
  them as a map {cell: lookup_value}!
  """
  def overlaps(claims) do
    Enum.reduce(claims, %{}, fn claim, fabric ->
      {_id, cells} = parse_to_cells(claim)
      Enum.reduce(cells, %{}, fn cell, acc ->
        # +1 Lookup values of cell -> overlap
        Map.update(acc, cell, 1, &(&1 + 1))
      end)
      |> Map.merge(fabric, fn _k, v1, v2 -> v1 + v2 end)
    end)
  end

  @doc """
  Returns the id of fabric with zero overlaps!
  """
  def seperated_fabric(overlapped_cells, parsed_data) do
    Enum.reduce_while(parsed_data, nil, fn claim, _ ->
      {id, cells} = parse_to_cells(claim)
      overlaps = MapSet.intersection(MapSet.new(cells), overlapped_cells)
      if MapSet.size(overlaps) != 0, do: {:cont, 0}, else: {:halt, id}
    end)
  end
end

# Read, Parse data from input-file
parsed_data = 
  File.stream!("./input.txt")
  |> Stream.map(&String.trim/1)
  |> Enum.to_list

# Format data from parsed-data
formated_data =
  parsed_data
  |> Claim.overlaps
  |> Map.to_list

# Part 1
formated_data
|> Enum.count(fn {_k, v} -> v > 1 end)
|> (&IO.puts("Number of overlaps of fabric cells: #{&1}")).()

# Part 2
formated_data
|> Enum.filter(fn {_k, v} -> v > 1 end) 
|> Enum.reduce(MapSet.new, fn {cell, _overlap}, acc -> MapSet.put(acc, cell) end)
|> Claim.seperated_fabric(parsed_data)
|> (&IO.puts("The only fabric-claim with no overlaps: #{&1}")).()
