/*
*
* Simple page rendering and DOM/view changes depend on websocket processing
*
* */
import {StorageGet, storageKeyUserUUID} from "./storage";

let viewNodePageSignIn = null;
let viewNodePageChat = null;

let viewNodeSys = null;
let viewNodeError = null;
let viewNodeUsers = null;
let viewNodeMessages = null;

const nodeIdPageSignIn = "page-signIn"
const nodeIdPageChat = "page-chat"

const nodeIdSys = "view-sys";
const nodeIdUsers = "view-users";
const nodeIdErrors = "view-errors";
const nodeIdMessages = "view-messages";

const node = (id) => {return document.getElementById(id);}

export function viewPagesInitialize() {
    if(viewNodePageSignIn == null){
        viewNodePageSignIn = node(nodeIdPageSignIn);
    }
    if(viewNodePageChat == null){
        viewNodePageChat = node(nodeIdPageChat);
    }
}

export function viewShowPageSignIn() {
    viewPagesInitialize();
    viewNodePageChat.style.display = "none";
    viewNodePageSignIn.style.display = "block";
}




export function viewSysAdd(message) {
    if(viewNodeSys == null) {
        viewNodeSys =node(nodeIdSys);
    }

    viewNodeSys.innerHTML = '<div class="log-sys">'+message+'</div>' + viewNodeSys.innerHTML;
}

export function viewErrorAdd(message) {
    if(viewNodeError == null){
        viewNodeError = node(nodeIdErrors);
    }
}




export function buildNodeMessage(message) {
    let divSender = "";
    if(typeof message.Sender != "undefined") {
        divSender = '<div class="message-' + message.UUID + '-sender">' + message.Sender.Username + '</div>';
    }
    let divRecipient = "";
    if(typeof message.Recipient != "undefined") {
        divRecipient = '<div class="message-' + message.UUID + '-recipient">' + message.Recipient.Username + '</div>';
    }
    if(divSender === "" && divRecipient === "") {
        return false;
    }
    let div = document.createElement("div");
    div.id = "message-" + message.UUID;
    div.innerHTML = divSender + divRecipient +
        '<div class="message-' + message.UUID + '-message">' + message.Message + '</div>' +
        '<div class="message-' + message.UUID + '-datetime">' + message.CreatedAt + '</div>' ;

    return div;
}

/*
* @params message Object
* */
export function viewMessagesAdd(message) {
    if(viewNodeMessages == null){
        viewNodeMessages = node(nodeIdMessages);
    }

    const messageNode = buildNodeMessage(message);

    if(messageNode === false) {
        return;
    }
    const nodes = viewNodeMessages.childNodes;
    if(nodes > 0){
        viewNodeMessages.insertBefore(nodes[0], messageNode);
    } else {
        viewNodeMessages.append(messageNode);
    }
}
