"use strict";
var page = require('webpage').create(),
    system = require('system'),
    address;

if (system.args.length === 1) {
    console.log('Usage: netlog.js <some URL>');
    phantom.exit(1);
}

address = system.args[1];
phantom.outputEncoding = '';
page.settings.javascriptEnabled = true;
page.settings.loadImages = false;
page.onResourceRequested = function (req) {
    console.log('requested: ' + JSON.stringify(req, undefined, 4));
};

page.onResourceReceived = function (res) {
    console.log('received: ' + JSON.stringify(res, undefined, 4));
};

page.open(address, function (status) {
    if (status !== 'success') {
        console.log('FAIL to load the address');
    }else{
        console.log(page.content);
    }
    phantom.exit();
});


