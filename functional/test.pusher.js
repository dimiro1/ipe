"use strict";

let assert = chai.assert;

let APP_KEY = "278d525bdf162c739803";
let HOST = "localhost";
let PORT = 8080;
let AUTH = "http://localhost:5000/pusher/private/auth";
let AUTH_PRESENCE = "http://localhost:5000/pusher/presence/auth";

Pusher.log = function (msg) {
    console.log(msg);
};

function getPusher(auth) {
    return new Pusher(APP_KEY, {
        wsHost: HOST,
        wsPort: PORT,
        authEndpoint: auth,
        enabledTransports: ["ws"],
        disabledTransports: ["flash"],
    });
}

describe("Pusher", function () {

    describe("connection", function () {
        it("should connect sucessfully with correct config", function (done) {
            let pusher = getPusher(AUTH);

            pusher.connection.bind('connected', function () {
                assert.ok(true, "Connected");
                done();
            });
        });

        it("should not connect without the correct config", function (done) {
            let pusher = new Pusher("INVALID_APP_KEY", {
                wsHost: HOST,
                wsPort: PORT,
                enabledTransports: ["ws"],
                disabledTransports: ["flash"]
            });

            pusher.connection.bind('disconnected', function () {
                assert.ok(true, "Not Connected");
                done();
            });
        });
    }); // connection

    describe("subscription", function () {
        it("should subscribe to a public channel", function (done) {
            let pusher = getPusher(AUTH);

            let channel = pusher.subscribe('public-channel');
            channel.bind("pusher:subscription_succeeded", function () {
                assert.ok(true, "Connected to the channel");
                done();
            });
        });

        it("should subscribe to a private channel", function (done) {
            let pusher = getPusher(AUTH);

            let channel = pusher.subscribe('private-channel');
            channel.bind("pusher:subscription_succeeded", function () {
                assert.ok(true, "Connected to the channel");
                done();
            });
        });

        it("should subscribe to a presence channel", function (done) {
            let pusher = getPusher(AUTH_PRESENCE);

            let channel = pusher.subscribe('presence-channel');
            channel.bind("pusher:subscription_succeeded", function () {
                assert.ok(true, "Connected to the channel");
                done();
            });
        });
    }); // subscription

    describe("hooks", function () {
        it('should receive hook', function (done) {
            let pusher = getPusher(AUTH);
            let channel = pusher.subscribe('private-webhook');

            channel.bind("pusher:subscription_succeeded", function () {
                console.log("subscribed");
            });

            channel.bind("channel_occupied", function (data) {
                assert.equal(data, "The Webhoook from server");
                pusher.unsubscribe('private-webhook');
                done();
            });
        });
    }); // hooks

    describe("events", function () {
        it('should not allowed client events on public channels', function (done) {
            let pusher = getPusher(AUTH);
            let channel = pusher.subscribe('public-channel');

            channel.bind("pusher:subscription_succeeded", function () {
                channel.trigger("client-message", "The message");
            });

            pusher.bind("pusher:error", function () {
                assert.ok(true, "Expected error");
                done();
            });
        });

        it('should allow client events on private channels', function (done) {
            let pusher_a = getPusher(AUTH);
            let pusher_b = getPusher(AUTH);

            let channel_a = pusher_a.subscribe('private-channel');
            let channel_b = pusher_b.subscribe('private-channel');

            channel_a.bind("pusher:subscription_succeeded", function () {
                channel_a.trigger("client-message", "The message");
            });

            channel_b.bind("client-message", function (data) {
                assert.equal(data, "The message");
                done();
            });
        });

        it('should publish event on private channel', function (done) {
            let pusher_a = getPusher(AUTH);
            let pusher_b = getPusher(AUTH);

            let channel_a = pusher_a.subscribe('private-messages');
            let channel_b = pusher_b.subscribe('private-messages');

            channel_a.bind("pusher:subscription_succeeded", function () {
                console.log("channel_a connected");
                let xhttp = new XMLHttpRequest();
                xhttp.open("GET", "/trigger", true);
                xhttp.send();
            });

            channel_b.bind("pusher:subscription_succeeded", function () {
                console.log("channel_b connected");
            });

            channel_b.bind("messages", function (data) {
                assert.equal(data, "The message from server");
                done();
            });
        });

    }); // events
});
