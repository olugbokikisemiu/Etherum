geth --datadir node2/ --syncmode 'full' --port 30311 --rpc --rpcaddr 'localhost' --rpcport 8501 --rpcapi 'personal,db,eth,net,web3,txpool,miner' --bootnodes 'enode://f9476d61724b137e919ccf32035924ccbd23daae5448771d7ed941b6066b6ace56e1fc553844d1eb0f21e337ddfd07ef8b13406021198587388eac3c4d88fa29@127.0.0.1:30310' --networkid 1415 --gasprice '0' --password node2/password.txt --mine