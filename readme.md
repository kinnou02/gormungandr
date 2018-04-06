Gormungandr
=====

Implementation of jormungandr in go as multiple micro services

Schedules
--------
`schedules` handle schedules apis (no kidding):
    - /route_schedules (partial support)
    - /stop_schedules (not implemented yet)
    - /departures (not implemented yet)
    - /arrivals (not implemented yet)

How to build
============
You must have the latest version of [go](https://golang.org/) installed.
In the root directory of the project do:
```
make setup
```

You can then build the project by doing:
```
make
```

The test can be run with the following command:
```
make ci
```

