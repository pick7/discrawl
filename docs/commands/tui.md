# `tui`

Opens the local terminal archive browser for stored messages.

## Usage

```bash
discrawl tui
discrawl tui --guild 123456789012345678 --channel general
discrawl tui --guilds 123,456 --author 1456464433768300635
discrawl tui --dm
discrawl --json tui --limit 50
```

## What it shows

The browser uses the shared crawlkit explorer:

- left pane: channel, person, or thread groups
- middle pane: newest matching message rows
- right pane: selected message detail, attachments, replies, and thread context
- footer: local DB or remote Git snapshot source

Mouse selection, right-click actions, sortable headers, refresh, and chat layout match the other crawlkit-backed archive tools.

## Flags

- `--guild <id>` / `--guilds <id,id>` - restrict the guild scope
- `--dm` - browse local direct messages under the synthetic `@me` guild
- `--channel <id|name|#name>` - restrict to one channel or DM conversation
- `--author <id|name>` - restrict to one author
- `--limit <n>` - newest rows to load (default 200)
- `--include-empty` - include rows with no displayable/searchable content
- `--json` - print crawlkit browser rows as JSON instead of opening the TUI

## Notes

- `tui` is read-only.
- without `--guild`, `--guilds`, or `--dm`, it uses `default_guild_id` when configured; otherwise it can browse all stored guild rows
- `--dm` only shows messages imported from the local Discord Desktop cache by [`wiretap`](wiretap.html)
- `--json` is useful for launchers and agents that want the same row shape without an interactive terminal

## See also

- [`messages`](messages.html)
- [`dms`](dms.html)
- [`wiretap`](wiretap.html)
