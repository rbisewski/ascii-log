/*
 * ASCII Log Tool
 *
 * Description: A simple tool written in golang, for the purposes of
 *              generating ASCII graphs detailing the current number of IP
 *              addresses.
 *
 *              Specifically, this requires an active apache or nginx
 *              server with logs enabled.
 *
 * Author: Robert Bisewski <contact@ibiscybernetics.com>
 */

//
// Package
//
package main

//
// Imports
//
import (
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
)

//
// Globals
//
var (

    // Current location of the log directory.
    log_directory = "/var/log/"

    // Name of the access and error log files
    access_log = "access.log"
    error_log  = "error.log"

    // Parameter for the server type
    serverType = ""
)

// Initialize the argument input flags.
func init() {

    // LED device flag
    flag.StringVar(&serverType, "server-type", "nginx",
      "Currently active server; e.g. 'nginx' ")
}

//
// PROGRAM MAIN
//
func main() {

    // String variable to hold eventual output, as well error variable.
    var output string = ""
    var err error     = nil

    // Parse the flags, if any.
    flag.Parse()

    // Print the usage
    flag.Usage()

    // If all is well, we can return quietly here.
    os.Exit(0)
}
