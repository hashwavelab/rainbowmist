// this is fully working！！！

const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const grpc_promise = require('grpc-promise');
 
const packageDefinition1 = protoLoader.loadSync(
    '/mnt/raid/rainbowmist/pb/pb.proto',
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  }
);
const protoDescriptor1 = grpc.loadPackageDefinition(packageDefinition1);
const pb = protoDescriptor1.rainbowmist;
 
async function main() {
  const client = new pb.Rainbowmist('localhost:8889', grpc.credentials.createInsecure());
  grpc_promise.promisifyAll(client);

  console.time('a')
  let res = await client.GetPrice()
    .sendMessage({
        base_asset: "AVAX",
        quote_asset: "USDT",
        decimals: "6",
    })

    console.log(res);
    console.timeEnd('a')
}
 
main();