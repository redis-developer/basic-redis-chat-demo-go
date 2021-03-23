/*
* Will return data object for websocket send to backend
*
*   function Data<MessageType>([argument, ...]) {
*       return {
*           <Property>: <Value>,
*           ...
*       }
*   }
*
* */

import {StorageGet, storageKeySessionUUID, storageKeyUserAccessKey, storageKeyUserUUID} from "./storage";

const dataTypeSignIn = "signIn";
const dataTypeSignOut = "signOut";
const dataTypeUsers = "users";
const dataTypeChannelJoin = "channelJoin";
const dataTypeChannelMessage = "channelMessage";
const dataTypeChannelLeave = "channelLeave";

export function DataSignIn(username, password) {
    return {
        SUUID: StorageGet(storageKeySessionUUID),
        type: dataTypeSignIn,
        signIn: {
            username: username,
            password: password
        }
    }
}

export function DataUsers() {
    return {
        SSUID: StorageGet(storageKeySessionUUID),
        type: dataTypeUsers,
        userUUID: StorageGet(storageKeyUserUUID),
        accessKey: StorageGet(storageKeyUserAccessKey),
    }
}

export function DataChannelJoin(recipientUUID) {
    return {
        SUUID: StorageGet(storageKeySessionUUID),
        type: dataTypeChannelJoin,
        userUUID: StorageGet(storageKeyUserUUID),
        accessKey: StorageGet(storageKeyUserAccessKey),
        channelJoin: {
            recipientUUID: recipientUUID
        }
    }
}

export function DataChannelLeave(recipientUUID) {
    return {
        SUUID: StorageGet(storageKeySessionUUID),
        type: dataTypeChannelLeave,
        userUUID: StorageGet(storageKeyUserUUID),
        userAccessKey: StorageGet(storageKeyUserAccessKey),
        channelLeave: {
            recipientUUID: recipientUUID,
            senderUUID: StorageGet(storageKeyUserUUID)
        }
    }
}

export function DataChannelMessage(recipientUUID, message) {
    return {
        SUUID: StorageGet(storageKeySessionUUID),
        type: dataTypeChannelMessage,
        userUUID: StorageGet(storageKeyUserUUID),
        accessKey: StorageGet(storageKeyUserAccessKey),
        channelMessage: {
            recipientUUID: recipientUUID,
            message: message
        }
    }
}
