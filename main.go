package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
_	"strings"
	"os/user"
	"log"
)

const (
	PW_MAGIC = 0xA3
	PW_FLAG  = 0xFF
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 && len(args) != 2 {
		fmt.Println("WinSCP stored password finder\n")
		fmt.Println("Registry:")
		fmt.Println("  Open regedit and navigate to [HKEY_CURRENT_USER\\Software\\Martin Prikryl\\WinSCP 2\\Sessions] to get the hostname, username and encrypted password\n")
		if runtime.GOOS == "windows" {
			fmt.Println("  Usage winscppasswd.exe <host> <username> <encrypted_password>")
		} else {
			fmt.Println("  Usage ./winscppasswd <host> <username> <encrypted_password>")
		}
		fmt.Println("\n\nWinSCP.ini:")
		if runtime.GOOS == "windows" {
			fmt.Println("  Usage winscppasswd.exe ini [<filepath>]")
		} else {
			fmt.Println("  Usage ./winscppasswd ini [<filepath>]")
		}
		fmt.Printf("  Default value <filepath>: %s\n", defaultWinSCPIniFilePath());
		return
	}
	fmt.Println(decrypt(args[0], args[1], args[2]))
}

func defaultWinSCPIniFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	return usr.HomeDir + "\\AppData\\Roaming\\winSCP.ini"
}

func decrypt(host, username, password string) string {
	key := username + host
	passbytes := []byte{}
	for i := 0; i < len(password); i++ {
		val, _ := strconv.ParseInt(string(password[i]), 16, 8)
		passbytes = append(passbytes, byte(val))
	}
	var flag byte
	flag, passbytes = dec_next_char(passbytes)
	var length byte = 0
	if flag == PW_FLAG {
		_, passbytes = dec_next_char(passbytes)

		length, passbytes = dec_next_char(passbytes)
	} else {
		length = flag
	}
	toBeDeleted, passbytes := dec_next_char(passbytes)
	passbytes = passbytes[toBeDeleted*2:]

	clearpass := ""
	var (
		i   byte
		val byte
	)
	for i = 0; i < length; i++ {
		val, passbytes = dec_next_char(passbytes)
		clearpass += string(val)
	}

	if flag == PW_FLAG {
		clearpass = clearpass[len(key):]
	}
	return clearpass
}

func dec_next_char(passbytes []byte) (byte, []byte) {
	if len(passbytes) <= 0 {
		return 0, passbytes
	}
	a := passbytes[0]
	b := passbytes[1]
	passbytes = passbytes[2:]
	return ^(((a << 4) + b) ^ PW_MAGIC) & 0xff, passbytes
}
