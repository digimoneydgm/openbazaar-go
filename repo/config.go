package repo

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/config"
)

var DefaultBootstrapAddresses = []string{
	"/ip4/159.203.115.78/tcp/5001/ipfs/QmPJuP4Myo8pGL1k56b85Q4rpaoSnmn5L3wLjYHTzbBrk1",
	"/ip4/104.131.19.44/tcp/5001/ipfs/QmZobfuuBvzQUA6xxuES63rpjQtGx3JbUbWxe7p3f7nQTM",
	"/ip4/45.32.171.119/tcp/5001/ipfs/QmbkmTrfSjniyMUfYt9tVUcfZC8yeshVFYsnJ4kuzveYWc",
}

var TestnetBootstrapAddresses = []string{}

var DataPushNodes = []string{}

type APIConfig struct {
	Authenticated bool
	AllowedIPs    []string
	Username      string
	Password      string
	CORS          *string
	Enabled       bool
	HTTPHeaders   map[string]interface{}
	SSL           bool
	SSLCert       string
	SSLKey        string
}

type TorConfig struct {
	Password   string
	TorControl string
}

type ResolverConfig struct {
	Id  string `json:".id"`
	Eth string `json:".eth"`
}

type WalletConfig struct {
	Type             string
	MaxFee           int
	HighFeeDefault   int
	MediumFeeDefault int
	LowFeeDefault    int
	RPCLocation      string
}

type DataSharing struct {
	AcceptStoreRequests bool
	PushTo              []string
}

var MalformedConfigError error = errors.New("Config file is malformed")

func GetAPIConfig(cfgBytes []byte) (*APIConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	apiIface, ok := cfg["JSON-API"]
	if !ok {
		return nil, MalformedConfigError
	}

	api, ok := apiIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	headers := make(map[string]interface{})
	h, ok := api["HTTPHeaders"]
	if h == nil || !ok {
		headers = nil
	} else {
		headers, ok = h.(map[string]interface{})
		if !ok {
			return nil, MalformedConfigError
		}
	}

	enabled, ok := api["Enabled"]
	if !ok {
		return nil, MalformedConfigError
	}
	enabledBool, ok := enabled.(bool)
	if !ok {
		return nil, MalformedConfigError
	}
	authenticated := api["Authenticated"]
	if !ok {
		return nil, MalformedConfigError
	}
	authenticatedBool, ok := authenticated.(bool)
	if !ok {
		return nil, MalformedConfigError
	}
	allowedIPs, ok := api["AllowedIPs"]
	if !ok {
		return nil, MalformedConfigError
	}
	allowedIPsIface, ok := allowedIPs.([]interface{})
	if !ok {
		return nil, MalformedConfigError
	}
	var allowedIPstrings []string
	for _, ip := range allowedIPsIface {
		ipStr, ok := ip.(string)
		if !ok {
			return nil, MalformedConfigError
		}
		allowedIPstrings = append(allowedIPstrings, ipStr)
	}

	username, ok := api["Username"]
	if !ok {
		return nil, MalformedConfigError
	}
	usernameStr, ok := username.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	password, ok := api["Password"]
	if !ok {
		return nil, MalformedConfigError
	}
	passwordStr, ok := password.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	c, ok := api["CORS"]
	var cors *string
	if c == nil || !ok {
		cors = nil
	} else {
		crs, ok := c.(string)
		if !ok {
			return nil, MalformedConfigError
		}
		cors = &crs
	}
	sslEnabled, ok := api["SSL"]
	if !ok {
		return nil, MalformedConfigError
	}
	sslEnabledBool, ok := sslEnabled.(bool)
	if !ok {
		return nil, MalformedConfigError
	}

	certFile, ok := api["SSLCert"]
	if !ok {
		return nil, MalformedConfigError
	}
	certFileStr, ok := certFile.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	keyFile, ok := api["SSLKey"]
	if !ok {
		return nil, MalformedConfigError
	}
	keyFileStr, ok := keyFile.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	apiConfig := &APIConfig{
		Authenticated: authenticatedBool,
		AllowedIPs:    allowedIPstrings,
		Username:      usernameStr,
		Password:      passwordStr,
		CORS:          cors,
		Enabled:       enabledBool,
		HTTPHeaders:   headers,
		SSL:           sslEnabledBool,
		SSLCert:       certFileStr,
		SSLKey:        keyFileStr,
	}

	return apiConfig, nil
}

func GetWalletConfig(cfgBytes []byte) (*WalletConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)
	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	walletIface, ok := cfg["Wallet"]
	if !ok {
		return nil, MalformedConfigError
	}
	wallet, ok := walletIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}
	low, ok := wallet["LowFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	lowFloat, ok := low.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	medium, ok := wallet["MediumFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	mediumFloat, ok := medium.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	high, ok := wallet["HighFeeDefault"]
	if !ok {
		return nil, MalformedConfigError
	}
	highFloat, ok := high.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	maxFee, ok := wallet["MaxFee"]
	if !ok {
		return nil, MalformedConfigError
	}
	maxFeeFloat, ok := maxFee.(float64)
	if !ok {
		return nil, MalformedConfigError
	}
	walletType, ok := wallet["Type"]
	if !ok {
		return nil, MalformedConfigError
	}
	walletTypeStr, ok := walletType.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	rpcLocation, ok := wallet["RPCLocation"]
	if !ok {
		return nil, MalformedConfigError
	}
	rpcLocationStr, ok := rpcLocation.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	wCfg := &WalletConfig{
		Type:             walletTypeStr,
		MaxFee:           int(maxFeeFloat),
		HighFeeDefault:   int(highFloat),
		MediumFeeDefault: int(mediumFloat),
		LowFeeDefault:    int(lowFloat),
		RPCLocation:      rpcLocationStr,
	}
	return wCfg, nil
}

func GetTorConfig(cfgBytes []byte) (*TorConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	tcIface, ok := cfg["Tor-config"]
	if !ok {
		return nil, MalformedConfigError
	}
	tc, ok := tcIface.(map[string]interface{})

	pw, ok := tc["Password"]
	if !ok {
		return nil, MalformedConfigError
	}
	pwStr, ok := pw.(string)
	if !ok {
		return nil, MalformedConfigError
	}
	controlUrl, ok := tc["TorControl"]
	if !ok {
		return nil, MalformedConfigError
	}
	controlUrlStr, ok := controlUrl.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	return &TorConfig{TorControl: controlUrlStr, Password: pwStr}, nil
}

func GetDropboxApiToken(cfgBytes []byte) (string, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return "", MalformedConfigError
	}

	token, ok := cfg["Dropbox-api-token"]
	if !ok {
		return "", MalformedConfigError
	}
	tokenStr, ok := token.(string)
	if !ok {
		return "", MalformedConfigError
	}

	return tokenStr, nil
}

func GetRepublishInterval(cfgBytes []byte) (time.Duration, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return time.Duration(0), MalformedConfigError
	}

	interval, ok := cfg["RepublishInterval"]
	if !ok {
		return time.Duration(0), MalformedConfigError
	}
	intervalStr, ok := interval.(string)
	if !ok {
		return time.Duration(0), MalformedConfigError
	}
	if intervalStr == "" {
		return time.Duration(0), nil
	}
	d, err := time.ParseDuration(intervalStr)
	if err != nil {
		return d, err
	}
	return d, nil
}

func GetDataSharing(cfgBytes []byte) (*DataSharing, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)
	dataSharing := new(DataSharing)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return dataSharing, MalformedConfigError
	}

	dscfg, ok := cfg["DataSharing"]
	if !ok {
		return dataSharing, MalformedConfigError
	}
	ds, ok := dscfg.(map[string]interface{})
	if !ok {
		return dataSharing, MalformedConfigError
	}

	acceptcfg, ok := ds["AcceptStoreRequests"]
	if !ok {
		return dataSharing, MalformedConfigError
	}
	accept, ok := acceptcfg.(bool)
	if !ok {
		return dataSharing, MalformedConfigError
	}
	dataSharing.AcceptStoreRequests = accept

	pushcfg, ok := ds["PushTo"]
	if !ok {
		return dataSharing, MalformedConfigError
	}
	pushList, ok := pushcfg.([]interface{})
	if !ok {
		return dataSharing, MalformedConfigError
	}

	for _, nd := range pushList {
		ndStr, ok := nd.(string)
		if !ok {
			return dataSharing, MalformedConfigError
		}
		dataSharing.PushTo = append(dataSharing.PushTo, ndStr)
	}
	return dataSharing, nil
}

func GetTestnetBootstrapAddrs(cfgBytes []byte) ([]string, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)
	var addrs []string

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return addrs, MalformedConfigError
	}

	bootstrap, ok := cfg["Bootstrap-testnet"]
	if !ok {
		return addrs, MalformedConfigError
	}
	addrList, ok := bootstrap.([]interface{})
	if !ok {
		return addrs, MalformedConfigError
	}

	for _, addr := range addrList {
		addrStr, ok := addr.(string)
		if !ok {
			return addrs, MalformedConfigError
		}
		addrs = append(addrs, addrStr)
	}

	return addrs, nil
}

func GetResolverConfig(cfgBytes []byte) (*ResolverConfig, error) {
	var cfgIface interface{}
	json.Unmarshal(cfgBytes, &cfgIface)

	cfg, ok := cfgIface.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}

	r, ok := cfg["Resolvers"]
	if !ok {
		return nil, MalformedConfigError
	}
	resolverMap, ok := r.(map[string]interface{})
	if !ok {
		return nil, MalformedConfigError
	}
	blockstack, ok := resolverMap[".id"]
	if !ok {
		return nil, MalformedConfigError
	}

	idStr, ok := blockstack.(string)
	if !ok {
		return nil, MalformedConfigError
	}

	resolvers := &ResolverConfig{
		Id: idStr,
	}

	return resolvers, nil
}

func extendConfigFile(r repo.Repo, key string, value interface{}) error {
	if err := r.SetConfigKey(key, value); err != nil {
		return err
	}
	return nil
}

func InitConfig(repoRoot string) (*config.Config, error) {
	bootstrapPeers, err := config.ParseBootstrapPeers(DefaultBootstrapAddresses)
	if err != nil {
		return nil, err
	}

	datastore := datastoreConfig(repoRoot)

	conf := &config.Config{

		// Setup the node's default addresses.
		// NOTE: two swarm listen addrs, one TCP, one UTP.
		Addresses: config.Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/5001",
				"/ip6/::/tcp/5001",
				"/ip4/0.0.0.0/tcp/10005/ws",
				"/ip6/::/tcp/10005/ws",
			},
			API:     "",
			Gateway: "/ip4/127.0.0.1/tcp/5002",
		},

		Datastore: datastore,
		Bootstrap: config.BootstrapPeerStrings(bootstrapPeers),
		Discovery: config.Discovery{config.MDNS{
			Enabled:  false,
			Interval: 10,
		}},

		// Setup the node mount points
		Mounts: config.Mounts{
			IPFS: "/ipfs",
			IPNS: "/ipns",
		},

		Ipns: config.Ipns{
			ResolveCacheSize:   128,
			RecordLifetime:     "7d",
			RepublishPeriod:    "24h",
			QuerySize:          5,
			UsePersistentCache: true,
		},

		Gateway: config.Gateway{
			RootRedirect: "",
			Writable:     false,
			PathPrefixes: []string{},
		},
	}

	return conf, nil
}

func datastoreConfig(repoRoot string) config.Datastore {
	return config.Datastore{
		StorageMax:         "10GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
		BloomFilterSize:    0,
		HashOnRead:         false,
		Spec: map[string]interface{}{
			"type": "mount",
			"mounts": []interface{}{
				map[string]interface{}{
					"mountpoint": "/blocks",
					"type":       "measure",
					"prefix":     "flatfs.datastore",
					"child": map[string]interface{}{
						"type":      "flatfs",
						"path":      "blocks",
						"sync":      true,
						"shardFunc": "/repo/flatfs/shard/v1/next-to-last/2",
					},
				},
				map[string]interface{}{
					"mountpoint": "/",
					"type":       "measure",
					"prefix":     "leveldb.datastore",
					"child": map[string]interface{}{
						"type":        "levelds",
						"path":        "datastore",
						"compression": "none",
					},
				},
			},
		},
	}
}
