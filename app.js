var stream, recording = false, encoder, ws, input, node;


function success(localMediaStream) {
	recording = true;
	$('#recording_sign').show();
	$('#start_btn').attr('disabled', 'disabled');
	$('#stop_btn').removeAttr('disabled');

	console.log('success grabbing microphone');
	stream = localMediaStream;

	var audio_context = new window.webkitAudioContext();

	input = audio_context.createMediaStreamSource(stream);
	node = input.context.createJavaScriptNode(4096, 1, 1);

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
	console.log('grabbing microphone failed: ' + code);
}

$(document).ready(function() {

	$('#start_btn').click(function(e) {
		console.log('pressed start button');
		e.preventDefault();

		encoder = new Worker('encoder.js');
		encoder.postMessage({ cmd: 'init', config: { samplerate: 22050, bitrate: 32 } });

		encoder.onmessage = function(e) {
			ws.send(e.data.buf);
			if (e.data.cmd == 'end') {
				ws.close();
				ws = null;
				encoder.terminate();
				encoder = null;
			}
		};

		var ws = new WebSocket("ws://" + window.location.host + "/ws/audio");
		ws.onopen = function() {
			navigator.webkitGetUserMedia({ vidoe: false, audio: true }, success, fail);
		};
	});

	$('#stop_btn').click(function(e) {
		console.log('pressed stop button');
		e.preventDefault();
		stream.stop();
		recording = false;
		encoder.postMessage({ cmd: 'finish' });

		$('#recording_sign').hide();
		$('#start_btn').removeAttr('disabled');
		$('#stop_btn').attr('disabled', 'disabled');

		input.disconnect();
		node.disconnect();
		input = node = null;
	});
});
