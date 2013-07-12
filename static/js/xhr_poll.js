function xhr_loop() {
	$.ajax({
		type : "POST",
		dataType : "json",
		url : "/xhr_poll/",
		timeout : 60000, // ajax����ʱʱ��80��
		data : {
			time : "80"
		}, // 40������۽������������������
		success : function(data, textStatus) {
			// �ӷ������õ����ݣ���ʾ���ݲ�������ѯ
			if (data.success == "1") {
				ws.onmessage({
					'data' : data.content
				});
			} else {
				console.log('No Data')
			}
			xhr_loop();
		},
		// Ajax����ʱ��������ѯ
		error : function(XMLHttpRequest, textStatus, errorThrown) {
			if(textStatus == "timeout") {
				xhr_loop();
			}
			console.log(textStatus)
		}
	});
}