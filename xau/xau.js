#!/usr/bin/env node

const fs = require('fs');
const crypto = require('crypto');
const execFile = require('child_process').execFile;

//console.log("***");
//console.log("UNO:", process.argv[1]);
//console.log("DOS:", process.argv[2]);
//console.log("***");

if (process.argv.length < 4) {
	console.error("not enough parameter");
	process.exit(-1);
}

const userName = process.argv[2];
const userPassword = process.argv[3];

const TargetFileName = 'xauconf.go';

function doStart() {
	return new Promise((r, c) => {
		fs.writeFile(TargetFileName, `

// This file is generated automatically.
// ${new Date()}
package xau

const _hashSecret = "${(() => {
				let sha256 = crypto.createHash('sha256');
				sha256.update(userName + userPassword);
				return sha256.digest('hex');
			})()}"

`, { encoding: 'utf8' },
			(err) => {
				if (err) {
					console.error(err);
				} else {
					r();
				}
			});
	});
};


// Master entry.
doStart().then(() => new Promise((r, c) => {
	execFile('gofmt', ['-w', '-s', TargetFileName], { encoding: 'utf8' }, (err) => {
		if (err) {
			c(err);
		} else {
			console.log("new key generated");
			r()
		}
	});
}))
	.catch((e) => {
		console.error(e);
	});

