# Session Enumeration

This example shows how to enumerate devices and their audio session

## Prerequisites

- Go 1.13 or later
- `go-ole` (https://github.com/go-ole/go-ole)

## Build

```console
go build
```

That's it. Then you'll get `SessionEnumeration.exe`. Note that your platform is not Windows, you need set the environment variable `GOOS='windows'`.

## Usage

```console
./SessionEnumeration
```

## Contributing

Bug reports and improving the documentation are welcome. (https://github.com/moutend/go-wca/issues)

The Windows Core Audio API was introduced Windows vista, so that the later than that version of Windows could run this example. However, I'm not sure because I've just tested this example on Windows 10 version 1909 at the moment.
