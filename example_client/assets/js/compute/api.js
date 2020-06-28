/*jshint esversion: 6 */

/*
Kaguya - Web Client

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

const API_HOST = "wss://localhost";

function kaguyaAPI() {
    this.client = new WebSocket(API_HOST);
    this.client.on_message = function (message) {
        if ("action" in message) {
            this.prototype["r_" + message["action"]];
        }
    }
    this.client.on_close = () => console.log("Closed");
}

kaguyaAPI.prototype = {
    getAccess: function (userId, userPw) {
        this.client.send(JSON.stringify({
            identity: userId,
            password: userPw
        }));
    },

    r_ReceivedMessage: function (){

    }
}