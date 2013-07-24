function db_connect(dbname, ok_cb, error_cb) {
	window.indexedDB = window.indexedDB || window.mozIndexedDB
			|| window.webkitIndexedDB || window.msIndexedDB;
	window.IDBTransaction = window.IDBTransaction
			|| window.webkitIDBTransaction || window.msIDBTransaction;
	window.IDBKeyRange = window.IDBKeyRange || window.webkitIDBKeyRange
			|| window.msIDBKeyRange;
	if (!window.indexedDB) {
		console.log("Your browser doesn't support a stable version of IndexedDB. Such and such feature will not be available.");
		return false;
	}
	var request = window.indexedDB.open("MyTestDatabase", 3);
	if(error_cb== undefined) {
		request.onerror = function(event) {
			console.log(event);
		}
	} else {
		request.onerror = error_cb;
	}
	request.onsuccess = function(event) {
		var db;
		var request2 = indexedDB.open(dbname);
		request2.onerror = error_cb;
		request2.onsuccess = function(event) {
			db = request2.result;
			ok_cb(db);
		};
	}
}

function db_table(tablename, primary, ok_cb, error_cb) {
	var objectStore = db.createObjectStore(tablename, {  
        // primary key  
        keyPath: primary,  
        // auto increment  
        autoIncrement: false  
    });  
	
}