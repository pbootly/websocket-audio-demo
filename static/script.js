const ws = new WebSocket('ws://localhost:8080/ws');

// Audio player
const audioPlayer = document.getElementById('audioPlayer');
const audioSource = new MediaSource();
audioPlayer.src = URL.createObjectURL(audioSource);
const audioBuffer = [];

// Bitrate calculation
let totalBytesReceived = 0;
let startTime = Date.now();

ws.onmessage = (event) => {
    // Handle incoming audio data
    const blob = event.data;
    audioBuffer.push(blob);

    // Update the bitrate
    totalBytesReceived += blob.size;
    const elapsedSeconds = (Date.now() - startTime) / 1000;
    const bitrate = (totalBytesReceived * 8) / (elapsedSeconds * 1000);
    document.getElementById('bitrate').innerHTML = `Bitrate: ${bitrate.toFixed(2)} kbps`;

    if (audioPlayer.paused) {
        // Start playing audio if not already playing
        playAudio();
    }
};

const playAudio = () => {
    if (audioBuffer.length > 0 && audioSource.readyState === 'open') {
        const audioData = audioBuffer.shift();
        const sourceBuffer = audioSource.addSourceBuffer('audio/mp4');
        const reader = new FileReader();

        reader.onload = (e) => {
            sourceBuffer.appendBuffer(new Uint8Array(e.target.result));
            sourceBuffer.addEventListener('updateend', () => {
                audioPlayer.play();
                playAudio();
            });
        };

        reader.readAsArrayBuffer(audioData);
    }
};
