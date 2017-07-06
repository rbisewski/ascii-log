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
    "regexp"
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

    // Name of the IP log file on the webserver.
    ip_log = "ip.log"

    // Name of the whois log file on the webserver.
    whois_log = "whois.log"

    // Parameter for the server type
    serverType = ""

    // Argument for enabling daemon mode
    daemonMode = false
)

// Initialize the argument input flags.
func init() {

    // Server type flag
    flag.StringVar(&serverType, "server-type", "nginx",
      "Currently active server; e.g. 'nginx' ")

    // Daemon mode flag
    flag.BoolVar(&daemonMode, "daemon-mode", false,
      "Whether or not to run this program as a background service.")
}

//
// PROGRAM MAIN
//
func main() {

    // String variable to hold eventual output, as well error variable.
    var err error = nil

    // Variable to hold the extracted IP addresses
    var ip_addresses = make(map[string] int)

    // Variable to hold the contents of the new ip.log and whois log until
    // it is written to disk.
    var ip_log_contents    = ""
    var whois_log_contents = ""

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
    access_log_location := log_directory + serverType + "/" + access_log

    // main infinite loop...
    for {

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

        // determine the last valid line
        last_line_num := len(lines)-2

        // safety check, ensure the value is at least zero
        if last_line_num < 0 {
            last_line_num = 0
        }

        // obtain the contents of the last line
        last_line := lines[last_line_num]

        // extract the date of the last line, this is so that the program can
        // gather data concerning only the latest entries
        latest_date_in_log, err := obtainLatestDate(last_line)

        // check if an error occurred
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // turn the latest data string into a regex
        re := regexp.MustCompile(latest_date_in_log)

        // for every line...
        for _, line := range lines {

            // verify that a match could be found
            verify := re.FindString(line)

            // skip a line if the entry is not the latest date
            if len(verify) < 1 {
                continue
            }

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

        // convert the ip addresses map into an array of strings
        ip_strings, err := convertIpAddressMapToString(ip_addresses)

        // if an error occurred, terminate from the program
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // attempt to stat() the ip.log file, else create it if it does
        // not currently exist
        err = statOrCreateFile(web_location + ip_log)

        // if an error occurred during stat(), yet the program was unable
        // to recover or recreate the file, then exit the program
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // append the title to the ip_log_contents
        ip_log_contents += "IP Address Counts Data\n\n"

        // attempt to grab the current day/month/year
        datetime := time.Now().Format(time.UnixDate)

        // append the date to the ip_log_contents on the next line
        ip_log_contents += "Generated on: " + datetime + "\n"
        ip_log_contents += "\n"
        ip_log_contents += "Log Data for " + latest_date_in_log + "\n"
        ip_log_contents += "-------------------------\n\n"

        // append the ip_strings content to this point of the log; it will
        // either contain the "IPv4 Address + Daily Count" or a message stating
        // that no addresses appear to be recorded today.
        ip_log_contents += ip_strings

        // attempt to write the string contents to the ip.log file
        err = ioutil.WriteFile(web_location + ip_log,
                               []byte(ip_log_contents),
                               0755)

        // if an error occurred, terminate the program
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // attempt to obtain the whois entries, as a string
        whois_strings, err := obtainWhoisEntries(ip_addresses)

        // if an error occurred, terminate the program
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // append the title to the whois_log_contents
        whois_log_contents += "Whois Entry Data\n\n"

        // append the date to the whois_log_contents on the next line
        whois_log_contents += "Generated on: " + datetime + "\n"
        whois_log_contents += "\n"
        whois_log_contents += "Log Data for " + latest_date_in_log + "\n"
        whois_log_contents += "-------------------------\n\n"

        // append the whois entry strings to the whois log contents
        whois_log_contents += whois_strings

        // attempt to stat() the whois.log file, else create it if it does
        // not currently exist
        err = statOrCreateFile(web_location + whois_log)

        // if an error occurred during stat(), yet the program was unable
        // to recover or recreate the file, then exit the program
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // attempt to write the string contents to the ip.log file
        err = ioutil.WriteFile(web_location + whois_log,
                               []byte(whois_log_contents),
                               0755)

        // if daemon mode is disabled, then exit this loop
        if !daemonMode {
            break
        }

        // take the current time and increment 24 hours
        currentTime := time.Now()
        oneDayLater := currentTime.AddDate(0,0,1)

        // since the user has selected daemon mode, wait 24 hours
        for {

            // grab the current time
            currentTime = time.Now()

            // if an entire day has passed, go ahead and break
            if currentTime.After(oneDayLater) {
                break
            }
        }
    }

    // If all is well, we can return quietly here.
    os.Exit(0)
}
