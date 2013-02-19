var stream = null;
var audio_context = new window.webkitAudioContext();
var recording = false;

var encoder = new Worker('encoder.js');

var ws = new WebSocket("ws://" + window.location.host + "/ws/audio");
ws.onopen = function() {
	console.log('ws.onopen called');
};

encoder.onmessage = function(e) {
	ws.send(e.data.buf);
	if (e.data.cmd == 'end') {
		ws.close();
	}
};

function success(localMediaStream) {
	stream = localMediaStream;
	console.log('success');

	var input = audio_context.createMediaStreamSource(stream);

	//input.connect(audio_context.destination);

	var node = input.context.createJavaScriptNode(4096, 1, 1);
	console.log('sampleRate: ' + input.context.sampleRate);

	node.onaudioprocess = function(e) {
		if (!recording)
			return;
		var channelLeft = e.inputBuffer.getChannelData(0);
		encoder.postMessage({ cmd: 'encode', buf: channelLeft });
	};

	input.connect(node);
	node.connect(audio_context.destination);
}

function fail(code) {
	console.log('fail: ' + code);
}

$(document).ready(function() {

	$('#start_btn').click(function(e) {
		console.log('start');
		e.preventDefault();
		recording = true;
		navigator.webkitGetUserMedia({ vidoe: false, audio: true }, success, fail);
	});

	$('#stop_btn').click(function(e) {
		e.preventDefault();
		stream.stop();
		recording = false;
		encoder.postMessage({ cmd: 'finish' });
	});
});
