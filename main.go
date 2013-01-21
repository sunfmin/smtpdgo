// 		Copyright 2010 Gary Sims. All rights reserved.
// 		http://www.garysims.co.uk
//
//    	This file is part of GoESMTP.
//		http://code.google.com/p/goesmtp/
//		http://goesmtp.posterous.com/
//
//    	GoESMTP is free software: you can redistribute it and/or modify
//    	it under the terms of the GNU General Public License as published by
//    	the Free Software Foundation, either version 2 of the License, or
//    	(at your option) any later version.
//
//    	GoESMTP is distributed in the hope that it will be useful,
//   	but WITHOUT ANY WARRANTY; without even the implied warranty of
//   	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    	GNU General Public License for more details.
//
//    	You should have received a copy of the GNU General Public License
//    	along with GoESMTP.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

const MAXRCPT = 100

var G_greetDomain string = ""

type SMTPStruct struct {
}

func NewSMTP() (mySMTP *SMTPStruct) {
	// Create and return a new instance of POP3Struct
	mySMTP = new(SMTPStruct)

	return
}

func (mySMTP *SMTPStruct) handleConnection(con *net.TCPConn, workerid int) {
	log.Print("IN Connection")
	// var rcpts [MAXRCPT]string
	var helodomain string = ""
	// var ehlodomain string = ""
	msgsForThisConnection := 0

	quitCmd, _ := regexp.Compile("^quit")
	heloCmd, _ := regexp.Compile("^helo ")
	ehloCmd, _ := regexp.Compile("^ehlo ")
	mailfromCmd, _ := regexp.Compile("^mail from:")
	rcpttoCmd, _ := regexp.Compile("^rcpt to:")
	dataCmd, _ := regexp.Compile("^data")
	rsetCmd, _ := regexp.Compile("^rset")
	noopCmd, _ := regexp.Compile("^noop")
	vrfyCmd, _ := regexp.Compile("^vrfy")
	authCmd, _ := regexp.Compile("^auth ")

	disconnected := false
	// authenticated := false

	// Check the remote server... Log IP and find out if SPAM server
	raddr := con.RemoteAddr()
	ipbits := strings.Split(raddr.String(), ":")
	log.Printf("SMTP connection from %s", ipbits[0])

	// Send greeting
	welcome := fmt.Sprintf("220 %s ESMTP\r\n", G_greetDomain)
	con.Write([]byte(welcome))

	buf := bufio.NewReader(con)
	for {
		lineofbytes, err := buf.ReadBytes('\n')
		log.Println(string(lineofbytes))
		if err != nil {
			con.Close()
			disconnected = true
			break
		} else {
			lineofbytes = TrimCRLF(lineofbytes)
			lineofbytesL := bytes.ToLower(lineofbytes)
			log.Printf(string(lineofbytes))

			if len(lineofbytes) > 0 {
				switch {
				case quitCmd.Match(lineofbytesL):
					con.Write([]byte("221 Bye for now\r\n"))
					con.Close()
					disconnected = true
					break
				case heloCmd.Match(lineofbytesL):
					helodomain = string(getDomainFromHELO(lineofbytes))
					helor := fmt.Sprintf("250 %s nice to meet you.\r\n", helodomain)
					con.Write([]byte(helor))
					// HELO / EHLO is also like a RSET
					break
				case ehloCmd.Match(lineofbytesL):
					// ehlodomain = string(getDomainFromHELO(lineofbytes))
					ehlor := fmt.Sprintf("250-%s\r\n", G_greetDomain)
					con.Write([]byte(ehlor))
					con.Write([]byte("250-AUTH=PLAIN CRAM-MD5\r\n"))
					con.Write([]byte("250 AUTH PLAIN CRAM-MD5\r\n"))
					// HELO / EHLO is also like a RSET
					break
				case mailfromCmd.Match(lineofbytesL):
					con.Write([]byte("250 OK\r\n"))
					// if (len(helodomain) > 0) || (len(ehlodomain) > 0) {
					// 	// Beginning of a new mail transaction
					// 	rcptsi = 0
					// 	// Check if spam server
					// 	if authenticated == false {
					// 		con.Close()
					// 		disconnected = true
					// 		break
					// 	}
					// mailfrom = string(lineofbytesL)
					// 	// mailfromaddr := string(getAddressFromMailFrom(lineofbytes))
					// 	// MAIL FROM:<> is valid
					// 	// Known as NULL return path; see RFC 2821 section 6.1
					// 	// if (len(mailfromaddr) == 0) || ((strings.Index(mailfromaddr, "@") != -1) && (strings.Index(mailfromaddr, ".") != -1)) {
					// 	con.Write([]byte("250 OK\r\n"))
					// 	// } else {
					// 	// 	con.Write([]byte("550 Bad email address\r\n"))
					// 	// }
					// 	msgsForThisConnection++
					// } else {
					// 	con.Write([]byte("503 Bad sequence of commands\r\n"))
					// }
					break
				case rcpttoCmd.Match(lineofbytesL):
					con.Write([]byte("250 OK\r\n"))
					// if (len(helodomain) > 0) || (len(ehlodomain) > 0) {
					// 	rcpts[rcptsi] = string(lineofbytesL)
					// 	rcpttoaddr := string(getAddressFromRcptTo(lineofbytesL))

					// 	if (strings.Index(rcpttoaddr, "@") != -1) && (strings.Index(rcpttoaddr, ".") != -1) {
					// 		if authenticated == false {
					// 			// If not authenticated only accept local mailbox recipients
					// 			ourpassword, _ := "", ""
					// 			if len(ourpassword) > 0 {
					// 				// OK, local mailbox
					// 				con.Write([]byte("250 OK\r\n"))
					// 				rcptsi += 1
					// 			} else {
					// 				con.Write([]byte("550 No such user here\r\n"))
					// 			}
					// 		} else {
					// 			con.Write([]byte("250 OK\r\n"))
					// 			rcptsi += 1
					// 		}
					// 	} else {
					// 		con.Write([]byte("550 Bad email address\r\n"))
					// 	}
					// } else {
					// 	con.Write([]byte("503 Bad sequence of commands\r\n"))
					// }
					break
				case dataCmd.Match(lineofbytesL):
					log.Println("DATA CMD")
					con.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))
					mySMTP.recvBodyToFile(con)
					con.Write([]byte("250 OK\r\n"))

					// if (len(mailfrom) > 0) && (rcptsi > 0) {
					// 	con.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))

					// 	msgFilename := "abc" //getFilenameForMsg(workerid, msgsForThisConnection)
					// 	// mySMTP.squirtHeaderToFile(helodomain, ehlodomain, mailfrom, rcpts, rcptsi, msgFilename)

					// 	if true {
					// 		log.Printf("New message received - %s (%s)\n", mailfrom, msgFilename)
					// 		con.Write([]byte("250 OK\r\n"))
					// 	} else {
					// 		con.Write([]byte("554 Transaction failed\r\n"))
					// 	}
					// } else {
					// 	con.Write([]byte("503 Bad sequence of commands\r\n"))
					// }
					break
				case noopCmd.Match(lineofbytesL):
					con.Write([]byte("250 OK\r\n"))
					break
				case rsetCmd.Match(lineofbytesL):
					msgsForThisConnection++
					con.Write([]byte("250 OK\r\n"))
					break
				case vrfyCmd.Match(lineofbytesL):
					con.Write([]byte("Cannot VRFY user\r\n"))
					break
				case authCmd.Match(lineofbytesL):
					con.Write([]byte("235 Authentication successful.\r\n"))

					// f := strings.Split(string(lineofbytes), " ")
					// if len(f) < 2 {
					// 	con.Write([]byte("504 Unrecognized authentication.\r\n"))
					// } else {
					// 	authtype := strings.ToLower(f[1])
					// 	if authtype == "plain" {
					// 		// The client can either give the authenication string here or in
					// 		// the next command
					// 		if len(f) == 2 {
					// 			con.Write([]byte("334 \r\n"))
					// 			lineofbytes, err := buf.ReadBytes('\n')
					// 			if err != nil {
					// 				con.Close()
					// 				disconnected = true
					// 				break
					// 			} else {
					// 				lineofbytes = TrimCRLF(lineofbytes)
					// 				_, u1, p := decodeSMTPAuthPlain(string(lineofbytes))
					// 				ourpassword, _ := "", ""
					// 				if (len(ourpassword) > 0) && (ourpassword == p) {
					// 					con.Write([]byte("235 Authentication successful.\r\n"))
					// 					authenticated = true
					// 				} else {
					// 					con.Write([]byte("535 Authentication failed.\r\n"))
					// 					authenticated = false
					// 				}
					// 			}
					// 		} else {
					// 			_, u1, p := decodeSMTPAuthPlain(f[2])
					// 			ourpassword, _ := "", ""
					// 			if (len(ourpassword) > 0) && (ourpassword == p) {
					// 				con.Write([]byte("235 Authentication successful.\r\n"))
					// 				authenticated = true
					// 			} else {
					// 				con.Write([]byte("535 Authentication failed.\r\n"))
					// 				authenticated = false
					// 			}
					// 		}
					// 	} else if authtype == "cram-md5" {
					// 		cram1 := fmt.Sprintf("<%d.%d@%s>", os.Getpid(), time.Now().Second(), "localhost")
					// 		cram2 := fmt.Sprintf("334 %s\r\n", encodeBase64String(cram1))
					// 		con.Write([]byte(cram2))
					// 		lineofbytes, err := buf.ReadBytes('\n')
					// 		if err != nil {
					// 			con.Close()
					// 			disconnected = true
					// 			break
					// 		} else {
					// 			lineofbytes = TrimCRLF(lineofbytes)
					// 			cram3 := decodeBase64String(string(lineofbytes))
					// 			cram3 = string(TrimCRLF([]byte(cram3)))
					// 			f := strings.Split(cram3, " ")
					// 			if len(f) != 2 {
					// 				con.Write([]byte("535 Authentication failed.\r\n"))
					// 				authenticated = false
					// 			} else {
					// 				ourpassword, _ := "", ""
					// 				if (len(ourpassword) > 0) && (keyedMD5(ourpassword, cram1) == f[1]) {
					// 					con.Write([]byte("235 Authentication successful.\r\n"))
					// 					authenticated = true
					// 				} else {
					// 					con.Write([]byte("535 Authentication failed.\r\n"))
					// 					authenticated = false
					// 				}
					// 			}
					// 		}
					// 	} else {
					// 		con.Write([]byte("504 Unrecognized authentication type.\r\n"))
					// 	}
					// }
					break
				default:
					con.Write([]byte("502 unimplemented\r\n"))
					break
				}
			}
		}

		if disconnected == true {
			break
		}
	}
}

func (mySMTP *SMTPStruct) recvBodyToFile(con *net.TCPConn) bool {
	disconnected := false
	sts := true

	buf := bufio.NewReader(con)
	for {
		log.Println("looping")
		lineofbytes, err := buf.ReadBytes('\n')
		if err != nil {
			con.Close()
			disconnected = true
			sts = false
			break
		} else {
			if (len(lineofbytes) == 3) && (lineofbytes[0] == '.') && (lineofbytes[1] == '\r') && (lineofbytes[2] == '\n') {
				disconnected = true
			} else {
				log.Print(string(lineofbytes))
			}
		}

		if disconnected == true {
			break
		}
	}
	return sts
}

func (mySMTP *SMTPStruct) startSMTP() {
	workerid := 0

	for {

		listener, err := net.Listen("tcp", ":8009")
		if err != nil {
			panic(err)
		}
		for {
			con, _ := listener.Accept()

			go mySMTP.handleConnection(con.(*net.TCPConn), workerid)
			workerid += 1

		}

		listener.Close()
	}
}

func main() {
	smtp := &SMTPStruct{}
	smtp.startSMTP()
}
