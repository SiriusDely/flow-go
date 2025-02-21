// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package middleware

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	ggio "github.com/gogo/protobuf/io"
	"github.com/ipfs/go-datastore"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/network/p2p"
	"github.com/onflow/flow-go/network/p2p/blob"
	"github.com/onflow/flow-go/network/p2p/internal/p2putils"
	"github.com/onflow/flow-go/network/p2p/p2pbuilder"
	"github.com/onflow/flow-go/network/p2p/p2pnode"
	"github.com/onflow/flow-go/network/p2p/ping"
	"github.com/onflow/flow-go/network/p2p/utils"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/component"
	"github.com/onflow/flow-go/module/irrecoverable"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/channels"
	"github.com/onflow/flow-go/network/codec"
	"github.com/onflow/flow-go/network/message"
	"github.com/onflow/flow-go/network/p2p/unicast"
	"github.com/onflow/flow-go/network/slashing"
	"github.com/onflow/flow-go/network/validator"
	_ "github.com/onflow/flow-go/utils/binstat"
)

const (
	_ = iota
	_ = 1 << (10 * iota)
	mb
	gb
)

const (
	// DefaultMaxUnicastMsgSize defines maximum message size in unicast mode for most messages
	DefaultMaxUnicastMsgSize = 10 * mb // 10 mb

	// LargeMsgMaxUnicastMsgSize defines maximum message size in unicast mode for large messages
	LargeMsgMaxUnicastMsgSize = gb // 1 gb

	// DefaultUnicastTimeout is the default maximum time to wait for a default unicast request to complete
	// assuming at least a 1mb/sec connection
	DefaultUnicastTimeout = 5 * time.Second

	// LargeMsgUnicastTimeout is the maximum time to wait for a unicast request to complete for large message size
	LargeMsgUnicastTimeout = 1000 * time.Second
)

var (
	_ network.Middleware = (*Middleware)(nil)
)

// Middleware handles the input & output on the direct connections we have to
// our neighbours on the peer-to-peer network.
type Middleware struct {
	sync.Mutex
	ctx context.Context
	log zerolog.Logger
	ov  network.Overlay
	// TODO: using a waitgroup here doesn't actually guarantee that we'll wait for all
	// goroutines to exit, because new goroutines could be started after we've already
	// returned from wg.Wait(). We need to solve this the right way using ComponentManager
	// and worker routines.
	wg                         *sync.WaitGroup
	libP2PNode                 *p2pnode.Node
	libP2PNodeFactory          p2pbuilder.LibP2PFactoryFunc
	preferredUnicasts          []unicast.ProtocolName
	me                         flow.Identifier
	metrics                    module.NetworkMetrics
	bitswapMetrics             module.BitswapMetrics
	rootBlockID                flow.Identifier
	validators                 []network.MessageValidator
	peerManagerFactory         p2p.PeerManagerFactoryFunc
	peerManager                p2p.PeerManager
	unicastMessageTimeout      time.Duration
	idTranslator               p2p.IDTranslator
	previousProtocolStatePeers []peer.AddrInfo
	codec                      network.Codec
	slashingViolationsConsumer slashing.ViolationsConsumer
	authorizedSenderValidator  *validator.AuthorizedSenderValidator
	component.Component
}

type MiddlewareOption func(*Middleware)

func WithMessageValidators(validators ...network.MessageValidator) MiddlewareOption {
	return func(mw *Middleware) {
		mw.validators = validators
	}
}

func WithPreferredUnicastProtocols(unicasts []unicast.ProtocolName) MiddlewareOption {
	return func(mw *Middleware) {
		mw.preferredUnicasts = unicasts
	}
}

func WithPeerManager(peerManagerFunc p2p.PeerManagerFactoryFunc) MiddlewareOption {
	return func(mw *Middleware) {
		mw.peerManagerFactory = peerManagerFunc
	}
}

// NewMiddleware creates a new middleware instance
// libP2PNodeFactory is the factory used to create a LibP2PNode
// flowID is this node's Flow ID
// metrics is the interface to report network related metrics
// peerUpdateInterval is the interval when the PeerManager's peer update runs
// unicastMessageTimeout is the timeout used for unicast messages
// connectionGating if set to True, restricts this node to only talk to other nodes which are part of the identity list
// managePeerConnections if set to True, enables the default PeerManager which continuously updates the node's peer connections
// validators are the set of the different message validators that each inbound messages is passed through
// During normal operations any error returned by Middleware.start is considered to be catastrophic
// and will be thrown by the irrecoverable.SignalerContext causing the node to crash.
func NewMiddleware(
	log zerolog.Logger,
	libP2PNodeFactory p2pbuilder.LibP2PFactoryFunc,
	flowID flow.Identifier,
	met module.NetworkMetrics,
	bitswapMet module.BitswapMetrics,
	rootBlockID flow.Identifier,
	unicastMessageTimeout time.Duration,
	idTranslator p2p.IDTranslator,
	codec network.Codec,
	slashingViolationsConsumer slashing.ViolationsConsumer,
	opts ...MiddlewareOption,
) *Middleware {

	if unicastMessageTimeout <= 0 {
		unicastMessageTimeout = DefaultUnicastTimeout
	}

	// create the node entity and inject dependencies & config
	mw := &Middleware{
		log:                        log,
		wg:                         &sync.WaitGroup{},
		me:                         flowID,
		libP2PNodeFactory:          libP2PNodeFactory,
		metrics:                    met,
		bitswapMetrics:             bitswapMet,
		rootBlockID:                rootBlockID,
		validators:                 DefaultValidators(log, flowID),
		unicastMessageTimeout:      unicastMessageTimeout,
		peerManagerFactory:         nil,
		idTranslator:               idTranslator,
		codec:                      codec,
		slashingViolationsConsumer: slashingViolationsConsumer,
	}

	for _, opt := range opts {
		opt(mw)
	}

	cm := component.NewComponentManagerBuilder().
		AddWorker(func(ctx irrecoverable.SignalerContext, ready component.ReadyFunc) {
			// TODO: refactor to avoid storing ctx altogether
			mw.ctx = ctx

			if err := mw.start(ctx); err != nil {
				ctx.Throw(err)
			}

			ready()

			<-ctx.Done()
			mw.stop()
		}).Build()

	mw.Component = cm

	return mw
}

func DefaultValidators(log zerolog.Logger, flowID flow.Identifier) []network.MessageValidator {
	return []network.MessageValidator{
		validator.ValidateNotSender(flowID),   // validator to filter out messages sent by this node itself
		validator.ValidateTarget(log, flowID), // validator to filter out messages not intended for this node
	}
}

// isProtocolParticipant returns a PeerFilter that returns true if a peer is a staked node.
func (m *Middleware) isProtocolParticipant() p2p.PeerFilter {
	return func(p peer.ID) error {
		if _, ok := m.ov.Identity(p); !ok {
			return fmt.Errorf("failed to get identity of unknown peer with peer id %s", p.Pretty())
		}
		return nil
	}
}

func (m *Middleware) NewBlobService(channel channels.Channel, ds datastore.Batching, opts ...network.BlobServiceOption) network.BlobService {
	return blob.NewBlobService(m.libP2PNode.Host(), m.libP2PNode.Routing(), channel.String(), ds, m.bitswapMetrics, m.log, opts...)
}

func (m *Middleware) NewPingService(pingProtocol protocol.ID, provider network.PingInfoProvider) network.PingService {
	return ping.NewPingService(m.libP2PNode.Host(), pingProtocol, m.log, provider)
}

func (m *Middleware) peerIDs(flowIDs flow.IdentifierList) peer.IDSlice {
	result := make([]peer.ID, len(flowIDs))

	for _, fid := range flowIDs {
		pid, err := m.idTranslator.GetPeerID(fid)
		if err != nil {
			// We probably don't need to fail the entire function here, since the other
			// translations may still succeed
			m.log.Err(err).Str("flowID", fid.String()).Msg("failed to translate to peer ID")
			continue
		}

		result = append(result, pid)
	}

	return result
}

// Me returns the flow identifier of this middleware
func (m *Middleware) Me() flow.Identifier {
	return m.me
}

// GetIPPort returns the ip address and port number associated with the middleware
// All errors returned from this function can be considered benign.
func (m *Middleware) GetIPPort() (string, string, error) {
	ipOrHostname, port, err := m.libP2PNode.GetIPPort()
	if err != nil {
		return "", "", fmt.Errorf("failed to get ip and port from libP2P node: %w", err)
	}

	return ipOrHostname, port, nil
}

func (m *Middleware) UpdateNodeAddresses() {
	m.log.Info().Msg("Updating protocol state node addresses")

	ids := m.ov.Identities()
	newInfos, invalid := utils.PeerInfosFromIDs(ids)

	for id, err := range invalid {
		m.log.Err(err).Str("node_id", id.String()).Msg("failed to extract peer info from identity")
	}

	m.Lock()
	defer m.Unlock()

	// set old addresses to expire
	for _, oldInfo := range m.previousProtocolStatePeers {
		m.libP2PNode.Host().Peerstore().SetAddrs(oldInfo.ID, oldInfo.Addrs, peerstore.TempAddrTTL)
	}

	for _, info := range newInfos {
		m.libP2PNode.Host().Peerstore().SetAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	}

	m.previousProtocolStatePeers = newInfos
}

func (m *Middleware) SetOverlay(ov network.Overlay) {
	m.ov = ov
}

// start will start the middleware.
// No errors are expected during normal operation.
func (m *Middleware) start(ctx context.Context) error {
	if m.ov == nil {
		return fmt.Errorf("could not start middleware: overlay must be configured by calling SetOverlay before middleware can be started")
	}

	m.authorizedSenderValidator = validator.NewAuthorizedSenderValidator(m.log, m.slashingViolationsConsumer, m.ov.Identity)

	libP2PNode, err := m.libP2PNodeFactory(ctx)
	if err != nil {
		return fmt.Errorf("could not create libp2p node: %w", err)
	}

	m.libP2PNode = libP2PNode
	err = m.libP2PNode.WithDefaultUnicastProtocol(m.handleIncomingStream, m.preferredUnicasts)
	if err != nil {
		return fmt.Errorf("could not register preferred unicast protocols on libp2p node: %w", err)
	}

	m.UpdateNodeAddresses()

	// create and use a peer manager if a peer manager factory was passed in during initialization
	if m.peerManagerFactory != nil {
		m.peerManager, err = m.peerManagerFactory(m.libP2PNode.Host(), func() peer.IDSlice {
			return m.topologyPeers()
		}, m.log)
		if err != nil {
			return fmt.Errorf("failed to create peer manager: %w", err)
		}

		select {
		case <-m.peerManager.Ready():
			m.log.Debug().Msg("peer manager successfully started")
		case <-time.After(30 * time.Second):
			return fmt.Errorf("could not start peer manager")
		}
	}

	return nil
}

// topologyPeers callback used by the peer manager to get the list of peer ID's
// which this node should be directly connected to as peers.
func (m *Middleware) topologyPeers() peer.IDSlice {
	return m.peerIDs(m.ov.Topology().NodeIDs())
}

// stop will end the execution of the middleware and wait for it to end.
func (m *Middleware) stop() {
	mgr, found := m.peerMgr()
	if found {
		// stops peer manager
		<-mgr.Done()
		m.log.Debug().Msg("peer manager successfully stopped")
	}

	// stops libp2p
	done, err := m.libP2PNode.Stop()
	if err != nil {
		m.log.Error().Err(err).Msg("could not stop libp2p node")
	} else {
		<-done
		m.log.Debug().Msg("libp2p node successfully stopped")
	}

	// wait for the readConnection and readSubscription routines to stop
	m.wg.Wait()
}

// SendDirect sends msg on a 1-1 direct connection to the target ID. It models a guaranteed delivery asynchronous
// direct one-to-one connection on the underlying network. No intermediate node on the overlay is utilized
// as the router.
//
// Dispatch should be used whenever guaranteed delivery to a specific target is required. Otherwise, Publish is
// a more efficient candidate.
//
// The following benign errors can be returned:
// - he peer ID for the target node ID cannot be found.
// - the msg size was too large.
// - failed to send message to peer.
//
// All errors returned from this function can be considered benign.
func (m *Middleware) SendDirect(msg *message.Message, targetID flow.Identifier) (err error) {
	// translates identifier to peer id
	peerID, err := m.idTranslator.GetPeerID(targetID)
	if err != nil {
		return fmt.Errorf("could not find peer id for target id: %w", err)
	}

	maxMsgSize := unicastMaxMsgSize(msg)
	if msg.Size() > maxMsgSize {
		// message size goes beyond maximum size that the serializer can handle.
		// proceeding with this message results in closing the connection by the target side, and
		// delivery failure.
		return fmt.Errorf("message size %d exceeds configured max message size %d", msg.Size(), maxMsgSize)
	}

	maxTimeout := m.unicastMaxMsgDuration(msg)

	m.metrics.DirectMessageStarted(msg.ChannelID)
	defer m.metrics.DirectMessageFinished(msg.ChannelID)

	// pass in a context with timeout to make the unicast call fail fast
	ctx, cancel := context.WithTimeout(m.ctx, maxTimeout)
	defer cancel()

	// protect the underlying connection from being inadvertently pruned by the peer manager while the stream and
	// connection creation is being attempted, and remove it from protected list once stream created.
	tag := fmt.Sprintf("%v:%v", msg.ChannelID, msg.Type)
	m.libP2PNode.Host().ConnManager().Protect(peerID, tag)
	defer m.libP2PNode.Host().ConnManager().Unprotect(peerID, tag)

	// create new stream
	// (streams don't need to be reused and are fairly inexpensive to be created for each send.
	// A stream creation does NOT incur an RTT as stream negotiation happens as part of the first message
	// sent out the receiver
	stream, err := m.libP2PNode.CreateStream(ctx, peerID)
	if err != nil {
		return fmt.Errorf("failed to create stream for %s: %w", targetID, err)
	}

	success := false

	defer func() {
		if success {
			// close the stream immediately
			err = stream.Close()
			if err != nil {
				err = fmt.Errorf("failed to close the stream for %s: %w", targetID, err)
			}

			// OneToOne communication metrics are reported with topic OneToOne
			m.metrics.NetworkMessageSent(msg.Size(), metrics.ChannelOneToOne, msg.Type)
		} else {
			resetErr := stream.Reset()
			if resetErr != nil {
				m.log.Err(resetErr).Msg("failed to reset stream")
			}
		}
	}()

	deadline, _ := ctx.Deadline()
	err = stream.SetWriteDeadline(deadline)
	if err != nil {
		return fmt.Errorf("failed to set write deadline for stream: %w", err)
	}

	// create a gogo protobuf writer
	bufw := bufio.NewWriter(stream)
	writer := ggio.NewDelimitedWriter(bufw)

	err = writer.WriteMsg(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to %s: %w", targetID, err)
	}

	// flush the stream
	err = bufw.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush stream for %s: %w", targetID, err)
	}

	success = true

	return nil
}

// handleIncomingStream handles an incoming stream from a remote peer
// it is a callback that gets called for each incoming stream by libp2p with a new stream object
func (m *Middleware) handleIncomingStream(s libp2pnetwork.Stream) {
	// qualify the logger with local and remote address
	log := p2putils.StreamLogger(m.log, s)

	log.Info().Msg("incoming stream received")

	success := false

	remotePeer := s.Conn().RemotePeer()

	defer func() {
		if success {
			err := s.Close()
			if err != nil {
				log.Err(err).Msg("failed to close stream")
			}
		} else {
			err := s.Reset()
			if err != nil {
				log.Err(err).Msg("failed to reset stream")
			}
		}
	}()

	// TODO: We need to allow per-topic timeouts and message size limits.
	// This allows us to configure higher limits for topics on which we expect
	// to receive large messages (e.g. Chunk Data Packs), and use the normal
	// limits for other topics. In order to enable this, we will need to register
	// a separate stream handler for each topic.
	ctx, cancel := context.WithTimeout(m.ctx, LargeMsgUnicastTimeout)
	defer cancel()

	deadline, _ := ctx.Deadline()

	err := s.SetReadDeadline(deadline)
	if err != nil {
		log.Err(err).Msg("failed to set read deadline for stream")
		return
	}

	// create the reader
	r := ggio.NewDelimitedReader(s, LargeMsgMaxUnicastMsgSize)

	for {
		if ctx.Err() != nil {
			return
		}

		var msg message.Message

		// read the next message (blocking call)
		err = r.ReadMsg(&msg)

		if err != nil {
			if err == io.EOF {
				break
			}

			m.log.Err(err).Msg("failed to read message")
			return
		}

		// TODO: once we've implemented per topic message size limits per the TODO above,
		// we can remove this check
		maxSize := unicastMaxMsgSize(&msg)
		if msg.Size() > maxSize {
			// message size exceeded
			m.log.Error().
				Hex("sender", msg.OriginID).
				Hex("event_id", msg.EventID).
				Str("event_type", msg.Type).
				Str("channel", msg.ChannelID).
				Int("maxSize", maxSize).
				Msg("received message exceeded permissible message maxSize")
			return
		}

		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.processUnicastStreamMessage(remotePeer, &msg)
		}()
	}

	success = true
}

// Subscribe subscribes the middleware to a channel.
// No errors are expected during normal operation.
func (m *Middleware) Subscribe(channel channels.Channel) error {

	topic := channels.TopicFromChannel(channel, m.rootBlockID)

	var peerFilter p2p.PeerFilter
	var validators []validator.PubSubMessageValidator
	if channels.IsPublicChannel(channel) {
		// NOTE: for public channels the callback used to check if a node is staked will
		// return true for every node.
		peerFilter = p2putils.AllowAllPeerFilter()
	} else {
		// for channels used by the staked nodes, add the topic validator to filter out messages from non-staked nodes
		validators = append(validators,
			m.authorizedSenderValidator.PubSubMessageValidator(channel),
		)

		// NOTE: For non-public channels the libP2P node topic validator will reject
		// messages from unstaked nodes.
		peerFilter = m.isProtocolParticipant()
	}

	s, err := m.libP2PNode.Subscribe(topic, m.codec, peerFilter, m.slashingViolationsConsumer, validators...)
	if err != nil {
		return fmt.Errorf("could not subscribe to topic (%s): %w", topic, err)
	}

	// create a new readSubscription with the context of the middleware
	rs := newReadSubscription(m.ctx, s, m.processAuthenticatedMessage, m.log, m.metrics)
	m.wg.Add(1)

	// kick off the receive loop to continuously receive messages
	go rs.receiveLoop(m.wg)

	// update peers to add some nodes interested in the same topic as direct peers
	m.peerManagerUpdate()

	return nil
}

// Unsubscribe unsubscribes the middleware from a channel.
// The following benign errors are expected during normal operations from libP2P:
// - the libP2P node fails to unsubscribe to the topic created from the provided channel.
//
// All errors returned from this function can be considered benign.
func (m *Middleware) Unsubscribe(channel channels.Channel) error {
	topic := channels.TopicFromChannel(channel, m.rootBlockID)
	err := m.libP2PNode.UnSubscribe(topic)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from channel (%s): %w", channel, err)
	}

	// update peers to remove nodes subscribed to channel
	m.peerManagerUpdate()

	return nil
}

// processUnicastStreamMessage will decode, perform authorized sender validation and process a message
// sent via unicast stream. This func should be invoked in a separate goroutine to avoid creating a message decoding bottleneck.
func (m *Middleware) processUnicastStreamMessage(remotePeer peer.ID, msg *message.Message) {
	channel := channels.Channel(msg.ChannelID)

	decodedMsgPayload, err := m.codec.Decode(msg.Payload)
	if codec.IsErrUnknownMsgCode(err) {
		// slash peer if message contains unknown message code byte
		violation := &slashing.Violation{PeerID: remotePeer.String(), Channel: channel, IsUnicast: true, Err: err}
		m.slashingViolationsConsumer.OnUnknownMsgTypeError(violation)
		return
	}
	if codec.IsErrMsgUnmarshal(err) {
		// slash if peer sent a message that could not be marshalled into the message type denoted by the message code byte
		violation := &slashing.Violation{PeerID: remotePeer.String(), Channel: channel, IsUnicast: true, Err: err}
		m.slashingViolationsConsumer.OnInvalidMsgError(violation)
		return
	}

	// unexpected error condition. this indicates there's a bug
	// don't crash as a result of external inputs since that creates a DoS vector.
	if err != nil {
		m.log.
			Error().
			Err(fmt.Errorf("unexpected error while decoding message: %w", err)).
			Hex("sender", msg.OriginID).
			Hex("event_id", msg.EventID).
			Str("event_type", msg.Type).
			Str("channel", msg.ChannelID).
			Msg("failed to decode message payload")
		return
	}

	// if message channel is not public perform authorized sender validation
	if !channels.IsPublicChannel(channel) {
		_, err := m.authorizedSenderValidator.Validate(remotePeer, decodedMsgPayload, channel, true)
		if err != nil {
			m.log.
				Error().
				Err(err).
				Hex("sender", msg.OriginID).
				Hex("event_id", msg.EventID).
				Str("event_type", msg.Type).
				Str("channel", msg.ChannelID).
				Msg("unicast authorized sender validation failed")
			return
		}
	}

	// message decoding and validation was successful now log metrics with the channel name as OneToOne and process message
	m.metrics.NetworkMessageReceived(msg.Size(), metrics.ChannelOneToOne, msg.Type)
	m.processAuthenticatedMessage(msg, decodedMsgPayload, remotePeer)
}

// processAuthenticatedMessage processes a message and a source (indicated by its peer ID) and eventually passes it to the overlay
// In particular, it populates the `OriginID` field of the message with a Flow ID translated from this source.
// The assumption is that the message has been authenticated at the network level (libp2p) to originate from the peer with ID `peerID`
// this requirement is fulfilled by e.g. the output of readConnection and readSubscription
func (m *Middleware) processAuthenticatedMessage(msg *message.Message, decodedMsgPayload interface{}, peerID peer.ID) {
	flowID, err := m.idTranslator.GetFlowID(peerID)
	if err != nil {
		m.log.Warn().Err(err).Msgf("received message from unknown peer %v, and was dropped", peerID.String())
		return
	}

	msg.OriginID = flowID[:]

	m.processMessage(msg, decodedMsgPayload)
}

// processMessage processes a message and eventually passes it to the overlay
func (m *Middleware) processMessage(msg *message.Message, decodedMsgPayload interface{}) {
	originID := flow.HashToID(msg.OriginID)

	m.log.Debug().
		Str("channel", msg.ChannelID).
		Str("type", msg.Type).
		Str("origin_id", originID.String()).
		Msg("processing new message")

	// run through all the message validators
	for _, v := range m.validators {
		// if any one fails, stop message propagation
		if !v.Validate(*msg) {
			return
		}
	}

	// if validation passed, send the message to the overlay
	err := m.ov.Receive(originID, msg, decodedMsgPayload)
	if err != nil {
		m.log.Error().Err(err).Msg("could not deliver payload")
	}
}

// Publish publishes a message on the channel. It models a distributed broadcast where the message is meant for all or
// a many nodes subscribing to the channel. It does not guarantee the delivery though, and operates on a best
// effort.
// The following benign errors are expected during normal operations:
// - the msg cannot be marshalled.
// - the msg size exceeds DefaultMaxPubSubMsgSize.
// - the libP2P node fails to publish the message.
//
// All errors returned from this function can be considered benign.
func (m *Middleware) Publish(msg *message.Message, channel channels.Channel) error {
	m.log.Debug().Str("channel", channel.String()).Interface("msg", msg).Msg("publishing new message")

	// convert the message to bytes to be put on the wire.
	//bs := binstat.EnterTime(binstat.BinNet + ":wire<4message2protobuf")
	data, err := msg.Marshal()
	//binstat.LeaveVal(bs, int64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to marshal the message: %w", err)
	}

	msgSize := len(data)
	if msgSize > p2pnode.DefaultMaxPubSubMsgSize {
		// libp2p pubsub will silently drop the message if its size is greater than the configured pubsub max message size
		// hence return an error as this message is undeliverable
		return fmt.Errorf("message size %d exceeds configured max message size %d", msgSize, p2pnode.DefaultMaxPubSubMsgSize)
	}

	topic := channels.TopicFromChannel(channel, m.rootBlockID)

	// publish the bytes on the topic
	err = m.libP2PNode.Publish(m.ctx, topic, data)
	if err != nil {
		return fmt.Errorf("failed to publish the message: %w", err)
	}

	m.metrics.NetworkMessageSent(len(data), string(channel), msg.Type)

	return nil
}

// IsConnected returns true if this node is connected to the node with id nodeID.
// All errors returned from this function can be considered benign.
func (m *Middleware) IsConnected(nodeID flow.Identifier) (bool, error) {
	peerID, err := m.idTranslator.GetPeerID(nodeID)
	if err != nil {
		return false, fmt.Errorf("could not find peer id for target id: %w", err)
	}
	return m.libP2PNode.IsConnected(peerID)
}

// unicastMaxMsgSize returns the max permissible size for a unicast message
func unicastMaxMsgSize(msg *message.Message) int {
	switch msg.Type {
	case "messages.ChunkDataResponse":
		return LargeMsgMaxUnicastMsgSize
	default:
		return DefaultMaxUnicastMsgSize
	}
}

// unicastMaxMsgDuration returns the max duration to allow for a unicast send to complete
func (m *Middleware) unicastMaxMsgDuration(msg *message.Message) time.Duration {
	switch msg.Type {
	case "messages.ChunkDataResponse":
		if LargeMsgUnicastTimeout > m.unicastMessageTimeout {
			return LargeMsgUnicastTimeout
		}
		return m.unicastMessageTimeout
	default:
		return m.unicastMessageTimeout
	}
}

// peerManagerUpdate request an update from the peer manager to connect to new peers and disconnect from unwanted peers
func (m *Middleware) peerManagerUpdate() {
	mgr, found := m.peerMgr()
	if found {
		mgr.RequestPeerUpdate()
	}
}

// peerMgr returns the PeerManager and true if this middleware was started with one, (nil, false) otherwise
func (m *Middleware) peerMgr() (p2p.PeerManager, bool) {
	if m.peerManager != nil {
		return m.peerManager, true
	}
	return nil, false
}
