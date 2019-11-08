geth --datadir node1/ --syncmode 'full' --port 30312 --rpcaddr 'localhost' --rpcport 8501 --rpcapi 'personal,db,eth,net,web3,txpool,miner' --bootnodes 'enode://f9476d61724b137e919ccf32035924ccbd23daae5448771d7ed941b6066b6ace56e1fc553844d1eb0f21e337ddfd07ef8b13406021198587388eac3c4d88fa29@127.0.0.1:30310' --networkid 1415 --gasprice '1' --unlock '0x3fdf7322063a1492737e135bc6ad1b8ee9ef18a3' --password node1/password.txt --mine