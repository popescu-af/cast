package chromecast

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"

	"chromecast/pb"
	"logger"
)

type transport struct {
	Conn net.Conn
}

func newTransport(
	ctx context.Context,
	ip net.IP,
	port int,
) (*transport, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	deadline, _ := ctx.Deadline()
	dialer := &net.Dialer{
		Deadline: deadline,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	return &transport{
		Conn: conn,
	}, nil
}

func (t *transport) close() error {
	return t.Conn.Close()
}

type messageCallback func(*pkg) error

func (t *transport) startReceiving(messageCallback messageCallback) {
	go t.receiveProc(messageCallback)
}

type pkg struct {
	source      string
	destination string
	namespace   string
	payload     []byte
}

func (t *transport) send(p *pkg) error {
	payloadString := string(p.payload)
	message := &pb.CastMessage{
		ProtocolVersion: pb.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &p.source,
		DestinationId:   &p.destination,
		Namespace:       &p.namespace,
		PayloadType:     pb.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}
	proto.SetDefaults(message)

	logger.Log.Printf("send(src: %s, dst: %s, ns: %s): %s\n", p.source, p.destination, p.namespace, payloadString)

	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %s", err)
	}

	err = binary.Write(t.Conn, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return fmt.Errorf("failed to write length: %s", err)
	}
	_, err = t.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %s", err)
	}
	return nil
}

func (t *transport) receiveProc(messageCallback messageCallback) {
	for {
		var length uint32
		err := binary.Read(t.Conn, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				// logger.Log.Printf("nothing on the pipe, sleeping for 1s...")
				time.Sleep(time.Second)
				continue
			}
			logger.Log.Printf("failed to read packet length: %s", err)
			return
		}
		if length == 0 {
			logger.Log.Printf("empty packet")
			continue
		}

		packet := make([]byte, length)
		_, err = io.ReadFull(t.Conn, packet)
		if err != nil {
			logger.Log.Printf("failed to read full packet: %s", err)
			return
		}

		m := new(pb.CastMessage)
		err = proto.Unmarshal(packet, m)
		if err != nil {
			logger.Log.Printf("failed to unmarshal packet: %s", err)
			return
		}

		logger.Log.Printf("recv(src: %s, dst: %s, ns: %s): %s\n", *m.SourceId, *m.DestinationId, *m.Namespace, *m.PayloadUtf8)

		p := &pkg{
			source:      *m.SourceId,
			destination: *m.DestinationId,
			namespace:   *m.Namespace,
			payload:     []byte(*m.PayloadUtf8),
		}
		if err := messageCallback(p); err != nil {
			logger.Log.Printf("failed to handle message: %s\n", err)
		}
	}
}
