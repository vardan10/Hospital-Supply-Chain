/**
 * Copyright 2017 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
'use strict';
var creds = require('./creds.json')
var path = require('path');
var fs = require('fs');
var util = require('util');
var Fabric_Client = require('fabric-client');
var helper = require('./helper.js');
var logger = helper.getLogger('invoke-chaincode');

var invokeChaincode = async function(peerNames, channelName, chaincodeName, fcn, args, username, org_name) {
	try{
		var fabric_client = new Fabric_Client();

		// setup the fabric network
		var channel = fabric_client.newChannel(channelName);
		var peer = fabric_client.newPeer(creds.peers["org1-peer1"].url, { pem: creds.peers["org1-peer1"].tlsCACerts.pem , 'ssl-target-name-override': null});
		channel.addPeer(peer);
		var order = fabric_client.newOrderer(creds.orderers.orderer.url, { pem: creds.orderers.orderer.tlsCACerts.pem , 'ssl-target-name-override': null})
		channel.addOrderer(order);

		console.log(channel);

		//
		var member_user = null;
		var store_path = path.join(__dirname, '/../fabric-client-kv-org1');
		console.log('Store path:'+store_path);
		var tx_id = null;

		// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
		var state_store = await Fabric_Client.newDefaultKeyValueStore({ path: store_path})
		// assign the store to the fabric client
		fabric_client.setStateStore(state_store);
		var crypto_suite = Fabric_Client.newCryptoSuite();
		// use the same location for the state store (where the users' certificate are kept)
		// and the crypto store (where the users' keys are kept)
		var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
		crypto_suite.setCryptoKeyStore(crypto_store);
		fabric_client.setCryptoSuite(crypto_suite);
		// get the enrolled user from persistence, this user will sign all requests
		var user_from_store = await fabric_client.getUserContext(username, true);

		if (user_from_store && user_from_store.isEnrolled()) {
			console.log('Successfully loaded user1 from persistence');
			member_user = user_from_store;
		} else {
			throw new Error('Failed to get user1.... run registerUser.js');
			return {"Success":false,"message":"User Not Enrolled to Blockchain"};
		}

		// get a transaction id object based on the current user assigned to fabric client
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);

		// createCar chaincode function - requires 5 args, ex: args: ['CAR12', 'Honda', 'Accord', 'Black', 'Tom'],
		// changeCarOwner chaincode function - requires 2 args , ex: args: ['CAR10', 'Dave'],
		// must send the proposal to endorsing peers

		console.log(channel);

		var request = {
			chaincodeId: chaincodeName,
			fcn: fcn,
			args: args,
			chainId: channelName,
			txId: tx_id
		};

		// send the transaction proposal to the peers
		let results = await channel.sendTransactionProposal(request);

		var proposalResponses = results[0];
		var proposal = results[1];
		let isProposalGood = false;
		if (proposalResponses && proposalResponses[0].response &&
			proposalResponses[0].response.status === 200) {
				isProposalGood = true;
				console.log('Transaction proposal was good');
			} else {
				console.error('Transaction proposal was bad' + results);
			}
		if (isProposalGood) {
			console.log(util.format(
				'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
				proposalResponses[0].response.status, proposalResponses[0].response.message));

			// build up the request for the orderer to have the transaction committed
			var request = {
				proposalResponses: proposalResponses,
				proposal: proposal
			};

			// set the transaction listener and set a timeout of 30 sec
			// if the transaction did not get committed within the timeout period,
			// report a TIMEOUT status
			var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
			var promises = [];

			var sendPromise = channel.sendTransaction(request);
			promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

			// get an eventhub once the fabric client has a user assigned. The user
			// is required bacause the event registration must be signed
			
				var promises = [];
				let event_hubs = channel.getChannelEventHubsForOrg();
				event_hubs.forEach((eh) => {
					logger.debug('invokeEventPromise - setting up event');
					let invokeEventPromise = new Promise((resolve, reject) => {
						let event_timeout = setTimeout(() => {
							let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
							logger.error(message);
							eh.disconnect();
						}, 3000);
						eh.registerTxEvent(tx_id.getTransactionID(), (tx, code, block_num) => {
							logger.info('The chaincode invoke chaincode transaction has been committed on peer %s',eh.getPeerAddr());
							logger.info('Transaction %s has status of %s in blocl %s', tx, code, block_num);
							clearTimeout(event_timeout);

							if (code !== 'VALID') {
								let message = util.format('The invoke chaincode transaction was invalid, code:%s',code);
								logger.error(message);
								reject(new Error(message));
							} else {
								let message = 'The invoke chaincode transaction was valid.';
								logger.info(message);
								resolve(message);
							}
						}, (err) => {
							clearTimeout(event_timeout);
							logger.error(err);
							reject(err);
						},
							// the default for 'unregister' is true for transaction listeners
							// so no real need to set here, however for 'disconnect'
							// the default is false as most event hubs are long running
							// in this use case we are using it only once
							{unregister: true, disconnect: true}
						);
						eh.connect();
					});
					promises.push(invokeEventPromise);
				});

			results = Promise.all(promises);
		} else {
			console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
			return {"Success":false,"message":'Failed to send Proposal or receive valid response. Response null or status is not 200.'};
		}

		// Send Success Message
		return {"Success":true,"transactionId":tx_id.getTransactionID()};

	} catch (error) {
		logger.error('Failed to invoke due to error: ' + error.stack ? error.stack : error);
		return {"Success":false,"message":error.toString()};
	}
};

exports.invokeChaincode = invokeChaincode;
