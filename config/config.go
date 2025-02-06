package config

type Symbol struct {
	FastexName      string  `json:"fastex_name"`
	HostSymbolName  string  `json:"host_symbol_name"`
	HostName        string  `json:"host_name"`
	Precision       int     `json:"precision"`
	Step            float64 `json:"step"`
	Volume          float64 `json:"volume"`
	VolumePrecision int     `json:"volume_precision"`
}

var Symbols = map[string]Symbol{
	"LTC-USDT": {
		FastexName:      "LTC-USDT",
		HostSymbolName:  "LTCUSDT",
		HostName:        "binance",
		Precision:       3,
		Step:            0.001,
		Volume:          1,
		VolumePrecision: 3,
	},
	"SOL-ETH": {
		FastexName:      "SOL-ETH",
		HostSymbolName:  "SOLETH",
		HostName:        "binance",
		Precision:       5,
		Step:            0.001,
		Volume:          1,
		VolumePrecision: 3,
	},
	"SHIB-USDT": {
		FastexName:      "SHIB-USDT",
		HostSymbolName:  "SHIBUSDT",
		HostName:        "binance",
		Precision:       8,
		Step:            0.00000001,
		Volume:          1000000,
		VolumePrecision: 0,
	},
	"DOGE-USDT": {
		FastexName:      "DOGE-USDT",
		HostSymbolName:  "DOGEUSDT",
		HostName:        "binance",
		Precision:       5,
		Step:            0.00001,
		Volume:          10,
		VolumePrecision: 0,
	},
	"XRP-USDT": {
		FastexName:      "XRP-USDT",
		HostSymbolName:  "XRPUSDT",
		HostName:        "binance",
		Precision:       5,
		Step:            0.00001,
		Volume:          10,
		VolumePrecision: 1,
	},
	"TRX-USDT": {
		FastexName:      "TRX-USDT",
		HostSymbolName:  "TRXUSDT",
		HostName:        "binance",
		Precision:       4,
		Step:            0.0001,
		Volume:          1,
		VolumePrecision: 1,
	},
	"FTN-USDT": {
		FastexName:      "FTN-USDT",
		HostSymbolName:  "FTNUSDT",
		HostName:        "bitget",
		Precision:       4,
		Step:            0.0001,
		Volume:          1,
		VolumePrecision: 2,
	},
	"ETH-USDT": {
		FastexName:      "ETH-USDT",
		HostSymbolName:  "ETHUSDT",
		HostName:        "binance",
		Precision:       2,
		Step:            0.01,
		Volume:          0.1,
		VolumePrecision: 2,
	},
	"BTC-USDT": {
		FastexName:      "BTC-USDT",
		HostSymbolName:  "BTCUSDT",
		HostName:        "binance",
		Precision:       2,
		Step:            0.01,
		Volume:          0.1,
		VolumePrecision: 2,
	},
	"FTN-USDC": {
		FastexName:      "FTN-USDC",
		HostSymbolName:  "FTNUSDT",
		HostName:        "bitget",
		Precision:       2,
		Step:            0.01,
		Volume:          0.1,
		VolumePrecision: 2,
	},
	"ETH-USDC": {
		FastexName:      "ETH-USDC",
		HostName:        "binance",
		HostSymbolName:  "ETHUSDC",
		Precision:       2,
		Step:            0.01,
		Volume:          0.1,
		VolumePrecision: 4,
	},
}
