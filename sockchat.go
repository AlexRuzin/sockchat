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
    "flag"
    "github.com/AlexRuzin/util"
    "github.com/AlexRuzin/websock"
)

const DEFAULT_GATE = "/gate.php"

func main() {
    util.DebugOut("[Welcome to sockchat, using the websock API]")
    mode := flag.String("mode", "", "[client|server]")
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
    util.DebugOut("Starting server mode on port: " + string(listenPort))
    _, err := websock.CreateServer(DEFAULT_GATE, listenPort, websock.FLAG_DEBUG, incomingClientHandler)
    if err != nil {
        return err
    }

    util.WaitForever()
    return nil
}

func incomingClientHandler(client *websock.NetInstance, server *websock.NetChannelService) error {
    util.DebugOut("Incoming client...")
    return nil
}

func clientMode(targetIP string, targetPort int16) error {
    util.DebugOut("Connecting to controller on: " + targetIP + ":" + string(targetPort))

    var gateURI = "http://" + targetIP + DEFAULT_GATE
    util.DebugOut("gate URI: " + gateURI)

    client, err := websock.BuildChannel(gateURI, targetPort, websock.FLAG_DEBUG)
    if err != nil {
        return err
    }

    if err := client.InitializeCircuit(); err != nil {
        return err
    }


    util.WaitForever()
    return nil
}













