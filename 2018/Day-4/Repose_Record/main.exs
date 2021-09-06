defmodule Log do
  def parse([log | logs], buffer_map) do
    capture = Regex.named_captures(~r/^\[(?<date>1518-\d+-\d+) (?<time>\d+:\d+)\] (?<msg> *(.+))$/, log) 
    pair = %{capture["time"] => capture["msg"]}
    buffer_map = Map.put(buffer_map, capture["date"],
      if buffer_map[capture["date"]] do
        Map.merge(buffer_map[capture["date"]], pair)
      else pair end)
    parse(logs, buffer_map)
# ...
# [1518-11-01 00:30] falls asleep
# [1518-11-01 00:55] wakes up
# [1518-11-01 23:58] Guard #99 begins shift
# [1518-11-02 00:40] falls asleep
# [1518-11-02 00:50] wakes up
# ...

    # See the ERROR above!!!

    # This function parse will parse the log
    # and store the below line
    # ```[1518-11-01 23:58] Guard #99 begins shift```
    # under the date of 1518-11-01 in our buffer_map.
    # But the messages related to this gaurd (Gaurd #99)
    # will be saved under 1518-11-02.

# To tackle we have to make a MAP
# with keys as IDs of gaurds and
# values as date, asleep/awake time.

# But if the situation is like `[1518-10-31 23:58]`
# then this data must be associated with the remaining
# message under [1518-11-01 ... : ...] for the corresponding
# guard.

# Also as we are thinking of saving the data
# such as time, date, messages associated to
# the ID of guard as key, we may not need to
# solve the above situation.

# Instead, we'll have to design a new data-structure here!

  end

  def parse([], buffer_map), do: buffer_map
end

File.stream!("./input.txt")
|> Stream.map(&String.trim/1)
|> Enum.to_list
|> Log.parse(%{})
|> IO.inspect
