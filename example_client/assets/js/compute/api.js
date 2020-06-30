/*jshint esversion: 6 */

/*
Kaguya - Web Client

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

const API_VERSION = 1;

function kaguyaAPI(API_HOST) {
    this.identity = "";
    this.client = new WebSocket(API_HOST);
    this.client.onclose = () => console.log("Closed");
}

kaguyaAPI.prototype = {
    _responseFactory: function (
        actionType, action, data
    ) {
        let actionID = this._uuid();
        return {
            id: actionID,
            data: JSON.stringify({
                version: API_VERSION,
                actionID: actionID,
                authToken: this.identity,
                actionType: actionType,
                action: action,
                data: data ? data : {}
            })
        }
    },

    _uuid: function () {
        let d = Date.now();
        if (typeof performance !== 'undefined' && typeof performance.now === 'function') {
            d += performance.now();
        }
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            let r = (d + Math.random() * 16) % 16 | 0;
            d = Math.floor(d / 16);
            return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
        });
    },

    setOnMessageHandle: function (func) {
        this.client.onmessage = func;
    },

    getAccess: function (userId, userPw) {
        let apiStmt = this._responseFactory(
            "authService",
            "getAccess", {
                identity: userId,
                password: userPw
            }
        );
        this.client.send(apiStmt.data);
        return apiStmt.id;
    },

    sendTextMessage: function (targetType, target, message) {
        let apiStmt = this._responseFactory(
            "talkService",
            "sendMessage", {
                contentType: 0,
                targetType: targetType,
                target: target,
                content: message
            }
        );
        this.client.send(apiStmt.data);
        return apiStmt.id;
    }
}