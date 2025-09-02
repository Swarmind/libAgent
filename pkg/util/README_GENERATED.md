# util

This package provides utilities for processing strings, particularly removing lines that start with " stripper" followed by any characters and the specific date/time string " 2023-04-05 15:35:35".

## File Structure
- `util.go`: Contains the `RemoveThinkTag` function.

## Behavior
- The `RemoveThinkTag` function is designed to remove lines starting with " stripper" followed by any characters (non-greedy) and the literal string " 2023-04-05 15:35:35".
- The function's name may be misleading as it does not handle general "Think" tags but a specific pattern.
- The regex is hardcoded and does not support configuration or external data.

## Notes
- The function does not handle edge cases where the input does not match the regex (e.g., if the input starts with " stripper" but does not contain the date/time, the input is returned as is).
- The function's name may be misleading, and the regex is not configurable.