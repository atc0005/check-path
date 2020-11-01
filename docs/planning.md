<!-- omit in toc -->
# check-path

<!-- omit in toc -->
## Table of Contents

- [Brain storming](#brain-storming)
  - [2020-10-18](#2020-10-18)
  - [2020-10-28](#2020-10-28)
    - [overview](#overview)
    - [common threads](#common-threads)
    - [existence check](#existence-check)
    - [age check](#age-check)
    - [size check](#size-check)
    - [owner, group check](#owner-group-check)
    - [permissions](#permissions)

## Brain storming

### 2020-10-18

How to report an individual file over a specific size vs an entire directory?

Does it matter? If a directory is too large, it is too large, whether a single
file or not. When a directory is flagged as a problem, perhaps the top X files
can be returned, or the top X subdirectories.

By default, if using os.Stat or os.Lstat an os.FileInfo will be returned which
contains:

```golang
FileInfo interface {
    Name() string       // base name of the file
    Size() int64        // length in bytes for regular files; system-dependent for others
    Mode() FileMode     // file mode bits
    ModTime() time.Time // modification time
    IsDir() bool        // abbreviation for Mode().IsDir()
    Sys() interface{}   // underlying data source (can return nil)
}
```

For the atc0005/elbow project I used `FileMatch` and `FileMatches` types:

```golang
// FileMatch represents a superset of statistics (including os.FileInfo) for a
// file matched by provided search criteria. This allows us to record the
// original full path while also recording file metadata used in later
// calculations.
type FileMatch struct {
  os.FileInfo
  Path string
}

// FileMatches is a slice of FileMatch objects that represents the search
// results based on user-specified criteria.
type FileMatches []FileMatch
```

in order to track a path and applicable metadata. From there I had multiple
methods to help process them in bulk. I will probably want to mirror much of
that functionality here.

We should probably not check the age of directories, instead limiting age
checks to just files within a directory path. We can then flag that directory
as a problem.

### 2020-10-28

#### overview

Currently the only "early exit" or "fail fast" logic being applied is for the
existence check. Everything else finishes the "crawl" of the first specified
path (no matter the depth) before evaluating the `FileMatches` entries:

- size
- age

Only at that point does the WARNING or CRITICAL thresholds apply. If a `path`
with millions of files is specified, this could result in a lot of wasted
time/effort by the app should the first (or even first thousand) files be
sufficient to cause the check to fail. Instead, we need to at least
periodically, if not for every file, check to see whether specified thresholds
are crossed before the entire collection of files for the specified path is
examined. If you take into account that *multiple* paths may be specified,
each with many, many files within, we really need to optimize for "fail fast".

#### common threads

Most of these points apply to all checks, with the notable exception of the
existence check. It is probably worth keeping that check separate from the
others; this choice makes sense because we treat it as incompatible with all
other checks.

- a collection of paths is specified

- recursive boolean flag *could* be specified for non-existence checks

- non-existence checks have potential need to walk a path

- all checks currently return a Nagios Exit State type
  - used at the end of execution to trigger the generation of a report/summary
    for Nagios' use
  - refactored checks could share the original as a pointer

- each check that is requested can register itself with a shared slice (e.g.,
  `checksApplied`) to be used with a final "app specified paths pass `%v`
  validation checks" line which is composed of a `strings.Join(checksApplied,
  ", ")` call

- one idea (probably won't implement) will be to setup a compatibility func
  that accepts registered checks thus far and the candidate check and returns
  a boolean value to indicate whether the candidate can be used with existing
  registered checks
  - e.g., age check and previously used existence check

#### existence check

- specified path
- path exists critical?
- path exists warning?
- if specified list of options is given, that is considered a config validation failure

outer logic is whether the result is critical or warning
inner logic is just whether the path exists, can fail early

#### age check

- specified path
- recursive?
- low age (warning)
- high age (critical)

We can ignore directories.

#### size check

- specified path
- recursive?
- low size (warning)
- high size (critical)

We can ignore directories.

#### owner, group check

- specified path
- recursive?
- username-warning
- username-critical
- group-name-warning
- group-name-critical

We must *not ignore* directories, but process them too alongside files.

#### permissions

Not implemented yet, needs further thought

- specified path
- recursive
- permissions (octal)
- permissions (select flags)
  - assert presence
    - `require-owner-write`
    - `require-owner-read`
    - `require-owner-execute`
    - `require-group-write`
    - `require-group-read`
    - `require-group-execute`
    - `require-other-write`
    - `require-other-read`
    - `require-other-execute`
  - assert missing
    - `no-owner-write`
    - `no-owner-read`
    - `no-owner-execute`
    - `no-group-write`
    - `no-group-read`
    - `no-group-execute`
    - `no-other-write`
    - `no-other-read`
    - `no-other-execute`

Perhaps a set of `require-*` flags and then an `invert-perms` flag to flip the
logic?
