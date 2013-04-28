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
		success: loadAlbumInfoCallback,
		complete: loadAlbumInfo,
		timeout: 9000000
	});
}	
function loadAlbumInfoCallback(data){
	if($("#songPic").attr("src") != data.pic){
		$("#songTitle").html(data.title);
		$("#songPic").attr("src",data.pic);
		$("#artist").html(data.artist);
		$("#summary").html(data.summary);
	}
}

