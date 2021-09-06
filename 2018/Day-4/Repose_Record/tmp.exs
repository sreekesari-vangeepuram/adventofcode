defmodule ElfAction do
  defstruct year: "", month: "", day: "", hour: "", minute: "", guard: -1, state: -1, action_string: ""
end

defmodule ElfSleepTimes do
  defstruct guard: -1, start_minute: -1, end_minute: -1
end

defmodule AdventOfCode2018.Day04 do
  @moduledoc false

  @doc """
  --- Day 4: Repose Record ---

  You've sneaked into another supply closet - this time, it's across from the 
  prototype suit manufacturing lab. You need to sneak inside and fix the issues 
  with the suit, but there's a guard stationed outside the lab, so this is as 
  close as you can safely get.

  As you search the closet for anything that might help, you discover that 
  you're not the first person to want to sneak in. Covering the walls, someone 
  has spent an hour starting every midnight for the past few months secretly 
  observing this guard post! They've been writing down the ID of the one guard 
  on duty that night - the Elves seem to have decided that one guard was enough 
  for the overnight shift - as well as when they fall asleep or wake up while 
  at their post (your puzzle input).

  For example, consider the following records, which have already been 
  organized into chronological order:

  [1518-11-01 00:00] Guard #10 begins shift
  [1518-11-01 00:05] falls asleep
  [1518-11-01 00:25] wakes up
  [1518-11-01 00:30] falls asleep
  [1518-11-01 00:55] wakes up
  [1518-11-01 23:58] Guard #99 begins shift
  [1518-11-02 00:40] falls asleep
  [1518-11-02 00:50] wakes up
  [1518-11-03 00:05] Guard #10 begins shift
  [1518-11-03 00:24] falls asleep
  [1518-11-03 00:29] wakes up
  [1518-11-04 00:02] Guard #99 begins shift
  [1518-11-04 00:36] falls asleep
  [1518-11-04 00:46] wakes up
  [1518-11-05 00:03] Guard #99 begins shift
  [1518-11-05 00:45] falls asleep
  [1518-11-05 00:55] wakes up

  Timestamps are written using year-month-day hour:minute format. The guard 
  falling asleep or waking up is always the one whose shift most recently 
  started. Because all asleep/awake times are during the midnight hour 
  (00:00 - 00:59), only the minute portion (00 - 59) is relevant for those 
  events.

  Visually, these records show that the guards are asleep at these times:

  Date   ID   Minute
              000000000011111111112222222222333333333344444444445555555555
              012345678901234567890123456789012345678901234567890123456789
  11-01  #10  .....####################.....#########################.....
  11-02  #99  ........................................##########..........
  11-03  #10  ........................#####...............................
  11-04  #99  ....................................##########..............
  11-05  #99  .............................................##########.....

  The columns are Date, which shows the month-day portion of the relevant day; 
  ID, which shows the guard on duty that day; and Minute, which shows the 
  minutes during which the guard was asleep within the midnight hour. (The 
  Minute column's header shows the minute's ten's digit in the first row and 
  the one's digit in the second row.) Awake is shown as ., and asleep is shown 
  as #.

  Note that guards count as asleep on the minute they fall asleep, and they 
  count as awake on the minute they wake up. For example, because Guard #10 
  wakes up at 00:25 on 1518-11-01, minute 25 is marked as awake.

  If you can figure out the guard most likely to be asleep at a specific time, 
  you might be able to trick that guard into working tonight so you can have 
  the best chance of sneaking in. You have two strategies for choosing the best 
  guard/minute combination.

  Strategy 1: Find the guard that has the most minutes asleep. What minute does 
  that guard spend asleep the most?

  In the example above, Guard #10 spent the most minutes asleep, a total of 50 
  minutes (20+25+5), while Guard #99 only slept for a total of 30 minutes 
  (10+10+10). Guard #10 was asleep most during minute 24 (on two days, whereas 
  any other minute the guard was asleep was only seen on one day).

  While this example listed the entries in chronological order, your entries 
  are in the order you found them. You'll need to organize them before they 
  can be analyzed.

  What is the ID of the guard you chose multiplied by the minute you chose? 
  (In the above example, the answer would be 10 * 24 = 240.)
  """
  def part1(file_path) do
    sleep_times = file_path
    |> File.stream!()
    |> Stream.map(&String.trim/1)
    |> Enum.to_list()
    |> data_to_structs()
    |> sort_structs()
    |> generate_sleep_times()

    sleepiest_guard = sleep_times
    |> sumarise_sleep_times()
    |> Enum.reduce({-1, -1}, fn guard_sleep_total, max ->
      new_max = if elem(guard_sleep_total, 1) > elem(max, 1) do
        { elem(guard_sleep_total, 0), elem(guard_sleep_total, 1) }
      else
        max
      end
    end)
    |> elem(0)

    {optimal_minute, _occurrence_count} = sleepiest_minute_for_guard(sleep_times, sleepiest_guard)

    sleepiest_guard * optimal_minute
  end

  @doc """
  Strategy 2: Of all guards, which guard is most frequently asleep on the same 
  minute?

  In the example above, Guard #99 spent minute 45 asleep more than any other 
  guard or minute - three times in total. (In all other cases, any guard spent 
  any minute asleep at most twice.)

  What is the ID of the guard you chose multiplied by the minute you chose? 
  (In the above example, the answer would be 99 * 45 = 4455.)
  """
  def part2(file_path) do
    sleep_times_for_guards = file_path
    |> File.stream!()
    |> Stream.map(&String.trim/1)
    |> Enum.to_list()
    |> data_to_structs()
    |> sort_structs()
    |> generate_sleep_times()

    {guard_id, minute, _occurrence} = sleep_times_for_guards
    |> Enum.reduce(MapSet.new(), fn sleep_time, acc ->
      MapSet.put(acc, sleep_time.guard)
    end)
    |> Enum.reduce(%{}, fn guard_id, acc ->
      Map.put(acc, guard_id, sleepiest_minute_for_guard(sleep_times_for_guards, guard_id))
    end)
    |> Enum.reduce({0, 0, 0}, fn guard_info, {acc_guard_id, acc_minute, acc_count} = acc ->
      {guard_id, {minute_number, occurrence_count}} = guard_info
      if occurrence_count > acc_count do
        {guard_id, minute_number, occurrence_count}
      else
        acc
      end
    end)

    guard_id * minute
  end

  defp data_to_structs(data_lines) do
    Enum.map(data_lines, fn data_line -> 
      elf_action_from_string(data_line)
    end)
  end

  defp sort_structs(structs_list) do
    Enum.sort(structs_list, fn x, y -> 
      date_time_x = x.year <> x.month <> x.day <> x.hour <> x.minute
      date_time_y = y.year <> y.month <> y.day <> y.hour <> y.minute
      date_time_x < date_time_y
    end)
  end

  # The accumulator in Enum.reduce is a tuple of:
  # - The Guard ID
  # - The Sleep Start Time
  # - A list of ElfSleepTime structs
  defp generate_sleep_times(structs_list) do
    Enum.reduce(structs_list, {0, -1, []}, fn elf_action, elf_state -> 
      elf_state_info(elf_action, elf_state)
    end)
    |> elem(2)
  end

  defp sumarise_sleep_times(sleep_times_list) do
    Enum.reduce(sleep_times_list, %{}, fn sleep_times, acc ->
      sleep_duration = sleep_times.end_minute - sleep_times.start_minute
      Map.update(acc, sleep_times.guard, sleep_duration, &(&1 + sleep_duration))
    end)
  end

  defp sleepiest_minute_for_guard(sleep_times_for_all_guards, guard_id) do
    sleep_times_for_all_guards
    |> Enum.filter(fn elf_sleep_time ->
      elf_sleep_time.guard == guard_id
    end)
    |> Enum.reduce(%{}, fn sleep_time, acc ->
      Enum.reduce(sleep_time.start_minute..sleep_time.end_minute, acc, fn minute, acc ->
        Map.update(acc, minute, 1, &(&1 + 1))
      end)
    end)
    |> Enum.reduce({0, 0}, fn {day, count}, acc -> 
      {acc_day, acc_count} = acc
      if count > acc_count do
        {day, count}
      else
        if count == acc_count and day < acc_day do
          {day, count}
        else
          acc
        end
      end
    end)
  end

  defp minutes_in_sleep_time(sleep_time) do
    Enum.reduce(sleep_time.start_minute..sleep_time.end_minute, %{}, fn minute, acc ->
      Map.update(acc, minute, 1, &(&1 + 1))
    end)
  end

  defp elf_state_info(%ElfAction{action_string: "falls asleep"} = elf_action, elf_state) do
    elf_state
    |> put_elem(1, String.to_integer(elf_action.minute))
  end

  defp elf_state_info(%ElfAction{action_string: "wakes up"} = elf_action, elf_state) do
    elf_id = elem(elf_state, 0)
    sleep_time = elem(elf_state, 1)
    wake_time = String.to_integer(elf_action.minute)

    elf_state = put_elem(elf_state, 1, -1)

    elf_sleeps = elem(elf_state, 2)
    elf_sleep = %ElfSleepTimes{guard: elf_id, start_minute: sleep_time, end_minute: wake_time}
    elf_state = put_elem(elf_state, 2, [elf_sleep] ++ elf_sleeps)
    elf_state
  end

  defp elf_state_info(elf_action, elf_state) do
    elf_id = elf_action.action_string
    |> String.slice(7..-14)
    |> String.to_integer()

    elf_state
    |> put_elem(0, elf_id)
    |> put_elem(1, -1)
  end

  defp elf_action_from_string(data_line) do
    year = String.slice(data_line, 1..4)
    month = String.slice(data_line, 6..7)
    day = String.slice(data_line, 9..10)
    hour = String.slice(data_line, 12..13)
    minute = String.slice(data_line, 15..16)
    action = String.slice(data_line, 19..-1)

    %ElfAction{year: year, month: month, day: day, hour: hour, minute: minute, action_string: action}
  end

end

IO.inspect AdventOfCode2018.Day04.part1("./input.txt")
IO.inspect AdventOfCode2018.Day04.part2("./input.txt")
