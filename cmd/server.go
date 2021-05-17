package main

import (
	_ "bytes"
	"log"
	"net"

	"github.com/gobwas/ws"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:1100")
	if err != nil {
		// handle error
	}
	e := wsflate.Extension{
		// We are using default parameters here since we use
		// wsflate.{Compress,Decompress}Frame helpers below in the code.
		// This assumes that we use standard compress/flate package as flate
		// implementation.
		Parameters: wsflate.DefaultParameters,
	}
	u := ws.Upgrader{
		Negotiate: e.Negotiate,
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Reset extension after previous upgrades.
		e.Reset()

		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
			continue
		}
		if _, ok := e.Accepted(); !ok {
			log.Printf("didn't negotiate compression for %s", conn.RemoteAddr())
			conn.Close()
			continue
		}

		go func() {
			defer conn.Close()
			for {
				frame, err := ws.ReadFrame(conn)
				if err != nil {
					// Handle error.
					return
				}

				frame = ws.UnmaskFrameInPlace(frame)

				if wsflate.IsCompressed(frame.Header) {
					// Note that even after successful negotiation of
					// compression extension, both sides are able to send
					// non-compressed messages.
					frame, err = wsflate.DecompressFrame(frame)
					if err != nil {
						// Handle error.
						return
					}
				}

				// Do something with frame...

				ack := ws.NewTextFrame([]byte("this is an acknowledgement"))

				// Compress response unconditionally.
				ack, err = wsflate.CompressFrame(ack)
				if err != nil {
					// Handle error.
					return
				}
				if err = ws.WriteFrame(conn, ack); err != nil {
					// Handle error.
					return
				}
			}
		}()
	}
}