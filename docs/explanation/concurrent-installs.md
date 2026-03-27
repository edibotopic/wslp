## Concurrent installs

Install multiple distros faster.

### Experimental concurrency flag

When installing multiple distros, you have the option to
install them concurrently with `wslp`:

```shell
wslp install --experimental-concurrent Ubuntu Debian archlinux
```

A default limit of 3 is set for simultaneous installations.
When a slot is freed up, another distro will begin installing if it is waiting.

### Performance

Below are the results of two experiments.

Each experiment compared bulk install of distros using
`wslp` with (concurrent) or without (sequential) the concurrency flag.

The speed-ups may be useful if you frequently create and teardown WSL instances.

#### Experiment 1: Three distros

_Tested with: Ubuntu-24.04, archlinux, Debian_

| Mode        | Run 1 | Run 2 | Run 3 | Average |
|-------------|------:|------:|------:|--------:|
| Concurrent  | 108.13 | 90.53 | 107.49 | 102 |
| Sequential  | 155.25 | 163.66 | 143.49 | 154 |

**Result**: 34% reduction in installation time

#### Experiment 2: Five distros

_Tested with: Ubuntu-24.04, archlinux, Debian_, kali-linux, FedoraLinux-42 

| Mode        | Run 1 | Run 2 | Average |
|-------------|------:|------:|--------:|
| Concurrent  | 153.83 | 145.67 | 150 |
| Sequential  | 299.92 | 244.21 | 272 |

**Result**: 45% reduction in installation time

#### Testing approach

The following command was run in PowerShell:

```powershell
Measure-Command { wslp.exe install <distro-1> <distro-2> <distro-n> --experimental-concurrent }
```
