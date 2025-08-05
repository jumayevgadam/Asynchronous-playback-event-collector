# Asynchronous Playback Event Collector

## Introduction
We need to manage the event processing department at Belet and are responsible for processing playback events that come from users. You regularly receive a batch of events that have pre-filled information about the playback event.

## Task

We need to write an implementation of the `EventCollector` interface. This interface can be found in the `collector.go` file. The service must process the events passed to the `Handle(...)` method and return an operation that is responsible for the progress of this processing.

Each event contains incomplete information upon receipt: it does not contain the user's region and device model. Before processing this event, this information must be completed. You can get the information in a separate service.

The implementation will be tested on various test cases.
The tests check the speed of the program and the correct collection of statistics on events.

A task that passes all checks within a limited time is considered successful.

## Instructions and details

To start the task, you need to clone the project to yourself, then create a separate
branch using the template `feature/<username>`, running the command:

```bash
git checkout -b feature/<username>
```

After that, you can start the task.

**Please note that for the tests to work properly, it is necessary to assign a working instance of your implementation to the CurrentCollector variable in the `init()` function of the `collector_impl.go` file, replacing the original `collectorNOOP`.**

It is recommended to put the implementation itself in the `collector_impl.go` file.

---

## ✅ Test Check
For intermediate testing, the file `collector_test.go` contains a simple test that checks the basic functionality. You can rely on it and/or add your own additional tests if they can help at the development stage.
Please check your function with given tests

### Run
```bash
  go test --race ./...
```

---

## Event Model
The event itself contains the following information:

* `id` — unique event identifier;
* `user_id` — user identifier;
* `video_id` — video identifier;
* `start_at`, `stop_at` — playback start and stop times;
* `bitrate_kbps` — playback bitrate;
* `device_type` — device type;
* `error_code` — playback error code;

When receiving an event, the user region and device model are empty. It is necessary to request this information, then fill the events with the received data and then pass them on for further processing.

It is also necessary to collect statistics for each batch of events in order to understand the workload.

Also, to monitor the work, it is necessary to be able to return event collector work statistics:

* How many events have been processed so far;
* Event processing speed per second.

There is no need to implement the `EnrichClient` interface yourself.

---

## 🧹 Code Quality

### Run
```bash
   golangci-lint run ./...
```
