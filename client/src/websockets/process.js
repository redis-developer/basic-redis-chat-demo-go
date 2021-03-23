/*
*   Basic processing examples
*
*   You can compare websocket messages sending or receiving together in
*   specific functions for make custom processing flow
*
* */

import {webSocketSend} from "../hooks";
import {DataChannelJoin, DataChannelLeave, DataChannelMessage, DataSignIn, DataUsers} from "./data";
import {node} from "prop-types";

let inputSignInUsername = null;
let inputSignInPassword = null;
let inputMessage = null;

const nodeIdInputSignInUsername = "input-username";
const nodeIdInputSignInPassword = "input-password";
const nodeIdInputMessage = "input-message";

let selectedRecipientUUID = "";

// signIn flow
export function processSignIn(username,password,setShowLogin) {
    let process = new Promise((resolve, reject) => {
        webSocketSend(DataSignIn(username, password));
        resolve();
    })

    process.then(() => {webSocketSend(DataUsers())});
    process.then(() => {webSocketSend(DataChannelJoin(""));});
    process.catch((err) => {console.log(err)});

    setShowLogin(false)
    return ;
}

export function processChannelJoin(recipientUUID) {
    let process = new Promise((resolve, reject) => {
        webSocketSend(DataChannelJoin(recipientUUID));
        resolve();
    });
    process.then(() => {selectedRecipientUUID = recipientUUID;})
    process.catch((err) => console.log(err));
}

export function processChannelLeave() {
    let process = new Promise((resolve, reject) => {
        webSocketSend(DataChannelLeave(selectedRecipientUUID));
        resolve();
    });
    process.catch((err) => console.log('error', err));
}

export function processChannelMessage(inputMessage) {
    if(inputMessage == null) {
        inputMessage = node(nodeIdInputMessage);
    }
    console.log('channelMessage', inputMessage)
    let process = new Promise((resolve, reject) => webSocketSend(DataChannelMessage(selectedRecipientUUID, inputMessage)));
    process.then(() => {});
    process.catch((err) => console.log(err));
    return false;
}
