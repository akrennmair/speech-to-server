importScripts('libmp3lame.js');

var mp3codec;

self.onmessage = function(e) {
	switch (e.data.cmd) {
	case 'init':
		mp3codec = Lame.init();
		Lame.set_mode(mp3codec, Lame.MONO);
		Lame.set_num_channels(mp3codec, 1);
		Lame.set_out_samplerate(mp3codec, 22050);
		Lame.set_bitrate(mp3codec, 32);
		Lame.init_params(mp3codec);
		break;
	case 'encode':
		var mp3data = Lame.encode_buffer_ieee_float(mp3codec, e.data.buf, e.data.buf);
		self.postMessage({cmd: 'data', buf: mp3data.data});
		break;
	case 'finish':
		var mp3data = Lame.encode_flush(mp3codec);
		self.postMessage({cmd: 'end', buf: mp3data.data});
		Lame.close(mp3codec);
		mp3codec = null;
		break;
	}
};
