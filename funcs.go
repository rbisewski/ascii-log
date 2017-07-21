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
    "regexp"
    "strings"
    "sort"
    "strconv"
)

//! Validate an IPv6 address
/*
 * @param     string    /path/to/file
 *
 * @return    bool      whether or not this is true
 *
 * TODO: add more logic to this function
 */
func isValidIPv6Address(ip string) (bool) {

    // input validation
    if len(ip) < 1 {
        return false
    }

    // attempt to split the string into pieces via the ':' char
    ip_pieces := strings.Split(ip, ":")

    // safety check, ensure there is at least one piece
    if len(ip_pieces) < 1 {
        return false
    }

    // for every hexadecimal piece of the IPv6 address...
    for _, hexa := range ip_pieces {

        // ensure it has a length between 1 and 4
        if len(hexa) < 1 || len(hexa) > 4 {
            return false
        }

        // convert the ip_piece string to an integer
        hexa_as_uint, err := strconv.ParseUint(hexa, 0, 16)

        // if an error occurs, go ahead and return false
        if err != nil {
            return false
        }

        // if greater than 0xFFFF pass back a false
        if hexa_as_uint > 65535 {
            return false
        }
    }

    // if all the tests passed, go ahead and return true
    return true
}

//! Validate an IPv4 address
/*
 * @param     string    /path/to/file
 *
 * @return    bool      whether or not this is true
 */
func isValidIPv4Address(ip string) (bool) {

    // input validation
    if len(ip) < 1 {
        return false
    }

    // ensure that the ip address is valid length
    //
    // 0.0.0.0 --> 8 chars (min)
    //
    // 127.123.123.123 --> 15 chars (max)
    //
    if len(ip) < 8 || len(ip) > 15 {
        return false
    }

    // attempt to split the string into pieces via the '.' char
    ip_pieces := strings.Split(ip, ".")

    // ensure that there are at least 4 pieces
    if len(ip_pieces) != 4 {
        return false
    }

    // for every oct piece of the IPv4 address...
    for _, oct := range ip_pieces {

        // ensure it has a length of at least 1
        if len(oct) < 1 {
            return false
        }

        // convert the ip_piece string to an integer
        oct_as_uint, err := strconv.ParseUint(oct, 0, 10)

        // if an error occurred, throw back a false
        if err != nil {
            return false
        }

        // ensure that the integer is between 0 and 255; actually it is a
        // unsigned int at this point, so only need check if > 255
        if oct_as_uint > 255 {
            return false
        }
    }

    // otherwise it appears to be a proper IPv4
    return true
}

//! Take a given IP address and space buffer it so that it is always 15
//! characters long.
/*
 * @param    string    IPv4 address
 *
 * @param    string    space-formatted IPv4 address
 * @param    error     error message, if any
 */
func spaceFormatIPv4(ip string) (string, error) {

    // input validation
    if len(ip) < 1 || len(ip) > 15 {
        return "", fmt.Errorf("spaceFormatIPv4() --> invalid input\n");
    }

    // ensure this is actually a IPv4 address
    if !isValidIPv4Address(ip) {
        return "", fmt.Errorf("spaceFormatIPv4() --> given IP is not " +
          "an IPv4 address\n")
    }

    // attempt to format the IPv4 address
    space_formatted_ip_address := ip
    for len(space_formatted_ip_address) < 16 {
        space_formatted_ip_address += " "
    }

    // return the formatted IPv4 string
    return space_formatted_ip_address, nil
}

//! Convert a given IPv4 address to a x.x.x.0/24 CIDR notation
/*
 * @param    string    an IPv4 address
 *
 * @return   string    result as a /24
 * @return   error     error message, if any
 */
func obtainSlash24FromIpv4(ip string) (string, error) {

    // input validation
    if len(ip) < 1 {
        return "", fmt.Errorf("obtainSlash24FromIpv4() --> invalid input")
    }

    // ensure the given value is actually an IP4 address
    if !isValidIPv4Address(ip) {
        return "", fmt.Errorf("obtainSlash24FromIpv4() --> improper " +
          "IPv4 address given")
    }

    // variable declaration
    ipv4_slash24_cidr := ""

    // separate the IPv4 address string into pieces
    ip_pieces := strings.Split(ip, ".")

    // ensure that there are at least 4 pieces
    if len(ip_pieces) != 4 {
        return "", fmt.Errorf("obtainSlash24FromIpv4() --> non-standard " +
          "IPv4 address")
    }

    // reconstruct the IPv4 address string
    ipv4_slash24_cidr += ip_pieces[0]
    ipv4_slash24_cidr += "."
    ipv4_slash24_cidr += ip_pieces[1]
    ipv4_slash24_cidr += "."
    ipv4_slash24_cidr += ip_pieces[2]
    ipv4_slash24_cidr += "."
    ipv4_slash24_cidr += "0/24"

    // having gone this far, return the adjusted result
    return ipv4_slash24_cidr, nil
}

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
 * @param     map        string map containing ip/whois country data
 *
 * @return    string     lines that contain "count | ip | country | host \n"
 *            error      error message, if any
 */
func convertIpAddressMapToString(ip_map map[string] int,
  whois_country_map map[string] string) (string, error) {

    // input validation
    if len(ip_map) < 1 || len(whois_country_map) < 1 {
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

        // lookup the country code
        country_code := whois_country_map[ip]

        // safety check, fallback to "--" if the country code is blank or
        // nil or unusual length
        if len(country_code) != 2 {
            country_code = "--"
        }

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
        space_formatted_ip_address, err := spaceFormatIPv4(ip)

        // if an error occurs, skip to the next element
        if err != nil {
           continue
        }

        // append that count | address |  hostname
        ip_strings += strconv.Itoa(count) + "\t"
        ip_strings += " | "
        ip_strings += space_formatted_ip_address
        ip_strings += " | "
        ip_strings += country_code
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
 * @param     map       string map containing ip addresses and counts
 *
 * @return    string    whois data of every given ip
 * @return    map       string map containing whois country data
 * @return    error     error message, if any
 *
 *
 * TODO: this function could use more testing
 */
func obtainWhoisEntries(ip_map map[string] int) (string, map[string] string,
  error) {

    // input validation
    if len(ip_map) < 1 {
        return "", nil, fmt.Errorf("obtainWhoisEntries() --> invalid input")
    }

    // variable declaration
    var whois_strings string  = ""
    var whois_summary_map     = make(map[string] string)
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

        // compile a regex that looks for "country: XX\n" or "Country: XX\n"
        re := regexp.MustCompile("[cC]ountry:[^\n]{2,32}\n")

        //
        // TODO: adjust this so it obtains all of the country data present
        //       in the whois entry.
        //
        // attempt to obtain the country of a given IP address
        whois_regex_country_result := re.FindString(trimmed_string)

        //
        // TODO: implement this pseudo code
        //

        // for every whois country line...

            // if there is more than one entry, take the last one since
            // the others are likely ARIN/RIPE/etc data and therefore not
            // quite as useful as the actual origin country network.

        // trim the result
        whois_regex_country_result =
          strings.Trim(whois_regex_country_result, " ")
        whois_regex_country_result =
          strings.Trim(whois_regex_country_result, "\n")

        // ensure that the result still has 2 letters
        if len(whois_regex_country_result) < 2 {
            whois_regex_country_result = "--"
        }

        // split up the string using spaces
        wr_pieces := strings.Split(whois_regex_country_result, " ")

        // safety check, ensure there is at least 1 pieces
        if len(wr_pieces) < 1 {
            whois_regex_country_result = "--"
        }

        // assemble a regex to test the country code
        re_country_code := regexp.MustCompile("^[A-Za-z]{2}$")

        // search thru the pieces for the country code result
        for _, code := range wr_pieces {

            // if the code is not equal to 2
            if len(code) != 2 {
                continue
            }

            // ensure the code is actually two alphabet characters
            verify := re_country_code.FindString(code)

            // skip a line if the entry is not the latest date
            if len(verify) != 2 {
                continue
            }

            // assign the code to the whois country result
            whois_regex_country_result = code

            // leave the loop
            break
        }

        // append it to the whois map
        whois_summary_map[ip] = whois_regex_country_result

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
    return whois_strings, whois_summary_map, nil
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
