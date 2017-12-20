/*
 * Copyright (c) 2017 AlexRuzin (stan.ruzin@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
    "io"
    "flag"
    "github.com/AlexRuzin/util"
    "github.com/AlexRuzin/websock"
    "time"
    "strconv"
)

const DEFAULT_GATE = "/gate.php"

var ClientInbound = make(chan *websock.NetInstance)
var mainClient *websock.NetInstance = nil

func main() {
    util.DebugOut("[Welcome to sockchat, using the websock API]")
    mode := flag.String("mode", "server", "[client|server]")
    listenPort := flag.Int("listen-port", 7777, "Listener port number")
    connectPort := flag.Int("connect-port", 7777, "Target connection port")
    connectIP := flag.String("connect-ip", "127.0.0.1", "Listener IP")

    flag.Parse()

    if *mode == "" {
        util.DebugOut("Check usage with -h")
    }

    var err error
    switch *mode {
    case "client":
        err = clientMode(*connectIP, int16(*connectPort))
    case "server":
        err = serverMode(int16(*listenPort))
    default:
        util.DebugOut("'mode' switch must be either 'client' or 'server'")
        return
    }

    if err != nil {
        util.DebugOut(err.Error())
        return
    }
}

func serverMode(listenPort int16) error {
    util.DebugOut("Starting server mode on port: " + strconv.Itoa(int(listenPort)))
    _, err := websock.CreateServer(DEFAULT_GATE, listenPort, websock.FLAG_DEBUG, incomingClientHandler)
    if err != nil {
        return err
    }

    mainClient := <- ClientInbound

    go func () {
        for {
            util.SleepSeconds(time.Duration(util.RandInt(1,10)))
            read_data := "Controller sends: " + util.RandomString(util.RandInt(10,50))
            util.DebugOut(read_data)
            if _, err := mainClient.Write([]byte(read_data)); err != io.EOF {
                util.DebugOut(err.Error())
                panic(err)
            }
        }
    } ()

    //util.WaitForever()

    /* Read listener -- from clients */
    go func () {
        for {
            util.Sleep(10 * time.Millisecond)
            if mainClient.Len() > 0 {
                data := make([]byte, mainClient.Len())
                _, err := mainClient.Read(data)
                if err != io.EOF {
                    panic(err.Error())
                }
                util.DebugOut(string(data))
            }
        }
    } ()

    util.WaitForever()
    return nil
}

func incomingClientHandler(client *websock.NetInstance, server *websock.NetChannelService) error {
    util.DebugOut("[+] Incoming client...")

    ClientInbound <- client

    return nil
}

func clientMode(targetIP string, targetPort int16) error {
    util.DebugOut("Connecting to controller on: " + targetIP + ":" + strconv.Itoa(int(targetPort)))

    var gateURI string
    if targetPort == 80 {
        gateURI = "http://" + targetIP + DEFAULT_GATE
    } else {
        gateURI = "http://" + targetIP + ":" + strconv.Itoa(int(targetPort)) + DEFAULT_GATE
    }

    util.DebugOut("gate URI: " + gateURI)

    client, err := websock.BuildChannel(gateURI, websock.FLAG_DEBUG)
    if err != nil {
        return err
    }

    if err := client.InitializeCircuit(); err != nil {
        return err
    }

    util.DebugOut("Connected to server, beginning i/o")
    util.SleepSeconds(1)

    /* Write to stdout -- (read from socket) */
    go func () {
        for {
            util.Sleep(10 * time.Millisecond)
            if client.Len() > 0 {
                data := make([]byte, client.Len())
                _, err := client.Read(data)
                if err != io.EOF {
                    panic(err.Error())
                }
                util.DebugOut(string(data))
            }
        }
    } ()

    //util.WaitForever()
    util.SleepSeconds(1)

    /* Read user input (write to socket) */
    go func () {
        for {
            util.SleepSeconds(time.Duration(util.RandInt(1,10)))
            read_data := "Client sends: " + util.RandomString(util.RandInt(10,50))
            util.DebugOut(read_data)
            wrote, err := client.Write([]byte(read_data))
            if err != io.EOF || wrote != len(read_data) {
                panic(err.Error())
            }
        }
    } ()

    util.WaitForever()
    return nil
}













