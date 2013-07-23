var ws = {};
var client_id = 0;
var userlist = {};
var GET = new Object();
GET = GetRequest();

$(document).ready(function(){
	if(window.WebSocket || window.MozWebSocket){
		ws = new WebSocket("ws://" + config.host + ":"+config.port+"/chatroom");
		ws.onopen = function(e) {
			if(config.debug) {
				console.log("connected"+config.host+":"+config.port);
			}
			if(GET['name']==undefined || GET['avatar']==undefined) {
				alert('非法请求');
				ws.close();
				return false;
			}
			ws.send("set name " + GET['name']+'|'+GET['avatar']);
			ws.send("get name list");
		};

		function selectUser(userid) {
			$('#userlist').val(userid);
		}

		ws.onmessage = function(e) {
			if(config.debug) {
				console.log("websocket recv: "+e.data);
			}
			var resp = e.data.split(" ", 3);
			var _userid;
			if (resp[0] == "setok") 
			{
				client_id = resp[1];
			} 
			else if (resp[0] == "names") 
			{
				userlist = eval("(" + resp[1] + ")");
				for (var i in userlist)
				{
					newUser(userlist[i].Id);
				}
			} 
			else if (resp[0] == "msgfrom") 
			{
				_userid = resp[1], 10
				newMsg(_userid, "对你说: " + resp[2]);
			} 
			else if (resp[0] == "msgall") 
			{
				_userid = resp[1]
				newMsg(_userid, "说: " + resp[2]);
			} 
			else if (resp[0] == "login") 
			{
				var _u = resp[2].split("|", 2);
				userlist[resp[1]] = {Name:_u[0], Avatar:_u[1], Id: resp[1]}
				newUser(resp[1])
				newMsg(0, _u[0] + "已登录");
			} 
			else if (resp[0] == "logout") 
			{
				newMsg(0, userlist[resp[1]].Name + "已退出");
				delUser(resp[1]);
			} 
			else 
			{
				console.log(e.data);
			}
		};
		ws.onclose = function(e) {
			console.log("onclose");
			alert('您已退出聊天室');
			if(!config.debug) {
				location.href = '/static/html/';
			}
		};
		ws.onerror = function(e) {
			console.log("onerror");
		};
	} else {
		ws = {}
		// jQuery.getScript("/static/js/xhr_poll.js", function(data, textStatus,
		// jqxhr){
		// xhr_loop();
		// });
		xhr_loop();
	}
});

function xssFilter(val) {
    val = val.toString();
    val = val.replace(/[<%3C]/g, "&lt;");
    val = val.replace(/[>%3E]/g, "&gt;");
    val = val.replace(/"/g, "&quot;");
    val = val.replace(/'/g, "&#39;");
    return val;
}

function GetDateT() {
	var d, s;
	s = '';
	d = new Date();
	//s = d.getYear() + "-";             // 取年份
	//s = s + (d.getMonth() + 1) + "-";// 取月份
	//s += d.getDate() + " ";         // 取日期
	s += d.getHours() + ":";       // 取小时
	s += d.getMinutes();           // 取分
	//s += d.getSeconds();         // 取秒
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

function selectUser(userid) {
	$('#userlist').val(userid);
}

function delUser(userid) {
	$('#user_' + userid).remove();
	$('#inroom_' + userid).remove();
	delete (userlist[userid])
}

function newUser(userid) {
	if(userid != client_id) {
		$('#userlist').append(
				"<option value='" + userid + "' id='user_" + userid + "'>" + userlist[userid].Name
						+ "</option>");
	}
	$('#left-userlist').append(
			"<li id='inroom_"+userid+"'><a href='javascript:selectUser("+userid+")'>" +
					"<img src='" + userlist[userid].Avatar + "' width='50' height='50'></a></li>"
	);
}

function newMsg(fromId, content, color) {
	content = xssFilter(content)
	$("#msg-template .msg-time").html(GetDateT());
	if(fromId == 0){
		$("#msg-template .userpic").html("");
		$("#msg-template .content").html("<span style='color: green'>【系统】</span>" + content);
	}
	else {
		var html = '';
		//$("#msg-template .userpic").html("<img src='" + userlist[fromId].Avatar + "' width='50' height='50'>")
		if(client_id == fromId){
			html += '<span style="color: orange">我说: </span>: ';
			
		} else {
			html += '<span style="color: orange"><a href="javascript:selectUser('+fromId+')">' + userlist[fromId].Name;
			html += '</a></span> '
		}
		html += content + '</span>';
		$("#msg-template .content").html(html);
	}
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
		newMsg(client_id, content);
		$('#msg').val('');
		return false;
	});
});