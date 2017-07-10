//
// Utility functions for the 
//

//
// Package
//
package main

//
// Imports
//
import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "strings"
    "sort"
    "strconv"
)

//! Convert a file into a string array as per a given separator
/*
 * @param     string      /path/to/file
 * @param     string      tokenizer character sequence
 *
 * @return    string[]    array of lines
 */
func tokenizeFile(filepath string, separator string) ([]string,
  error) {

    // input validation
    if len(filepath) < 1 || len(separator) < 1 {
        return nil, fmt.Errorf("tokenizeFile() --> invalid input")
    }

    // Check if access log file actually exists.
    byte_contents, err := ioutil.ReadFile(filepath)

    // if an error occurs at this point, it is due to the program being
    // unable to access a read the file, so pass back an error
    if err != nil {
        return nil, fmt.Errorf("tokenizeFile() --> An error occurred " +
          "while trying to read the following file: ", filepath)
    }

    // dump the contents of the file to a string
    string_contents := string(byte_contents)

    // if the contents are less than 1 byte, mention that via error
    if len(string_contents) < 1 {
        return nil, fmt.Errorf("tokenizeFile() --> the following file " +
          "was empty: ", filepath)
    }

    // attempt to break up the file into an array of strings
    str_array := strings.Split(string_contents, separator)

    // terminate the program if the array has less than 1 element
    if len(str_array) < 1 {
        return nil, fmt.Errorf("tokenizeFile() --> no string data was " +
          "found")
    }

    // having obtained the lines of data, pass them back
    return str_array, nil
}

//! Stat if a given file exists at specified path, else create it.
/*
 * @param     string    /path/to/filename
 *
 * @return    error     error message, if any
 */
func statOrCreateFile(path string) error {

    // input validation, ensure the file location is sane
    if len(path) < 1 {
        return fmt.Errorf("statOrCreateFile() --> invalid input")
    }

    // variable declaration
    var fileNotFoundAndWasCreated bool = false

    // attempt to stat() if the whois.log file even exists
    _, err := os.Stat(path)

    // attempt check if the file exists at the given path
    if os.IsNotExist(err) {

        // if not, then create it
        f, creation_err := os.Create(path)

        // if an error occurred during creation, terminate program
        if creation_err != nil {
            return creation_err
        }

        // then go ahead and close the file connection for the time being
        f.Close()

        // if the program go to actually create the file, go ahead and
        // set this flag to true
        fileNotFoundAndWasCreated = true
    }

    // if an error occurred during stat(), yet the program was unable
    // to recover or recreate the file, then exit the program
    if err != nil && !fileNotFoundAndWasCreated {
        return err
    }

    // else everything worked, so go ahead and return nil
    return nil
}

//! Determine the latest date present in the logs
/*
 * @param     string    line data
 *
 * @return    string    latest time-date, in the form of DD/MMM/YYYY
 *            error     error message, if any
 */
func obtainLatestDate(line_data string) (string, error) {

    // input validation
    if len(line_data) < 1 {
        return "", fmt.Errorf("obtainLatestDate() --> invalid input")
    }

    // variable declaration
    var result string = ""

    // attempt to split that line via spaces
    elements := strings.Split(line_data, " ")

    // safety check, ensure there are at least 4 elements
    if len(elements) < 4 {

        // otherwise send back an error
        return "", fmt.Errorf("obtainLatestDate() --> poorly formatted line")
    }

    // attempt to grab the fourth element
    datetime := elements[3]

    // attempt to trim the string of [ and ] brackets
    datetime = strings.Trim(datetime, "[]")

    // ensure the string actually has at least a length of 1
    if len(datetime) < 1 {
        return "", fmt.Errorf("obtainLatestDate() --> date-time string " +
          "is of improper length")
    }

    // split the string via the ':' characters
    time_pieces := strings.SplitAfter(datetime, ":")

    // ensure there is at least one element
    if len(time_pieces) < 1 || len(time_pieces[0]) < 1 {
        return "", fmt.Errorf("obtainLatestDate() --> unable to use _:_ " +
          "chars to separate time into pieces")
    }

    // trim away the remaining : chars
    result = strings.Trim(time_pieces[0], ":")

    // final safety check, ensure that the result has a len > 0
    if len(result) < 1 {

        // otherwise send back an error
        return "", fmt.Errorf("obtainLatestDate() --> unable to " +
          "assemble string result")
    }

    // if everything turned out fine, go ahead and return
    return result, nil
}

//! Convert the global IP address map to an array of sorted ipEntry objects
/*
 * @param     map        string map containing ip addresses and counts
 *
 * @return    string     ip address + count strings, with newlines
 *                       separating them.
 *            error      error message, if any
 */
func convertIpAddressMapToString(ip_map map[string] int) (string, error) {

    // input validation
    if len(ip_map) < 1 {
        return "", fmt.Errorf("convertIpAddressMapToString() --> " +
          "invalid input")
    }

    // variable declaration
    var ip_strings string     = ""
    var tmp_str_array         = make([]string, 0)
    var lines_appended uint   = 0
    var first_hostname string = ""

    // for every IPv4 address in the given map...
    for ip, _ := range ip_map {

        // workaround, to better sort IP addresses
        if strings.Index(ip, ".") == 1 {
            ip = "00" + ip
        } else if strings.Index(ip, ".") == 2 {
            ip = "0" + ip
        }

        // append that address to the temp string array
        tmp_str_array = append(tmp_str_array, ip)
    }

    // sort the given list of IPv4 addresses
    sort.Strings(tmp_str_array)

    // for every ip address
    for _, ip := range tmp_str_array {

        // workaround, trim away any LHS zeros
        ip = strings.TrimLeft(ip, "0")

        // grab the count
        count := ip_map[ip]

        // take the given IP address and attempt to grab the hostname
        hostnames, err := net.LookupAddr(ip)

        // default to "N/A" as the default hostname if an error occurred
        // or no hostnames could be currently found...
        if err != nil || len(hostnames) < 1 {
            first_hostname = "N/A"

        // default to "N/A" as the default hostname if the hostname is
        // blank or currently NXDOMAIN and etc.
        } else if len(hostnames[0]) < 1 {
            first_hostname = "N/A"

        // Otherwise go ahead and use the first available hostname
        } else {
            first_hostname = hostnames[0]
        }

        // since the \t character tends to get mangled easily, add a buffer
        // of single-space characters instead to the IPv4 addresses
        space_formatted_ip_address := ip
        for len(space_formatted_ip_address) < 16 {
            space_formatted_ip_address += " "
        }

        // append that count | address |  hostname
        ip_strings += strconv.Itoa(count) + "\t"
        ip_strings += " | "
        ip_strings += space_formatted_ip_address
        ip_strings += " | "
        ip_strings += first_hostname
        ip_strings += "\n"

        // add a line counter for internal use
        lines_appended++
    }

    // if no ip addresses present, instead append a line about there being
    // no data for today.
    if lines_appended == 0 {
        ip_strings += "No IP addressed listed at this time."
    }

    // everything worked fine, so return the completed string contents
    return ip_strings, nil
}

//! Convert the global IP address map to string containing whois entries
/*
 * @param     map        string map containing ip addresses and counts
 *
 * @return    string     whois + ip, with hyphens and newlines
 *                       separating them.
 *            error      error message, if any
 *
 *
 * TODO: this function could use more testing
 */
func obtainWhoisEntries(ip_map map[string] int) (string, error) {

    // input validation
    if len(ip_map) < 1 {
        return "", fmt.Errorf("obtainWhoisEntries() --> invalid input")
    }

    // variable declaration
    var whois_strings string  = ""
    var entries_appended uint = 0
    var tmp_str_array         = make([]string, 0)
    var tmp_str_buffer        = ""
    var trimmed_string        = ""
    var err error
    var result bytes.Buffer

    // for every IPv4 address in the given map...
    for ip, _ := range ip_map {

        // workaround, to better sort IP addresses
        if strings.Index(ip, ".") == 1 {
            ip = "00" + ip
        } else if strings.Index(ip, ".") == 2 {
            ip = "0" + ip
        }

        // append that address to the temp string array
        tmp_str_array = append(tmp_str_array, ip)
    }

    // sort the given list of IPv4 addresses
    sort.Strings(tmp_str_array)

    // for every ip address
    for _, ip := range tmp_str_array {

        // workaround, trim away any LHS zeros
        ip = strings.TrimLeft(ip, "0")

        // safety check, skip to the next entry if this one is of length
        // zero
        if len(ip) < 1 {
            continue
        }

        // attempt to obtain the whois record
        result, err = runWhoisCommand(ip)

        // if an error occurs, break out of the loop
        if err != nil {
            break
        }

        // convert the byte buffer to a string
        tmp_str_buffer = result.String()

        // if no record is present, pass back a "N/A"
        if len(tmp_str_buffer) < 1 || tmp_str_buffer == "<nil>" {
            whois_strings += "Whois Entry for the following: "
            whois_strings += ip
            whois_strings += "\n"
            whois_strings += "N/A\n\n"
            whois_strings += "---------------------\n\n"
            continue
        }

        // trim it to remove potential whitespace
        trimmed_string = strings.Trim(tmp_str_buffer, " ")

        // ensure it still has a length of zero
        if len(trimmed_string) < 1 {
            whois_strings += "Whois Entry for the following: "
            whois_strings += ip
            whois_strings += "\n"
            whois_strings += "N/A\n\n"
            whois_strings += "---------------------\n\n"
            continue
        }

        // otherwise it's probably good, then go ahead and append it
        whois_strings += "Whois Entry for the following: "
        whois_strings += ip
        whois_strings += "\n"
        whois_strings += trimmed_string
        whois_strings += "\n\n"
        whois_strings += "---------------------\n\n"

        // since an entry was appended, make a note of it
        entries_appended++
    }

    // if no ip addresses present, instead append a line about there being
    // no data for today.
    if entries_appended == 0 {
        whois_strings += "No whois entries given at this time."
    }

    // everything worked fine, so return the completed string contents
    return whois_strings, nil
}

//! Attempt to execute the whois command.
/*
 *  @param    ...string    list of arguments
 *
 *  @return   bytes[]      array of byte buffer data
 */
func runWhoisCommand(args ...string) (bytes.Buffer, error) {

    // variable declaration
    var output bytes.Buffer

    // assemble the command from the list of string arguments
    cmd := exec.Command("whois", args...)
    cmd.Stdout = &output
    cmd.Stderr = &output

    // attempt to execute the command
    err := cmd.Run()

    // if an error occurred, go ahead and pass it back
    if err != nil {
        return output, err
    }

    // having ran the command, pass back the result if no error has
    // occurred
    return output, nil
}