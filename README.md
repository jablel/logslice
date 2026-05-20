# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to a new file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log > output.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | stdin |
| `--from` | Start of time range (RFC3339) | required |
| `--to` | End of time range (RFC3339) | required |
| `--format` | Timestamp format in log entries | `RFC3339` |
| `--field` | Timestamp field name (JSON logs) | `time` |

### Example with JSON logs

```bash
logslice --file service.log \
  --from "2024-01-15T08:00:00Z" \
  --to "2024-01-15T08:30:00Z" \
  --field timestamp
```

## How It Works

logslice uses binary search to efficiently locate the start of the target time range, avoiding the need to scan the entire file. This makes it significantly faster than `grep` on large log files.

## Building from Source

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build ./...
```

## License

MIT — see [LICENSE](LICENSE) for details.