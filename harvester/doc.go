// Package harvester harvest log files.
// It will register itself onto the central server when started.
//
// Once the central server requested for a log file,
// the harvester will read (like `tail -f`) the file
// then send the new contents to the central server.
package harvester
