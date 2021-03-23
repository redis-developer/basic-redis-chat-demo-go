/*
* Websocket onmessage event handlers
*
* ws.onmessage event has one argument `event`,
* we should read `event.data` value, it contained message from backend
*
* `event.data` is stringify JSON object, each of received JSON object contained `type:"<Message Type>"` property
*
* `messages` - is a key:value map, where:
*   - key: received message type from websocket, message type contained in `event.data.type`
*   - value: handler function for processing received message
* */

import {storageKeySessionUUID, storageKeyUserAccessKey, storageKeyUserUUID, StorageSet} from "./storage";
import {viewMessagesAdd, viewShowPageChat, viewUsersAdd, viewUsersClean} from "./view";
import useAppStateContext from '../state';

export const messages = {
    // system message from backend, usually it contained system info for debug
    "sys": sys,
    // error from backend when backend received message and can't processed it
    "error": error,
    // backend return session UUID when websocket connection successful
    "ready": ready,
    // backend return user data on signIn successful
    "authorized": authorized,
    // backend return all users
    "users": users,
    // backend said that somebody joined to channel
    "channelJoin": channelJoin,
    // backend said that channel accepted new message
    "channelMessage": channelMessage
}

const messagesSys = {
    "signIn": sysSignIn(),
}

function sys(data, {dispatch}) {
    if(typeof messagesSys[data.sys.type] == "function") {
        messagesSys[data.sys.type](data.sys);
    } else {
        console.log("Unknown message sys.type: " + data.sys.type)
    }
    if (data.sys.signIn)
        dispatch({type: 'set user', payload: data.sys.signIn});
}

function error(data) {
    console.log("error: " + data.error.code + " - " + data.error.message, data)
}

function ready(data) {
    StorageSet(storageKeySessionUUID, data.ready.sessionUUID);
}

function authorized(data) {
    StorageSet(storageKeyUserUUID,data.authorized.userUUID);
    StorageSet(storageKeyUserAccessKey, data.authorized.accessKey);
}


function users(data, {dispatch}) {
    dispatch({type: 'set users', payload: data.users.users})
    dispatch({type: 'set rooms', payload: data.users.users})
}

function channelJoin(data, {dispatch}) {
    dispatch({type: 'set messages', payload: data.channelJoin.messages || []});
    dispatch({type: 'set users', payload: data.channelJoin.users});
}

function channelMessage(data, {dispatch}) {
    console.log('append message', data);
    dispatch({type: 'append message', payload: {id: data.channelMessage.RecipientUUID, message: data.channelMessage}});
}

function sysSignIn(data){
    console.log('sysSignIn')
}
