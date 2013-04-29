function togglePause(){
	$.ajax({
		dataType: "json",
		url : "/togglePause.html",
		success: togglePauseCallback
	});
}
var togglePauseText="播放";
function togglePauseCallback(data){
	var oldTogglePauseText=$("#togglePause").html()
	if (data.success){
		$("#togglePause").html(togglePauseText);
		togglePauseText=oldTogglePauseText;
	}
}

function next(){
	$.ajax({
		dataType: "json",
		url : "/next.html"
	});
}

function loadAlbumInfo(){
	var wsuri ="ws://"+location.host+"/song.html"
	sock = new WebSocket(wsuri);

	sock.onopen = function() {
		console.log("connected to " + wsuri);
	}

	sock.onclose = function(e) {
		console.log("connection closed (" + e.code + ")");
	}

	sock.onmessage = function(e) {
		console.log("message received: " + e.data);
		loadAlbumInfoCallback(eval(e.data));
	}
}	
function loadAlbumInfoCallback(data){
	$("#songTitle").html(data.title);
	$("#songPic").attr("src",data.pic);
	$("#artist").html(data.artist);
	$("#summary").html(data.summary);
}

