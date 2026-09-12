# bytetime

A command-line calculator for the question that comes up every time you're
staring at a copy progress bar or planning a backup window: given two of
size, rate, and duration, what's the third?

`bytetime` doesn't move any bytes itself. It just does the arithmetic, and
it parses the same unit strings you'd already write by hand - `4.7GB`,
`650MiB`, `25MB/s`, `100Mbps`, `1h30m` - so you don't have to convert units
in your head first.

## Usage

Give it exactly two of `--size`, `--rate`, `--duration` and it solves for
the one you left out.

How long will a 4.7GB disc image take over a 25MB/s link?

```
$ bytetime --size 4.7GB --rate 25MB/s
4.70 GB at 25.00 MB/s takes 3m8s
```

What throughput do I need to move 650MiB in 90 seconds?

```
$ bytetime --size 650MiB --duration 90s
681.57 MB at 7.57 MB/s takes 1m30s
```

How much data fits down a 100Mbps line in an hour?

```
$ bytetime --rate 100Mbps --duration 1h
45.00 GB at 12.50 MB/s takes 1h0m0s
```

Add `--json` to any of the above for a machine-readable result instead:

```
$ bytetime --size 4.7GB --rate 25MB/s --json
{
  "size_bytes": 4700000000,
  "size_human": "4.70 GB",
  "rate_bytes_per_sec": 25000000,
  "rate_human": "25.00 MB/s",
  "duration_seconds": 188,
  "duration_human": "3m8s"
}
```

### Size units

Decimal (`KB`, `MB`, `GB`, `TB`, `PB`, powers of 1000) and binary (`KiB`,
`MiB`, `GiB`, `TiB`, `PiB`, powers of 1024) are both accepted. A bare number
with no unit is read as bytes.

### Rate units

Two forms: a size followed by `/s` (`25MB/s`, `1.5GiB/s`) for bytes per
second, or a decimal prefix followed by `bps` (`100Mbps`, `56kbps`) for bits
per second, matching how network speeds are usually advertised. Bit rates
are converted to bytes internally by dividing by 8.

### Durations

Anything Go's `time.ParseDuration` accepts: `90s`, `3m20s`, `1h30m`, and so
on, plus `d` and `w` for days and weeks, since backup windows and retention
periods are usually thought about that way: `2d`, `1w`, `1w3d12h`.

## Building

Standard library only, no dependencies to fetch:

```
go build -o bytetime .
```

## License

MIT, see LICENSE.
