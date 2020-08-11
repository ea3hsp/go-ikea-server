package tradfri

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dustin/go-coap"
	"github.com/eriklupander/dtls"
)

// DtlsClient provides an domain-agnostic CoAP-client with DTLS transport.
type DtlsClient struct {
	logger         log.Logger
	peer           *dtls.Peer
	msgID          uint16
	gatewayAddress string
	clientID       string
	psk            string
}

// NewDtlsClient acts as factory function, returns a pointer to a connected (or will panic) DtlsClient.
func NewDtlsClient(gatewayAddress, clientID, psk string, logger log.Logger) *DtlsClient {
	client := &DtlsClient{
		logger:         logger,
		gatewayAddress: gatewayAddress,
		clientID:       clientID,
		psk:            psk,
	}
	client.connect()
	return client
}

func (dc *DtlsClient) connect() {
	//
	dc.setupKeystore()
	//
	listener, err := dtls.NewUdpListener(":0", time.Second*900)
	if err != nil {
		panic(err.Error())
	}
	//
	peerParams := &dtls.PeerParams{
		Addr:             dc.gatewayAddress,
		Identity:         dc.clientID,
		HandshakeTimeout: time.Second * 15,
	}

	dc.logger.Log("[info]", "connecting to ", "peer", dc.gatewayAddress)

	dc.peer, err = listener.AddPeerWithParams(peerParams)
	if err != nil {
		dc.logger.Log("[error]", "unable to connect to gateway", "gateway", dc.gatewayAddress, "msg", err.Error())
		os.Exit(1)
	}
	dc.peer.UseQueue(true)
	dc.logger.Log("[info]", "DTLS connection established", "gateway", dc.gatewayAddress)
}

// Call writes the supplied coap.Message to the peer
func (dc *DtlsClient) Call(req coap.Message) (coap.Message, error) {
	dc.logger.Log("[info]", "calling", "code", req.Code.String(), "path", req.PathString())
	data, err := req.MarshalBinary()
	if err != nil {
		return coap.Message{}, err
	}
	err = dc.peer.Write(data)

	if err != nil {
		return coap.Message{}, err
	}

	respData, err := dc.peer.Read(time.Second)
	if err != nil {
		return coap.Message{}, err
	}

	msg, err := coap.ParseMessage(respData)
	if err != nil {
		return coap.Message{}, err
	}

	/* 	logrus.Info("Response: ")
	   	dc.logger.Log("MessageID: %v\n", msg.MessageID)
	   	dc.logger.Log("Type: %v\n", msg.Type)
	   	dc.logger.Log("Code: %v\n", msg.Code)
	   	dc.logger.Log("Token: %v\n", msg.Token)
	   	dc.logger.Log("Payload: %v\n", string(msg.Payload)) */

	return msg, nil
}

// BuildGETMessage produces a CoAP GET message with the next msgID set.
func (dc *DtlsClient) BuildGETMessage(path string) coap.Message {
	dc.msgID++
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: dc.msgID,
	}
	req.SetPathString(path)
	return req
}

// BuildPUTMessage produces a CoAP PUT message with the next msgID set.
func (dc *DtlsClient) BuildPUTMessage(path string, payload string) coap.Message {
	dc.msgID++

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.PUT,
		MessageID: dc.msgID,
		Payload:   []byte(payload),
	}
	req.SetPathString(path)

	return req
}

// BuildPOSTMessage produces a CoAP POST message with the next msgID set.
func (dc *DtlsClient) BuildPOSTMessage(path string, payload string) coap.Message {
	dc.msgID++

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: dc.msgID,
		Payload:   []byte(payload),
	}
	req.SetPathString(path)

	return req
}

func (dc *DtlsClient) setupKeystore() {
	mks := dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{mks})
	mks.AddKey(dc.clientID, []byte(dc.psk))
}
