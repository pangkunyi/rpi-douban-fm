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
	$.ajax({
		dataType: "json",
		url : "/song.html",
		success: loadAlbumInfoCallback
	});
}	
function loadAlbumInfoCallback(data){
	$("#songTitle").html(data.title);
	$("#songPic").html(data.pic);
	$("#artist").html(data.artist);
	$("#summary").html(data.summary);
	loadAlbumInfo();
}

