
var assert = chai.assert;

var APP_KEY = "278d525bdf162c739803";
var HOST = "localhost";
var PORT = 8080;
var AUTH = "http://localhost:5000/pusher/private/auth"
var AUTH_PRESENCE = "http://localhost:5000/pusher/presence/auth"

Pusher.log = function(msg) {
	console.log(msg);
};

function getPusher(auth) {
  return new Pusher(APP_KEY, {
    wsHost: HOST,
    wsPort: PORT,
    authEndpoint: auth,
    enabledTransports: ["ws"],
    disabledTransports: ["flash"],
    cluster: "hello", // Should be ignored
  });
}

describe("Pusher", function() {

  describe("connection", function() {
    it("should connect sucessfully with correct config", function(done) {
      var pusher = getPusher(AUTH);

    	pusher.connection.bind('connected', function() {
    		assert.ok(true, "Connected");
    		done();
    	});
    });

    it("should not connect without the correct config", function(done) {
      var pusher = new Pusher("INVALID_APP_KEY", {
    		wsHost: HOST,
    		wsPort: PORT,
    		enabledTransports: ["ws"],
    		disabledTransports: ["flash"]
    	});

      pusher.connection.bind('disconnected', function() {
    		assert.ok(true, "Not Connected");
    		done();
    	});
    });
  }); // connection

  describe("subscription", function() {
    it("should subscribe to a public channel", function(done) {
      var pusher = getPusher(AUTH);

      var channel = pusher.subscribe('public-channel');
    	channel.bind("pusher:subscription_succeeded", function(data) {
    		assert.ok(true, "Connected to the channel");
    		done();
    	});
    });

    it("should subscribe to a private channel", function(done) {
      var pusher = getPusher(AUTH);

      var channel = pusher.subscribe('private-channel');
      channel.bind("pusher:subscription_succeeded", function(data) {
        assert.ok(true, "Connected to the channel");
        done();
      });
    });

    it("should subscribe to a presence channel", function(done) {
      var pusher = getPusher(AUTH_PRESENCE);

      var channel = pusher.subscribe('presence-channel');
      channel.bind("pusher:subscription_succeeded", function(data) {
        assert.ok(true, "Connected to the channel");
        done();
      });
    });
  }); // subscription

  describe("events", function() {
    it('should not allowed client events on public channels', function(done) {
      var pusher = getPusher(AUTH);
      var channel = pusher.subscribe('public-channel');

    	channel.bind("pusher:subscription_succeeded", function(data) {
    		channel.trigger("client-message", "The message");
    	});

    	pusher.bind("pusher:error", function(data) {
    		assert.ok(true, "Expected error");
    		done();
    	});
    });

    it('should allow client events on private channels', function(done) {
      var pusher_a = getPusher(AUTH);
      var pusher_b = getPusher(AUTH);

      var channel_a = pusher_a.subscribe('private-channel');
      var channel_b = pusher_b.subscribe('private-channel');

      channel_a.bind("pusher:subscription_succeeded", function() {
        channel_a.trigger("client-message", "The message");
      });

      channel_b.bind("client-message", function(data) {
        assert.equal(data, "The message");
        done();
      });
    });

    it('should publish event on private channel', function(done) {
      var pusher_a = getPusher(AUTH);
      var pusher_b = getPusher(AUTH);

      var channel_a = pusher_a.subscribe('private-messages');
      var channel_b = pusher_b.subscribe('private-messages');

      channel_a.bind("pusher:subscription_succeeded", function() {
        var xhttp = new XMLHttpRequest();
        xhttp.open("GET", "/trigger", true);
        xhttp.send();
      });

      channel_b.bind("messages", function(data) {
        assert.equal(data, "The message from server");
        done();
      });
    });

  }); // events
});
