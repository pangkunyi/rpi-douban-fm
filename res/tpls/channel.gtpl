<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black" />
        <title>
        </title>
        <link rel="stylesheet" href="https://s3.amazonaws.com/codiqa-cdn/mobile/1.2.0/jquery.mobile-1.2.0.min.css" />
        <link rel="stylesheet" href="/static/my.css" />
        <script src="https://s3.amazonaws.com/codiqa-cdn/jquery-1.7.2.min.js">
        </script>
        <script src="https://s3.amazonaws.com/codiqa-cdn/mobile/1.2.0/jquery.mobile-1.2.0.min.js">
        </script>
        <script src="/static/my.js">
        </script>
    </head>
    <body>
        <!-- Home -->
	<div data-role="page" id="playPage" data-cid="playPage" class="ui-page ui-body-a ui-page-footer-fixed" data-theme="a">
		<div data-theme="a" data-role="header" data-cid="pageheader1" class="ui-sortable ui-header ui-bar-a ui-sortable-disabled" style="" role="banner">
			<h3 data-cid="heading1" class="ui-title" role="heading" aria-level="1"><span id="songTitle">{{.Title}}</span>/<span id="artist">{{.Artist}}</span></h3>
		</div>
		<div data-role="content" style="padding: 0px;" data-cid="pagecontent1" class="ui-sortable ui-content ui-sortable-disabled" role="main">
	<div style="" data-cid="image1">
		<img style="width: 320px; height: 320px" id="songPic" src="{{.Picture}}">
	</div>
	<div data-cid="text1">
		<pre id="summary">{{.AlbumInfo.Summary}} </pre>	
	</div>

		<div data-role="tabbar" data-iconpos="top" data-theme="a" data-cid="tabbar1" class="ui-footer ui-footer-fixed ui-bar-a ui-navbar ui-mini" role="navigation">
			<ul class="ui-grid-c">
				<li class="ui-block-a">
					<a href="/index.html" data-transition="fade" data-theme="a" data-icon="arrow-l" class="ui-btn ui-btn-inline ui-btn-icon-top ui-btn-up-a" data-corners="false" data-shadow="false" data-iconshadow="true" data-iconsize="18" data-wrapperels="span" data-iconpos="top" data-inline="true">
						<span class="ui-btn-inner">
							<span class="ui-btn-text">返回</span>
							<span class="ui-icon ui-icon-back ui-icon-shadow ui-iconsize-18">&nbsp;</span>
						</span>
					</a>
				</li>

				<li class="ui-block-b">
					<a href="#" data-transition="fade" data-theme="a" data-icon="star" class="ui-btn ui-btn-inline ui-btn-icon-top ui-btn-up-a" data-corners="false" data-shadow="false" data-iconshadow="true" data-iconsize="18" data-wrapperels="span" data-iconpos="top" data-inline="true">
						<span class="ui-btn-inner">
							<span class="ui-btn-text">红心</span>
							<span class="ui-icon ui-icon-star ui-icon-shadow ui-iconsize-18">&nbsp;</span>
						</span>
					</a>
				</li>
				<li class="ui-block-c">
					<a href="#" onclick="javascript:togglePause()" data-transition="fade" data-theme="a" data-icon="star" class="ui-btn ui-btn-inline ui-btn-icon-top ui-btn-up-a" data-corners="false" data-shadow="false" data-iconshadow="true" data-iconsize="18" data-wrapperels="span" data-iconpos="top" data-inline="true">
						<span class="ui-btn-inner">
							<span id="togglePause" class="ui-btn-text">暂停</span>
							<span class="ui-icon ui-icon-pause ui-icon-shadow ui-iconsize-18">&nbsp;</span>
						</span>
					</a>
				</li>

				<li class="ui-block-d">
					<a href="#" onclick="javascript:next()" data-transition="fade" data-theme="a" data-icon="info" class="ui-btn ui-btn-up-a ui-btn-inline ui-btn-icon-top" data-corners="false" data-shadow="false" data-iconshadow="true" data-iconsize="18" data-wrapperels="span" data-iconpos="top" data-inline="true">
						<span class="ui-btn-inner">
							<span class="ui-btn-text">下一首</span>
							<span class="ui-icon ui-icon-arrow-r ui-icon-shadow ui-iconsize-18">&nbsp;</span>
						</span>
					</a>
				</li>
			</ul>
		</div>
</div>
    </body>
</html>
<script>
	$(document).ready(function(){
		loadAlbumInfo();	
	});
</script>
