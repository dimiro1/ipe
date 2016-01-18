
var assert = chai.assert;

var APP_KEY = "7324d55a5eeb8f554761";
var HOST = "localhost";
var PORT = 8080;

Pusher.log = function(msg) {
	console.log(msg);
};

function getPusher() {
  return new Pusher(APP_KEY, {
    wsHost: HOST,
    wsPort: PORT,
    enabledTransports: ["ws"],
    disabledTransports: ["flash"]
  });
}

describe("Pusher", function() {

  describe("connection", function() {
    it("should connect sucessfully", function(done) {
      var pusher = getPusher();

    	pusher.connection.bind('connected', function() {
    		assert.ok(true, "Connected");
    		done();
    	});
    });

    it("should not connect", function(done) {
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
    it("should connect to a channel", function(done) {
      var pusher = getPusher();

      var channel = pusher.subscribe('public-channel');
    	channel.bind("pusher:subscription_succeeded", function(data) {
    		assert.ok(true, "Connected to the channel");
    		done();
    	});
    });
  }); // subscription

  describe("events", function() {
    it('should not allowed client events on public channels', function(done) {
      var pusher = getPusher();
      var channel = pusher.subscribe('public-channel');

    	channel.bind("pusher:subscription_succeeded", function(data) {
    		channel.trigger("client-message", "The message");
    	});

    	pusher.bind("pusher:error", function(data) {
    		assert.ok(true, "Expected error");
    		done();
    	});
    });
  }); // events
});
