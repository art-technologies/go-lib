package events

// Methods
const (
	MethodHead           = "SubscribeToHead"
	MethodBlocks         = "SubscribeToBlocks"
	MethodOperations     = "SubscribeToOperations"
	MethodBigMap         = "SubscribeToBigMaps"
	MethodAccounts       = "SubscribeToAccounts"
	MethodTokenTransfers = "SubscribeToTokenTransfers"
	MethodEvents         = "SubscribeToEvents"
)

// Channels
const (
	ChannelHead       = "head"
	ChannelBlocks     = "blocks"
	ChannelOperations = "operations"
	ChannelBigMap     = "bigmaps"
	ChannelAccounts   = "accounts"
	ChannelTransfers  = "transfers"
	ChannelEvents     = "events"
)

// Big map tags
const (
	BigMapTagMetadata      = "metadata"
	BigMapTagTokenMetadata = "token_metadata"
)
