var ws = {};
$(document).ready(function(){
	if(window.WebSocket || window.MozWebSocket){
		ws = new WebSocket("ws://" + config.host + ":"+config.port+"/chatroom");
	} else {
		ws = {}
		//jQuery.getScript("/static/js/xhr_poll.js", function(data, textStatus, jqxhr){
		//	xhr_loop();
		//});
		xhr_loop();
	}
});

var client_id;
var userlist;
var GET = new Object();
GET = GetRequest();

function GetDateT() {
	var d,s;
	d = new Date();
	s = d.getYear() + "-";             // 取年份
	s = s + (d.getMonth() + 1) + "-";// 取月份
	s += d.getDate() + " ";         // 取日期
	s += d.getHours() + ":";       // 取小时
	s += d.getMinutes() + ":";    // 取分
	s += d.getSeconds();         // 取秒
	return(s);  
}

function GetRequest() {
	   var url = location.search; // 获取url中"?"符后的字串
	   var theRequest = new Object();
	   if (url.indexOf("?") != -1) {
	      var str = url.substr(1);
	      strs = str.split("&");
	      for(var i = 0; i < strs.length; i ++) {
	         theRequest[strs[i].split("=")[0]] = decodeURIComponent(strs[i].split("=")[1]);
	      }
	   }
	   return theRequest;
	}

ws.onopen = function(e) {
	console.log("onopen");
	ws.send("set name " + GET['name']);
	ws.send("get name list");
};

function selectUser(userid) {
	$('#userlist').val(userid);
}

ws.onmessage = function(e) {
	var resp = e.data.split(" ", 3);
	var _userid;
	console.log(e.data)
	if (resp[0] == "setok") {
		client_id = parseInt(resp[1], 10);
	} else if (resp[0] == "names") {
		userlist = eval("(" + resp[1] + ")");
		for ( var i in userlist) {
			if(i == client_id) continue
			newUser(i, userlist[i]);
		}
	} else if (resp[0] == "msgto") {
		_userid = parseInt(resp[1], 10)
		newMsg("<span style='color: orange'><a href='javascript:selectUser("+_userid+")'>" + userlist[_userid]
				+ "</a></span> 对你说: " + resp[2]);
	} else if (resp[0] == "msgall") {
		_userid = parseInt(resp[1], 10)
		newMsg("<span style='color: orange'><a href='javascript:selectUser("+_userid+")'>" + userlist[_userid]
				+ "</a></span> 对所有人说: " + resp[2]);
	} else if (resp[0] == "login") {
		resp[1] = parseInt(resp[1], 10)
		newUser(resp[1], resp[2])
		userlist[resp[1]] = resp[2]
		newMsg("<span style='color: green'>【系统】</span> " + resp[2] + "已登录");
	} else if (resp[0] == "logout") {
		newMsg("<span style='color: green'>【系统】</span> "
				+ userlist[parseInt(resp[1], 10)] + "已退出");
		delUser(resp[1]);
	} else {
		console.log(e.data);
	}
};

ws.onclose = function(e) {
	console.log("onclose");
	alert('您已退出聊天室');
	location.href = '/static/html/';
};

ws.onerror = function(e) {
	console.log("onerror");
};

function delUser(userid) {
	$('#user_' + userid).remove();
	delete (userlist[userid])
}

function newUser(userid, name) {
	$('#userlist').append(
			"<option value='" + userid + "' id='user_" + userid + "'>" + name
					+ "</option>");
}

function newMsg(content, color) {
	$("#msg-template .userpic").html("<img src='" + GET['avatar'] + "'>")
	$("#msg-template .msg-time").html(GetDateT());
	$("#msg-template .content").html('<span style="color:' + color + ';">' + content + '</span>');
	$("#chat-messages").append($("#msg-template").html());
	$('#chat-messages')[0].scrollTop = $('#chat-column')[0].scrollHeight;
}

$(function() {
	$('#msgform').submit(function() {
		var content = $('#msg').val();
		if ($('#userlist').val() == 0) {
			ws.send("sendm 0 " + content);

		} else {
			ws.send("sendto " + $('#userlist').val() + " " + content);
		}
		newMsg("我：" + content, 'blue');
		$('#msg').val('');
		return false;
	});
});