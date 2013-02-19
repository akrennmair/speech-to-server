importScripts('libmp3lame.js');

var mp3codec = Lame.init();
Lame.init_params(mp3codec);

self.onmessage = function(e) {
	switch (e.data.cmd) {
	case 'encode':
		var mp3data = Lame.encode_buffer_ieee_float(mp3codec, e.data.buf, e.data.buf);
		self.postMessage({buf: mp3data.data});
		break;
	case 'finish':
		var mp3data = Lame.encode_flush(mp3codec);
		self.postMessage({buf: mp3data.data});
		Lame.close(mp3codec);
		break;
	}
};
