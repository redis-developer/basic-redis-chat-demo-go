/*
*
* Data storage for store data (session uuid, user uuid, access key, etc.) and reuse it, as example: localStorage
*
* */
const storage = window.localStorage;

export const storageKeySessionUUID = "session.uuid";

export const storageKeyUserUUID = "user.uuid";
export const storageKeyUserAccessKey = "user.accessKey";

export function StorageSet(key, value) {
    storage.setItem(key, value);
}

export function StorageGet(key) {
    const value =  storage.getItem(key);
    return value == null?"":value;
}
