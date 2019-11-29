{
    "version": "1",
    "predistribution": [
        {
            "address": "{{.PredistributionAddr}}",
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
            "proposer_num": "{{.ProposerNum}}",
            "period": "3000",
            "alternate_interval": "9000",
            "term_interval": "9000",
            "block_num": "5",
            "vote_unit_price": "1",
            "init_proposer": {
                "1": [
                    {{.InitProposer}}
                ]
            },
            "init_proposer_neturl": {
                "1": [
		    {{.InitProposerNeturl}}
                ]
            },
            "bft_config": {
            }
        }
    }
}
