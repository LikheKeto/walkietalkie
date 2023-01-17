<script lang="ts">
	import { onMount } from 'svelte';

	let blinker = false;
	let value = 0;

	$: value = parseFloat(value.toFixed(2));
	setInterval(() => {
		blinker = !blinker;
	}, 1000);

	$: roundedValue = Math.round((value / 1000) * 270);

	let userStream: MediaStream, peerStream: MediaStream;
	let ws: WebSocket;
	let peerRef: RTCPeerConnection;

	onMount(async () => {

		navigator.mediaDevices.getUserMedia({ audio: true }).then((mediaStream) => {

			userStream = mediaStream;
			peerStream = userStream.clone();
			peerStream.removeTrack(peerStream.getTracks()[0]);
			document.querySelector('audio').srcObject = peerStream;
		});
	});

	let once = false;

	let prevTimeoutID: NodeJS.Timeout;
	$: {
		if (value > 0) {
			clearTimeout(prevTimeoutID);
			prevTimeoutID = setTimeout(() => {
				if (!once) {
					once = true;
				} else {
					return;
				}

				// send request to backend
				ws = new WebSocket(`ws://35.79.227.8:8080/connect?freq=${value}`);
				ws.addEventListener('open', () => {
					ws.send(JSON.stringify({ join: 'true' }));
				});
				ws.addEventListener('message', async (msg) => {
					const message = JSON.parse(msg.data);
					if (message.join) {
						joinFreq();
					}
					if (message.iceCandidate) {
						try {
							await peerRef.addIceCandidate(message.iceCandidate);
						} catch (err) {
							console.log(err);
						}
					}
					if (message.offer) {
						await handleOffer(message.offer);
					}
					if (message.answer) {
						peerRef.setRemoteDescription(new RTCSessionDescription(message.answer));
					}
				});
			}, 2000);
		}
	}

	const handleOffer = async (offer: any) => {
		peerRef = createPeer();

		await peerRef.setRemoteDescription(new RTCSessionDescription(offer));

		userStream.getTracks().forEach((track) => {
			peerRef.addTrack(track, userStream);
		});

		const answer = await peerRef.createAnswer();
		await peerRef.setLocalDescription(answer);

		ws.send(JSON.stringify({ answer: peerRef.localDescription }));
	};

	const joinFreq = () => {
		peerRef = createPeer();

		userStream.getTracks().forEach((track) => {
			peerRef.addTrack(track, userStream);
		});
	};

	const createPeer = () => {
		const peer = new RTCPeerConnection({
			iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
		});

		peer.onnegotiationneeded = async (e: Event) => {
			try {
				const myOffer = await peerRef.createOffer();
				await peerRef.setLocalDescription(myOffer);
				ws.send(JSON.stringify({ offer: myOffer }));
			} catch (err) {
				console.log(err);
			}
		};

		peer.onicecandidate = (e: RTCPeerConnectionIceEvent) => {
			if (e.candidate) {
				ws.send(JSON.stringify({ iceCandidate: e.candidate }));
			}
		};
		peer.ontrack = (e: RTCTrackEvent) => {
			peerStream.addTrack(e.track);
		};

		return peer;
	};

	let dragging = false;
	let startY: number;

	const dragstart = (e: any) => {
		dragging = true;
		startY = e.screenY;
	};
	const dragend = (e: any) => {
		dragging = false;
	};

	const mousemoving = (e: any) => {
		if (dragging) {
			let endY = e.screenY;
			let dy = startY - endY;

			value -= dy * 2;
			startY = endY;
			value = parseFloat(value.toFixed(2));
			if (value > 1000) value = 1000;
			if (value < 0) value = 0;
		}
	};
</script>

<div class="w-screen h-screen flex justify-center items-center bg-purple-900">
	<div
		class="w-3 h-16 bg-black absolute bottom-[85%] right-[%] border-t-[1rem] {blinker
			? 'border-t-red-500'
			: 'border-t-black'} rounded-t-md"
	/>
	<div class="bg-purple-700 grid grid-rows-6 w-[90%] md:w-[50%] max-w-[25rem] h-[70%] shadow-md">
		<input
			max="1000"
			min="0"
			bind:value
			class="row-span-2 px-4 text-2xl outline-none min-w-0"
			type="number"
		/>
		<div
			on:mousedown={dragstart}
			on:mouseleave={() => (dragging = false)}
			on:mousemove={mousemoving}
			on:mouseup={dragend}
			class="bg-purple-400 cursor-pointer relative row-start-4 row-span-2 m-auto h-48 w-48 max-w-[15rem] shadow-md flex items-center justify-center rounded-[50%]"
		>
			<p class="z-10 text-2xl select-none text-white font-bold">
				{value} Hz
			</p>
			<div
				style="transform: rotate({roundedValue}deg);"
				class="absolute w-2 bg-red-500 origin-bottom h-[50%] border-b-[2rem] border-b-purple-400 top-0"
			/>
		</div>
	</div>
	<audio autoplay src="" />
</div>
