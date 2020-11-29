const { dialog } = require('electron').remote

let index = {
	init: function() {
		// Wait for astilectron to be ready
		document.addEventListener('astilectron-ready', function() {
		    // Listen
		    index.listen();
		    astilectron.sendMessage({name: "getVersion", payload: ""}, function(message) {
			    index.next("version", message.payload);
		    });
		    astilectron.sendMessage({name: "getAuthed", payload: ""}, function(message) {
			    index.next("auth_status", message.payload);
		    });
		    astilectron.sendMessage({name: "getContactList", payload: ""}, function(message) {
			    index.send("get_clist", "getContactList", "");
		    });
		    astilectron.sendMessage({name: "keyExists", payload: ""}, function(message) {
			    index.send("key_exists", "keyExists", "")
		    });

		})
	},
	listen: function() {
		astilectron.onMessage(function(message) {
		    return {payload: "payload"};
		});
    	},
	send: function(t, n, p) {
		console.log("Send: " + t);
		console.log("Send: " + n);
		console.log("Send: " + p);
		astilectron.sendMessage({name: n, payload: p}, function(message) {
			console.log("Inside send message...");
			console.log("t: " + t);
			console.log("received: " + message.payload);
			index.next(t, message.payload);
		});
	},
	next: function(type, s) {
		switch(type) {
			case "sleep":
				console.log("Done Sleeping...");
				var x = document.getElementById("auth_input");
				var y = document.getElementById("auth_msg");
				y.style.display = "none";
				x.style.display = "none";
				break;
			case "version":
				document.getElementById("version").innerHTML = s;
				break;
			case "auth_status":
				var z = document.getElementById("status");	
				if (s == true) {
					z.innerHTML = "Authorized";
				} else {
					show("info", "Authorization", "Your app has not been authorized yet. Do the following to authorize this app:\n\n1. Generate and copy an app access secret from InstaCrypt Vault in the \"App Access\" section\n\n2. Come back here to click on the \"Authorize\" button to paste the app secret in the field that appears\n\n3. Hit \"Save\".");
				}
				break;
			/*
			 case "dialog":
				console.log("Dialog Status: " + s);
				document.getElementById(s.type).innerText = s.file;
				document.getElementById(s.type).style.border = "3px dashed black";
				break;
			*/
			case "get_clist":
				console.log("Contact List:");
				console.log(s);
				var clist = s.List;
				var count = clist.length;
				var contacts = document.getElementById("contacts");
				var i;
				for (i=0; i < count; i++) {
					var option = document.createElement("option");
					option.text = clist[i];
					option.value = i;
					contacts.options.add(option, i+1);
				}
				break;
			case "auth_save":
				console.log(s);
				var x = document.getElementById("auth_input");
				var y = document.getElementById("auth_msg");
				var z = document.getElementById("status");
				if (s == false) {
					y.innerHTML = "Unable to save access secret..."
					y.style.display = "block";
					z.inner.HTML = "Not Authorized";
					return;
				} else {
					y.innerHTML = "Succesfully saved!"
					y.style.display = "block";
					index.send("sleep", "sleep", "");
					z.innerHTML = "Authorized";
					index.send("get_server_con_to_verify", "getServerContacts", "");
				}
				break;
			case "post_keygen":
				if (s.Ok == true) {
					ask(s.Title, s.Message, "catchall", "keygened");
					var gk = document.getElementById("genkey");
					gk.innerHTML = "<p class='icon' ><svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><g><path d='M16.57,1.5h-4.16C12.36,0.66,11.66,0,10.82,0H9.18C8.34,0,7.64,0.66,7.59,1.5H3.43c-0.43,0-0.77,0.35-0.77,0.77s0.35,0.77,0.77,0.77h4.93c0.43,0,0.77-0.35,0.77-0.77V1.61c0-0.04,0.03-0.06,0.05-0.06h1.64c0.02,0,0.05,0.02,0.05,0.06v0.67c0,0.43,0.35,0.77,0.77,0.77h4.93c0.43,0,0.77-0.35,0.77-0.77S16.99,1.5,16.57,1.5z'/><path d='M3.48,12.89C3.48,15.16,5.3,17,7.54,17h4.93c2.24,0,4.06-1.84,4.06-4.11V3.78H3.48V12.89z M12.08,6.79c0-0.21,0.17-0.39,0.39-0.39s0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39s-0.39-0.17-0.39-0.39V6.79z M9.61,6.79C9.61,6.57,9.79,6.4,10,6.4s0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39s-0.39-0.17-0.39-0.39V6.79z M7.15,6.79c0-0.21,0.17-0.39,0.39-0.39c0.21,0,0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39c-0.21,0-0.39-0.17-0.39-0.39V6.79z'/></g></svg></p>Delete<br />My Lock";
					gk.onclick = function(){preKeyDel()};
				}
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				}
				break;
			case "post_keydel":
				if (s.Ok == true) {
					show("info", s.Title, s.Message)
					var dk = document.getElementById("genkey");
					dk.innerHTML = "<p class='icon' ><svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><g><path d='M17.96,13.88L11.98,7.9l0,0c0,0,0,0,0,0L9.54,5.46c-0.35-0.35-0.8-0.54-1.29-0.54c-0.49,0-0.95,0.19-1.29,0.54C6.61,5.8,6.42,6.26,6.42,6.75c0,0.49,0.19,0.95,0.54,1.29l8.43,8.43c0.35,0.35,0.8,0.54,1.29,0.54c0.49,0,0.95-0.19,1.29-0.54s0.54-0.8,0.54-1.29C18.5,14.68,18.31,14.23,17.96,13.88z M9.4,7.9C9.26,8.03,9.16,8.19,9.07,8.36L7.85,7.14c-0.22-0.22-0.22-0.57,0-0.79c0.22-0.22,0.57-0.22,0.79,0l1.22,1.22C9.69,7.66,9.53,7.76,9.4,7.9z'/><path d='M8.73,2.52c0.23,0,0.45,0.1,0.61,0.29C9.52,3,9.62,3.27,9.62,3.55c0,0.22,0.18,0.4,0.4,0.4s0.4-0.18,0.4-0.4c0-0.28,0.1-0.55,0.27-0.75c0.17-0.19,0.38-0.29,0.61-0.29c0.22,0,0.4-0.18,0.4-0.4c0-0.22-0.18-0.4-0.4-0.4c-0.23,0-0.45-0.1-0.61-0.29c-0.18-0.2-0.27-0.46-0.27-0.75c0-0.22-0.18-0.4-0.4-0.4s-0.4,0.18-0.4,0.4c0,0.28-0.1,0.55-0.27,0.75C9.17,1.62,8.96,1.73,8.73,1.73c-0.22,0-0.4,0.18-0.4,0.4C8.33,2.34,8.51,2.52,8.73,2.52z M9.93,1.97c0.03-0.03,0.05-0.07,0.08-0.11c0.03,0.04,0.05,0.08,0.08,0.11c0.05,0.06,0.1,0.11,0.16,0.15c-0.05,0.05-0.11,0.1-0.16,0.15c-0.03,0.03-0.05,0.07-0.08,0.11C9.98,2.35,9.96,2.31,9.93,2.28c-0.05-0.06-0.1-0.11-0.16-0.15C9.83,2.07,9.88,2.02,9.93,1.97z'/><path d='M5.79,1.91c-0.38,0-0.74-0.15-1.01-0.42C4.51,1.22,4.36,0.86,4.36,0.48C4.36,0.21,4.15,0,3.88,0S3.41,0.21,3.41,0.48c0,0.38-0.15,0.74-0.42,1.01C2.72,1.76,2.36,1.91,1.98,1.91c-0.26,0-0.48,0.21-0.48,0.48s0.21,0.48,0.48,0.48c0.38,0,0.74,0.15,1.01,0.42c0.27,0.27,0.42,0.63,0.42,1.01c0,0.26,0.21,0.48,0.48,0.48s0.48-0.21,0.48-0.48c0-0.38,0.15-0.74,0.42-1.01c0.27-0.27,0.63-0.42,1.01-0.42c0.26,0,0.48-0.21,0.48-0.48S6.05,1.91,5.79,1.91z M4.1,2.6C4.02,2.69,3.95,2.77,3.88,2.86C3.82,2.77,3.74,2.69,3.66,2.6C3.58,2.52,3.49,2.45,3.4,2.38c0.09-0.07,0.18-0.14,0.26-0.22C3.74,2.08,3.82,1.99,3.88,1.9C3.95,1.99,4.02,2.08,4.1,2.16c0.08,0.08,0.17,0.15,0.26,0.22C4.27,2.45,4.19,2.52,4.1,2.6z'/><path d='M5.32,6.22c-0.2,0-0.38-0.09-0.52-0.24C4.65,5.81,4.57,5.59,4.57,5.35c0-0.18-0.14-0.32-0.32-0.32c-0.18,0-0.32,0.14-0.32,0.32c0,0.24-0.08,0.46-0.23,0.63C3.56,6.13,3.37,6.22,3.18,6.22C3,6.22,2.86,6.36,2.86,6.54c0,0.18,0.14,0.32,0.32,0.32c0.2,0,0.38,0.09,0.52,0.24c0.15,0.17,0.23,0.39,0.23,0.63c0,0.18,0.14,0.32,0.32,0.32c0.18,0,0.32-0.14,0.32-0.32c0-0.24,0.08-0.46,0.23-0.63c0.14-0.16,0.33-0.24,0.52-0.24c0.18,0,0.32-0.14,0.32-0.32C5.64,6.36,5.5,6.22,5.32,6.22z M4.33,6.67C4.3,6.71,4.28,6.74,4.25,6.78C4.22,6.74,4.2,6.71,4.17,6.67c-0.05-0.05-0.1-0.09-0.15-0.14c0.05-0.04,0.1-0.09,0.15-0.14C4.2,6.37,4.22,6.33,4.25,6.3C4.28,6.33,4.3,6.37,4.33,6.4c0.05,0.05,0.1,0.09,0.15,0.14C4.43,6.58,4.37,6.62,4.33,6.67z'/></g></svg></p>Generate<br />Lock";
					dk.onclick = function(){preKeyGen()};
				}
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				}
				// document.getElementById('loader').style.visibility = "hidden";
				break;
			case "clear_key":
				index.send("catchall", "clearKey", "");
				break;
			case "set_file":	
				console.log(s.Type);
				console.log(s.File);
				var k = document.getElementById("file");
				k.style.border = '3px dashed #000000';
				k.innerHTML = s.File;
				setDir(s.Type);	
				break;
			case "set_3df":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				}
				console.log(s);
				break;
			case "set_local":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				}
				console.log(s);
				break;
			case "set_selected_contact":
				var k = document.getElementById("3dfkey");	
				var status = false;
				if (s.Ok == false) {
					status =  false;
					var c = document.getElementById("contacts");
					c.selectedIndex = 0;
					show("error", s.Title, s.Message);
					return;
				}
				k.style.border = "0px";
				break;
			case "key_exists":
				var el = document.getElementById('genkey');
				if(s){
					el.innerHTML = "<p class='icon'><svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><g><path d='M16.57,1.5h-4.16C12.36,0.66,11.66,0,10.82,0H9.18C8.34,0,7.64,0.66,7.59,1.5H3.43c-0.43,0-0.77,0.35-0.77,0.77s0.35,0.77,0.77,0.77h4.93c0.43,0,0.77-0.35,0.77-0.77V1.61c0-0.04,0.03-0.06,0.05-0.06h1.64c0.02,0,0.05,0.02,0.05,0.06v0.67c0,0.43,0.35,0.77,0.77,0.77h4.93c0.43,0,0.77-0.35,0.77-0.77S16.99,1.5,16.57,1.5z'/><path d='M3.48,12.89C3.48,15.16,5.3,17,7.54,17h4.93c2.24,0,4.06-1.84,4.06-4.11V3.78H3.48V12.89z M12.08,6.79c0-0.21,0.17-0.39,0.39-0.39s0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39s-0.39-0.17-0.39-0.39V6.79z M9.61,6.79C9.61,6.57,9.79,6.4,10,6.4s0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39s-0.39-0.17-0.39-0.39V6.79z M7.15,6.79c0-0.21,0.17-0.39,0.39-0.39c0.21,0,0.39,0.17,0.39,0.39v6.67c0,0.21-0.17,0.39-0.39,0.39c-0.21,0-0.39-0.17-0.39-0.39V6.79z'/></g></svg></p>Delete<br />My Lock";
					el.onclick = function(){preKeyDel()};
				}
				break;
			case "check_refreshable":
				if (s == false) {
					show("info", "unauthorized", "You must authorize this app before using this feature.");
				} else {
					index.send("get_server_con", "getServerContacts", "");
					var contacts = document.getElementById("contacts").options.length = 1;
				}
				break;
			case "get_server_con_to_verify":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
					document.getElementById('loader').style.visibility = "hidden";
					document.getElementById('submit').style.visibility = "visible";
				} else {
					index.send("get_new_clist", "getContactList", "");
				}
				break;
			case "get_server_con":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
					document.getElementById('loader').style.visibility = "hidden";
					document.getElementById('submit').style.visibility = "visible";
				} else {
					index.send("get_refreshed_clist", "getContactList", "");
				}
				break;
			case "get_new_clist":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
					break;
				}
				var clist = s.List;
				var count = clist.length;
				var contacts = document.getElementById("contacts");
				var i;
				for (i=0; i < count; i++) {
					var option = document.createElement("option");
					option.text = clist[i];
					option.value = i;
					contacts.options.add(option, i+1);
				}
				show("info", "Authorized!", "Your app is now authorized and your contact list is loaded! Make sure you click the \"Generate Lock\" button now to generate your unique Lock and Key so other users can encrypt files that only you can decrypt.");
				//show("info", s.Title, s.Message);
				document.getElementById('loader').style.visibility = "hidden";
				document.getElementById('submit').style.visibility = "visible";
				break;
			case "get_refreshed_clist":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
					break;
				}
				var clist = s.List;
				var count = clist.length;
				var contacts = document.getElementById("contacts");
				var i;
				for (i=0; i < count; i++) {
					var option = document.createElement("option");
					option.text = clist[i];
					option.value = i;
					contacts.options.add(option, i+1);
				}
				show("info", s.Title, s.Message);
				document.getElementById('loader').style.visibility = "hidden";
				document.getElementById('submit').style.visibility = "visible";
				break;
			case "submitted":
				document.getElementById('loader').style.visibility = "hidden";
				document.getElementById('submit').style.visibility = "visible";
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				} 
				if (s.Ok == true) {
					ask(s.Title, s.Message, "catchall", "done");
				}
				break;
			case "return_res":
				return s;
			case "catchall":
				if (s.Ok == false) {
					show("error", s.Title, s.Message);
				}
				document.getElementById('loader').style.visibility = "hidden";
				break;

		}

	}

}

function allowDrop(ev) {
	ev.preventDefault();
	file.focus();
}

function drag(ev, type) {
	ev.preventDefault();

	var file = document.getElementById(type);
	file.style.border = '3px dashed #000000';
}

function leave(ev, type) {
	ev.preventDefault();

	var file = document.getElementById(type);
	file.style.border = '0px';
}

function drop(ev, type) {
	var file = ev.dataTransfer.files[0].path;	
	
	console.log("Full Path: ");
	console.log(file);
	console.log(type);
	index.send("catchall", "clear_key", "");
	if (type == "file") {
		index.send("set_file", "setFile", file);
	}
	if (type == "key") {
		index.send("set_local_key", "setLocalKey", file);
	}
}

function fileDialog(type) {
	dialog.showOpenDialog({
		properties: ['openFile']
	}).then(result => {
	  	console.log("Cancelled: " + result.canceled);
	  	console.log("File: " + result.filePaths[0]);
		index.send("catchall", "clear_key", "");
		if (type == "file" && result.canceled == false) {
			index.send("set_file", "setFile", result.filePaths[0]);
		}

		if (type == "key" && result.canceled == false) {
			index.send("set_local_key", "setLocalKey", result.filePaths[0]);
		}
	}).catch(err => {
	  console.log(err)
	})
}

function ask(t, msg, after, next) {
	dialog.showMessageBox({
		type: "question",
		title: t,
		message: msg,
		buttons: ["Yes", "No"]
	}).then(result => {
		console.log(result.response);
		if (result.response == 0) {
			index.send(after, next, "");
		}
		if (result.response == 1) {
			document.getElementById('loader').style.visibility = "hidden";
		}
	}).catch(err => {
	  console.log(err)
	}); 
}

function show(type, t, msg) {
	dialog.showMessageBox({
		type: type,
		title: t,
		message: msg,
		buttons: ["OK"]
	}).then(result => {
		if (result.response == 0) {
			document.getElementById('loader').style.visibility = "hidden";
		}
	});
}

function authInput() {
	var x = document.getElementById("auth_input");
	if (x.style.display === "none") {
		x.style.display = "block";
	} else {
		x.style.display = "none";
	}
}

function setContact(selected) {
	var k = document.getElementById("key");
	var cname = selected.options[selected.selectedIndex].text;
	var check = selected.options[selected.selectedIndex].value;
	console.log(check);
	if (check != "none") {
		var id = Number(check);
		console.log(id);
		console.log(cname);
		index.send("set_selected_contact", "setContactSelected", id);
	} else {
		k.style.border = "0px";
		var cd = document.getElementById("direction");
		var checkd = cd.options[cd.selectedIndex].value;
		if (checkd  == "encrypt") {
			k.innerHTML = "<svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><path d='M15.49,6.33C15.36,3.69,14.38,0,10,0S4.64,3.69,4.51,6.33C3.49,6.55,2.73,7.47,2.73,8.57v7.66c0,0.42,0.34,0.77,0.77,0.77h13.02c0.42,0,0.77-0.34,0.77-0.77V8.57C17.27,7.47,16.51,6.55,15.49,6.33z M10,1.84c1.49,0,3.43,0.48,3.65,4.44H6.35C6.57,2.32,8.51,1.84,10,1.84z M11.07,13.9H8.93l0.73-2.28c-0.46-0.15-0.8-0.58-0.8-1.09c0-0.63,0.51-1.15,1.15-1.15c0.63,0,1.15,0.51,1.15,1.15c0,0.51-0.34,0.94-0.8,1.09L11.07,13.9z'/></svg><br />LOCK";
		} else {
			k.innerHTML = "<svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><path d='M15.28,3.78c-2.18,0-4,1.48-4.55,3.48c-0.06-0.01-0.11-0.02-0.17-0.02H1.13C0.51,7.25,0,7.75,0,8.38v2.83c0,0.63,0.51,1.13,1.13,1.13s1.13-0.51,1.13-1.13v-1.7h1.11v0.99c0,0.63,0.51,1.13,1.13,1.13s1.13-0.51,1.13-1.13V9.51h4.93c0.04,0,0.07-0.01,0.11-0.01c0.46,2.12,2.35,3.72,4.61,3.72c2.61,0,4.72-2.11,4.72-4.72S17.89,3.78,15.28,3.78z M16.28,9.5c-0.55,0.55-1.45,0.55-2,0c-0.55-0.55-0.55-1.45,0-2s1.45-0.55,2,0C16.84,8.05,16.84,8.95,16.28,9.5z'/></svg><br />KEY";
		}
	}


}

function setDir(d) {
	var s3d = document.getElementById("3dfkey");
	var con = document.getElementById("contacts");

	if (d == ".icfx") {
		s3d.innerHTML = "<p class='icon' ><svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><path d='M15.28,3.78c-2.18,0-4,1.48-4.55,3.48c-0.06-0.01-0.11-0.02-0.17-0.02H1.13C0.51,7.25,0,7.75,0,8.38v2.83c0,0.63,0.51,1.13,1.13,1.13s1.13-0.51,1.13-1.13v-1.7h1.11v0.99c0,0.63,0.51,1.13,1.13,1.13s1.13-0.51,1.13-1.13V9.51h4.93c0.04,0,0.07-0.01,0.11-0.01c0.46,2.12,2.35,3.72,4.61,3.72c2.61,0,4.72-2.11,4.72-4.72S17.89,3.78,15.28,3.78z M16.28,9.5c-0.55,0.55-1.45,0.55-2,0c-0.55-0.55-0.55-1.45,0-2s1.45-0.55,2,0C16.84,8.05,16.84,8.95,16.28,9.5z'/></svg></p>Use<br />My Key";
		s3d.onclick = function(){set3DF(s, "Key")};
		con.style.visibility = "hidden";
	} else {
		s3d.innerHTML = "<p class='icon' ><svg xmlns='http://www.w3.org/2000/svg' width='20' height='17' wBox='0 0 20 17'><path d='M15.49,6.33C15.36,3.69,14.38,0,10,0S4.64,3.69,4.51,6.33C3.49,6.55,2.73,7.47,2.73,8.57v7.66c0,0.42,0.34,0.77,0.77,0.77h13.02c0.42,0,0.77-0.34,0.77-0.77V8.57C17.27,7.47,16.51,6.55,15.49,6.33z M10,1.84c1.49,0,3.43,0.48,3.65,4.44H6.35C6.57,2.32,8.51,1.84,10,1.84z M11.07,13.9H8.93l0.73-2.28c-0.46-0.15-0.8-0.58-0.8-1.09c0-0.63,0.51-1.15,1.15-1.15c0.63,0,1.15,0.51,1.15,1.15c0,0.51-0.34,0.94-0.8,1.09L11.07,13.9z'/></svg></p>Use<br />My Lock";
		s3d.onclick = function(){set3DF(s, "Lock")};
		con.style.visibility = "visible";
	}
	s3d.style.border = "3px dashed #000000";
}

function setd(d) {
	index.send("set_dir", "setDirection", d);
}

function disableDragAndDrop(ev) {
	ev.preventDefault();
}

function disableDragLeave(ev) {
	ev.preventDefault();
}

function preKeyGen() {
	document.getElementById('loader').style.visibility = "visible";
	index.send("post_keygen", "keygen", "");
}

function preKeyDel() {
	document.getElementById('loader').style.visibility = "visible";
	ask("Delete Lock & Key", "\nAre you sure you want to delete your lock & key? This is unrecoverable...", "post_keydel", "keydel");
}

function set3DF(d, t) {
	var k = document.getElementById("3dfkey");
	k.style.border = '3px dashed #000000';
	index.send("set_3df", "setKey", {key: "3df", dir: d});
}

function refreshCList() {
	document.getElementById('submit').style.visibility = "hidden";
	document.getElementById('loader').style.visibility = "visible";
	var authed = index.send("check_refreshable", "getAuthed", "");
}

function preSubmit() {
	document.getElementById('submit').style.visibility = "hidden";
	document.getElementById('loader').style.visibility = "visible";
	index.send("submitted", "submit", "");
}

