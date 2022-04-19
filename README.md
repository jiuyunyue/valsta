# valsta
For cosmos blockchain validator to collect validator information

## Easy Start
```bash
# mysql docker
bash mysql.sh
# build valsta
make build
# init database
./build/valsta init
# start at height 400000
./build/valsta start 400000 400000 -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
# query val
./build/valsta q val
# query voters
./build/valsta q voters -g peer0.testnet.uptick.network:9090 -r http://peer0.testnet.uptick.network:26657/
# clean database
./build/valsta clean
```
## Use Guide
There are only two steps required to use valsta
### Init
use `valsta init` to create database `valsta` and table `valdator_infos`
### Start
use `valsta start [startHeight] [endHeight] -g <grpc address> -r <rpc address>` to start valsta.
If you want to have a test with your local cosmos-base blockchain , just use `valsta start [startHeight] [endHeight]`,
the flag `-g` have default value `localhost:9090` and the flag `-r` have default value `http://localhost:26657`

