# valsta
For cosmos blockchain validator to collect validator information

## Easy Start
```bash
# build valsta
make build
# start at height 400000
./build/valsta start 841500 1412246  -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
# query voters
./build/valsta q voters -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
# query sign times 
./build/valsta q signTimes <address> -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
# query sign first height 
./build/valsta q signHeight <address> -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
```
uptick1g9qle0zayz3mjnel3f9du7wdrddclh4hxnwcal
## Use Guide
There are only two steps required to use valsta
### Init
use `valsta init` to create database `valsta` and table `valdator_infos`
### Start
use `valsta start [startHeight] [endHeight] -g <grpc address> -r <rpc address>` to start valsta.
If you want to have a test with your local cosmos-base blockchain , just use `valsta start [startHeight] [endHeight]`,
the flag `-g` have default value `localhost:9090` and the flag `-r` have default value `http://localhost:26657`
### Start with nohup
```bash
nohup ./build/valsta start 841500 1412246  -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/  > work.log 2>&1 & 
```