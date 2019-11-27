{
    "version": "1",
    "predistribution": [
        {
            "address": "f3prTg9itaZY6m48wXXikXdcxiByW7zgk",
            "quota": "100000000000000000000"
        }
    ],
    "maxblocksize": "128",
    "award": "1000000",
    "decimals": "8",
    "award_decay": {
        "height_gap": 31536000,
        "ratio": 1
    },
    "genesis_consensus": {
        "name": "tdpos",
        "config": {
            "timestamp": "1566830669000000000",
            "proposer_num": "3",
            "period": "3000",
            "alternate_interval": "9000",
            "term_interval": "9000",
            "block_num": "5",
            "vote_unit_price": "1",
            "init_proposer": {
                "1": [
                    "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN",
                    "f3prTg9itaZY6m48wXXikXdcxiByW7zgk",
                    "U9sKwFmgJVfzgWcfAG47dKn1kLQTqeZN3"
                ]
            },
            "init_proposer_neturl": {
                "1": [
		    "/ip4/127.0.0.1/tcp/47101/p2p/QmTdHf8shB5MY1Ur8w7mPjHcxKGL6CCfE43EpYn2Z8fMzy",
		    "/ip4/127.0.0.1/tcp/47102/p2p/QmeSjPfpn9D6oa7Kpr4aMR4yUhWTX9vre53PjjqzJFYZYW",
		    "/ip4/127.0.0.1/tcp/47103/p2p/QmPYPuXsNXYRfVyr6FsingX8gwFW3FwgJDtuJW1iYugwJn"
                ]
            },
            "bft_config": {
            }
        }
    }
}
