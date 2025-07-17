package poolStratum

import (
	"encoding/json"
	"fmt"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/snipa22/go-tari-p2pool-interface/subsystems/config"
	"github.com/snipa22/go-tari-p2pool-interface/subsystems/messages"
	"io"
	"net"
	"sync"
	"strings"
	"time"
)

const httpDummy = ` 200 OK
Content-Type: text/plain
Content-Length: 18

Mining Pool Online`

func ClientConn(hashChan chan messages.HashToVerify, milieu *core.Milieu) func(net.Conn) {
	return func(inConn net.Conn) {
		// The core loop for a connected client is as follows
		// Register the client has connected
		// Init the miner data structure
		// Listen for traffic on the connection, prepping the natural timeout system
		var closeOnce sync.Once
		quit := make(chan bool, 1)
		miner := newMiner(inConn, hashChan, milieu, quit)
		newBt := make(chan bool)

		// TODO: Add tracking for connection counts for miner counts

		// cleanup logic — only runs once
		cleanup := func() {
			closeOnce.Do(func() {
				// TODO: Add tracking for disconnection counts for miner counts
				fmt.Println("Closing miner")
				for _, v := range miner.cronJobs {
					config.SystemCrons.DeleteCronJob(v)
				}
				_ = inConn.Close()
				miner = nil
			})
		}

		time.AfterFunc(30*time.Second, func() {
			if miner != nil && !miner.Active {
				select {
				case quit <- true:
				default:
				}
			}
		})
		time.AfterFunc(45*time.Second, func() {
			if miner != nil && miner.Address == "" {
				select {
				case quit <- true:
				default:
				}
			}
		})

		tmp := make([]byte, 2048)
		data := make([]byte, 0)

		for {
			select {
			case <-newBt:
				// TODO: Call the new job subroutine
			case <-quit:
				fmt.Println("ClientConn quit")
				miner.Active = false
				cleanup()
				return
			default:
				// Perform data read and traffic management
				n, err := inConn.Read(tmp)
				if err != nil {
					// log if not normal error
					if err != io.EOF {
						if strings.Contains(err.Error(), "connection timed out") {
							milieu.Error("Exiting due to a connection timeout")
							cleanup()
							return
						}
						if strings.Contains(err.Error(), "connection reset by peer") {
							milieu.Error("Exiting due to a connection reset by peer")
							cleanup()
							return
						}
						milieu.Error(err.Error())
					}
					cleanup()
					return
				}
				found := 0
				data = append(data, tmp[:n]...)

				for k, v := range data {
					if v == '\n' {
						found = k
					}
				}
				if found == 0 {
					if len(data) >= 20480 {
						cleanup()
						return
					}
				}

				// TODO: Add IP bans
				//if config.GlobalConfig.IsIPBanned(miner.IpAddress) {
				//	miner.RPCResponse("IP Address currently banned for using an invalid mining protocol, please check your miner", "")
				//	return
				//}

				rawJson := string(data[:found])
				data = data[found:]
				for _, v := range strings.Split(rawJson, "\n") {
					// Close any HTTP only request
					if strings.HasPrefix(v, "GET /") {
						if strings.Contains(v, "HTTP/1.1") {
							_, _ = inConn.Write([]byte(fmt.Sprintf("HTTP/1.1%v", httpDummy)))
						} else if strings.Contains(v, "HTTP/1.0") {
							_, _ = inConn.Write([]byte(fmt.Sprintf("HTTP/1.0%v", httpDummy)))
						}
						cleanup()
						return
					}
					// Parse and handle the real work
					if len(v) < 10 {
						continue
					}
					parsedJson := messages.MinerRPCRequest{}
					if err = json.Unmarshal([]byte(v), &parsedJson); err != nil {
						milieu.Error(err.Error())
						cleanup()
						return
					}
					miner.Active = true
					miner.UpdateLastComm()
					miner.rpcID = parsedJson.ID
					// Miner is now set active, parse and handle various cases
					switch parsedJson.Method {
					case `login`:
						// Do login stuff
						miner.Login(parsedJson.Params)
						break
					case `getjob`:
						// Do getjob stuff
						miner.SendNewJob(false)
						break
					case `submit`:
						// Do share submission stuff
						// TODO: Build the share submission system
						miner.SubmitJob(parsedJson.Params)
						break
					case `keepalived`:
						// Do keepalive
						if miner.Address == "" {
							miner.RPCResponse("Unauthenticated", "")
							cleanup()
							return
						}
						miner.RPCResponse("", "KEEPALIVED")
						break
					default:
						// Invalid message, log it, but do not kill conn
						milieu.Info(fmt.Sprintf("Received invalid message: %v", v))
						break
					}
				}
			}
		}
	}
}
