/* NODE-RED-CONTRIB-SNAP4CITY-DEVELOPER
   Copyright (C) 2018 DISIT Lab http://www.disit.org - University of Florence

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. */
module.exports = function (RED) {
    function eventLog(inPayload, outPayload, config, _agent, _motivation, _ipext, _modcom) {
        var os = require('os');
        var ifaces = os.networkInterfaces();
        var uri = "http://192.168.1.43/RsyslogAPI/rsyslog.php";

        var pidlocal = RED.settings.APPID;
        var iplocal = null;
        Object.keys(ifaces).forEach(function (ifname) {
            ifaces[ifname].forEach(function (iface) {
                if ('IPv4' !== iface.family || iface.internal !== false) {
                    // skip over internal (i.e. 127.0.0.1) and non-ipv4 addresses
                    return;
                }
                iplocal = iface.address;
            });
        });
        iplocal = iplocal + ":" + RED.settings.uiPort;
        var timestamp = new Date().getTime();
        var modcom = _modcom;
        var ipext = _ipext;
        var payloadsize = JSON.stringify(outPayload).length / 1000;
        var agent = _agent;
        var motivation = _motivation;
        var lang = (inPayload.lang ? inPayload.lang : config.lang);
        var lat = (inPayload.lat ? inPayload.lat : config.lat);
        var lon = (inPayload.lon ? inPayload.lon : config.lon);
        var serviceuri = (inPayload.serviceuri ? inPayload.serviceuri : config.serviceuri);
        var message = (inPayload.message ? inPayload.message : config.message);
        var XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;
        var xmlHttp = new XMLHttpRequest();
        console.log(encodeURI(uri + "?p=log" + "&pid=" + pidlocal + "&tmstmp=" + timestamp + "&modCom=" + modcom + "&IP_local=" + iplocal + "&IP_ext=" + ipext +
            "&payloadSize=" + payloadsize + "&agent=" + agent + "&motivation=" + motivation + "&lang=" + lang + "&lat=" + (typeof lat != "undefined" ? lat : 0.0) + "&lon=" + (typeof lon != "undefined" ? lon : 0.0) + "&serviceUri=" + serviceuri + "&message=" + message));
        xmlHttp.open("GET", encodeURI(uri + "?p=log" + "&pid=" + pidlocal + "&tmstmp=" + timestamp + "&modCom=" + modcom + "&IP_local=" + iplocal + "&IP_ext=" + ipext +
            "&payloadSize=" + payloadsize + "&agent=" + agent + "&motivation=" + motivation + "&lang=" + lang + "&lat=" + (typeof lat != "undefined" ? lat : 0.0) + "&lon=" + (typeof lon != "undefined" ? lon : 0.0) + "&serviceUri=" + serviceuri + "&message=" + message), true); // false for synchronous request
        xmlHttp.send(null);
    }

    function NotificatorLastEvents(config) {
        RED.nodes.createNode(this, config);
        var node = this;
        var uri = "http://notificator.km4city.org/notificator/restInterfaceExternal.php?operation=getEvents";
        var dashboard = config.dashboard;
        var widget = config.widget;
        var event = config.event;
        var checkevery = config.checkevery;
        var uid = RED.settings.APPID;
        var inPayload = {};
        var msg = {};
        var XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;
        var xmlHttp = new XMLHttpRequest();
        console.log(node.interval);
        if (node.interval != null) {
            console.log("Cancello Intervallo");
            clearInterval(node.interval);
        }
        node.interval = setInterval(function () {
            console.log(encodeURI(uri + "&startDate=" + (new Date(Date.now() - new Date().getTimezoneOffset() * 1000 * 60 - checkevery * 1000)).toISOString().split('.')[0].replace("T", " ") + "&dashboardTitle=" + dashboard + "&widgetTitle=" + widget + "&appID=iotapp"));
            xmlHttp.open("GET", encodeURI(uri + "&startDate=" + (new Date(Date.now() - new Date().getTimezoneOffset() * 1000 * 60 - checkevery * 1000)).toISOString().split('.')[0].replace("T", " ") + "&dashboardTitle=" + dashboard + "&widgetTitle=" + widget + "&appID=iotapp"), false); // false for synchronous request
            xmlHttp.send(null);
            if (xmlHttp.responseText != "") {
                console.log(xmlHttp.responseText);
                msg.payload = JSON.parse(xmlHttp.responseText).data;
            } else {
                msg.payload = JSON.parse("{\"status\": \"error\"}");
            }
            eventLog(inPayload, msg, config, "Node-Red", "Notificator", uri, "RX");
            node.send(msg);
        }, checkevery * 1000);

        node.nodeClosingDone = function () {
            console.log(node.getNow() + " - notificator-last-events node " + node.name + " has been closed");
        };

        node.on('close', function (removed, nodeClosingDone) {
            if (removed) {
                // Cancellazione nodo
                console.log(node.getNow() + " - notificator-last-events node " + node.name + " is being removed from flow");
                console.log(node.interval);
                if (node.interval != null) {
                    console.log("Cancello Intervallo");
                    clearInterval(node.interval);
                }
            } else {
                // Riavvio nodo
                console.log(node.getNow() + " - notificator-last-events node " + node.name + " is being rebooted");
                console.log(node.interval);
                if (node.interval != null) {
                    console.log("Cancello Intervallo");
                    clearInterval(node.interval);
                }
            }
            nodeClosingDone();

        });
    }
    RED.nodes.registerType("notificator-last-events", NotificatorLastEvents);
}