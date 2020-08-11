package tradfri

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-coap"
	"github.com/eriklupander/dtls"
	"github.com/eriklupander/tradfri-go/dtlscoap"
	"github.com/sirupsen/logrus"

	"github.com/go-kit/kit/log"
)

const (
	token = "o9BS7ebX0gAmfBIq"
)

// Client tradfri client
type Client struct {
	logger         log.Logger
	clientID       string
	gatewayAddress string
	dtlsclient     *dtlscoap.DtlsClient
}

// NewClient creates new tradfri client
func NewClient(gatewayAddress, clientID, psk string, logger log.Logger) ITradfri {
	return &Client{
		logger:         logger,
		clientID:       clientID,
		gatewayAddress: gatewayAddress,
	}
}

// Connect gateway connection
func (c *Client) Connect() error {
	// creates listener
	listener, err := dtls.NewUdpListener(":0", time.Second*900)
	if err != nil {
		c.logger.Log("[error]", "listener", "msg", err.Error)
		os.Exit(1)
	}
	// peer params
	peerParams := &dtls.PeerParams{
		Addr:             c.gatewayAddress,
		Identity:         c.clientID,
		HandshakeTimeout: time.Second * 15}
	c.logger.Log("[info]", "connecting", "peer", c.gatewayAddress)

	peer, err = listener.AddPeerWithParams(peerParams)
	if err != nil {
		logrus.Infof("Unable to connect to Gateway at %v: %v\n", c.gatewayAddress, err.Error())
		os.Exit(1)
	}
	peer.UseQueue(true)
	logrus.Infof("DTLS connection established to %v\n", c.gatewayAddress)
}

// DevicePower puts on/off a tradfri device
func (c *Client) DevicePower(id, power int) (interface{}, error) {
	return nil, nil
}

func (c *Client) call(msg coap.Message) (coap.Message, error) {
	return c.dtlsclient.Call(msg)
}

// AuthExchange ...
func (c *Client) AuthExchange(clientID string) (interface{}, error) {
	// request CoAP message
	req := c.dtlsclient.BuildPOSTMessage("/15011/9063", fmt.Sprintf(`{"9090":"%s"}`, clientID))
	// Send CoAP message for token exchange
	resp, err := c.call(req)
	if err != nil {
		c.logger.Log("[error]", "performing gateway token exchange request", "msg", err.Error)
		return nil, err
	}
	// creates json decoder
	dec := json.NewDecoder(bytes.NewReader(resp.Payload))
	// json map holder
	var m map[string]interface{}
	dec.Decode(&m)
	if err != nil {
		c.logger.Log("error decoding json response from Gateway for token exchange")
		return nil, err
	}
	// return msg
	return m, nil
}
