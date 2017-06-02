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

    // Web location
    web_location = "/var/www/html/data/"

    // Parameter for the server type
    serverType = ""
)

// Initialize the argument input flags.
func init() {

    // Server type flag
    flag.StringVar(&serverType, "server-type", "nginx",
      "Currently active server; e.g. 'nginx' ")
}

//
// PROGRAM MAIN
//
func main() {

    // String variable to hold eventual output, as well error variable.
    var err error     = nil

    // Parse the flags, if any.
    flag.Parse()

    // Lower case the serverType variable value.
    serverType = strings.ToLower(serverType)

    // Print the usage message if not nginx or apache.
    if (serverType != "nginx") && (serverType != "apache") {
        flag.Usage()
        os.Exit(1)
    }

    // Check if the web data directory actually exists.
    _, err = ioutil.ReadDir(web_location)

    // ensure no error occurred
    if err != nil {
        fmt.Println("The following directory does not exist: ",
          web_location)
        os.Exit(1)
    }

    // Assemble the access.log file location.
    access_log_location := log_directory + serverType + access_log

    // Check if access log file actually exists.
    byte_contents_of_access_log, err := ioutil.ReadFile(access_log_location)

    // ensure no error occurred
    if err != nil {
        fmt.Println("An error occurred while trying to read the " +
          "following file: ", access_log_location)
        fmt.Println(err)
        os.Exit(1)
    }

    // TODO: delete this
    byte_contents_of_access_log = byte_contents_of_access_log

    // If all is well, we can return quietly here.
    os.Exit(0)
}
