// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

// Constants
const PORT = 8080;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

async function cc_call(fn_name, args){

    // chaincode 호출 파트

    // load the network configuration
    const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
    let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), "wallet");
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const identity = await wallet.get("appUser1");
    if (!identity) {
        console.log(
            'An identity for the user "appUser1" does not exist in the wallet'
        );
        console.log("Run the registerUser.js application before retrying");
        return;
    }

    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: "appUser1",
        discovery: { enabled: true, asLocalhost: true },
    });

    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork("mychannel");

    // Get the contract from the network.
    const contract = network.getContract("sYacht");



    // Submit the specified transaction.
    // var result;

    // var result

    if (fn_name == "registYacht") {
        await contract.submitTransaction("registYacht", args[0],args[1],args[2],args[3],args[4]);       
    } else if (fn_name == "readYacht"){
        result = await contract.evaluateTransaction("readYacht", args);
    }
     else result = "not supported fuction!"
     
    await gateway.disconnect();

    return result
}

app.post("/registYacht", async(req, res) => {
    const username = req.body.username
    const name = req.body.yname
    const color = req.body.ycolor
    const size = req.body.ysize
    const price = req.body.yprice

    var args = [username, name, color, size, price];

    result = await cc_call("registYacht", args)

    const resultPath = path.join(process.cwd(), "/views/result.html")
    var resultHTML =  fs.readFileSync(resultPath, "utf-8")

    // console.log(result)    
          
    resultHTML = resultHTML.replace("<dir></dir>", "요트등록에 성공하였습니다!")        
    
    res.status(200).send(resultHTML)

    console.log("reqest : " + username + ":" + name + ":" + color + ":" + size + ":" + price )
});

app.get("/readYacht", async(req, res) => {
    const name = req.query.name

    result = await cc_call("readYacht", name)


    console.log("reqest : " + name )
    console.log("response : " + result.toString())
    const result_str = `{"result":"Done /readYacht !!", "msg":${result.toString()}}`
    res.status(200).json(JSON.parse(result_str))
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
