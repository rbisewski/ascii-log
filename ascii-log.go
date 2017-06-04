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
    "time"
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

    // Name of the IPs file on the webserver.
    ip_log = "ip.log"

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
    var err error = nil

    // Variable to hold the extracted IP addresses
    var ip_addresses = make(map[string] int)

    // Variable to hold the contents of the new ip.log until it is
    // written to disk.
    var ip_log_contents = ""

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

    // dump the contents of the access log to a string
    string_contents_of_access_log := string(byte_contents_of_access_log)

    // if the contents are less than 1 byte, end the program
    if len(string_contents_of_access_log) < 1 {
        fmt.Println("Note: the following file was empty: ",
          access_log_location)
        fmt.Println("Exiting...")
        os.Exit(0)
    }

    // Attempt to break up the file into an array of strings a demarked by
    // the newline character.
    lines := strings.Split(string_contents_of_access_log, "\n")

    // terminate the program if the array has less than 1 element
    if len(lines) < 1 {
        fmt.Println("Warning: no line data was found, exiting...")
        os.Exit(0)
    }

    // for every line...
    for _, line := range lines {

        // attempt to split that line via spaces
        elements := strings.Split(line, " ")

        // safety check, ensure that element actually has a length
        // of at least 1
        if len(elements) < 1 {
            continue
        }

        // grab the first element, that is the IP address
        ip := elements[0]

        // ensure that the ip address is valid length
        //
        // 0.0.0.0 --> 8 chars (min)
        //
        // 123.123.123.123 --> 15 chars (max)
        //
        if len(ip) < 8 || len(ip) > 15 {
            continue
        }

        // Since the IP address is indeed roughly valid, go ahead and add
        // it to the global array.
        ip_addresses[ip]++
    }

    // do a natural sort on the array of IP addresses
    // TODO: implement this here

    // attempt check if the ip.log file exists in the "web_location"
    if _, err = os.Stat(web_location + ip_log); os.IsNotExist(err) {

        // if not, then create it
        f, creation_err := os.Create(web_location + ip_log)

        // if an error occurred during creation, terminate program
        if creation_err != nil {
            fmt.Println(creation_err)
            os.Exit(1)
        }

        // then go ahead and close the file connection for the time being
        f.Close()
    }

    // if an error occurred during stat() then exit the program
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // append the title to the ip_log_contents
    ip_log_contents += "IP Address Counts Data\n\n"

    // attempt to grab the current day/month/year
    datetime := time.Now().Format("2017-Jan-01 14:02")

    // append the date to the ip_log_contents on the next line
    ip_log_contents += "Log Generated on: " + datetime + "\n"
    ip_log_contents += "--------------------------------"

    // TODO: implement the below pseudo code

    // for every ip address

        // append that address + \t + count

    // if no ip addresses present, instead append a line about there being
    // no data for today.

    // attempt to write the string contents to the ip.log file

        // if an error occurred, terminate the program

    // If all is well, we can return quietly here.
    os.Exit(0)
}
