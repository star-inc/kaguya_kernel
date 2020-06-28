/*jshint esversion: 6 */

/*
Kaguya - Web Client

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

import Loader from './loader.js';
import Login from './login.js';

$(".disabled").removeClass("disabled");
const Initialize = {
    template: '<div><div v-if="!ready" style="width: 10px;margin:0 auto;"><Loader></Loader></div><div  v-else><Login></Login></div></div>',
    computed: {},
    components: {
        "Loader": Loader,
        "Login": Login
    },
    data() {
        return {
            ready: false
        }
    }
};

export default Initialize;