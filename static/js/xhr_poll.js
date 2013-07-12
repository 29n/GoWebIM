function xhr_loop() {
	$.ajax({
		type : "POST",
		dataType : "json",
		url : "/xhr_poll/",
		timeout : 60000, // ajax请求超时时间80秒
		data : {
			time : "80"
		}, // 40秒后无论结果服务器都返回数据
		success : function(data, textStatus) {
			// 从服务器得到数据，显示数据并继续查询
			if (data.success == "1") {
				ws.onmessage({
					'data' : data.content
				});
			} else {
				console.log('No Data')
			}
			xhr_loop();
		},
		// Ajax请求超时，继续查询
		error : function(XMLHttpRequest, textStatus, errorThrown) {
			if(textStatus == "timeout") {
				xhr_loop();
			}
			console.log(textStatus)
		}
	});
}