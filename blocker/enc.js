

const execFile = require('child_process').execFile;
const fs = require('fs');
const path = require('path');
const doEncrypt = require('./aesutil').encrypt;
const crypto = require('crypto');

const nowdir = __dirname;
const outName = "allow.go";

const finOutName = () => path.join(nowdir, outName);


function randHintKey() {
    let rv = "";
    for (let i = 0; i < 6; i++) {
        rv += String.fromCharCode(Math.floor(Math.random() * 26) + 65);
    }
    return rv;
}

function loadChunk(filename, varname, curstr) {
    curstr = curstr ? curstr : "";
    return new Promise((r, c) => {
        fs.readFile(filename, (err, chunk) => {
            if (err) {
                c(err);
            } else {
                let hintKey = randHintKey() + "k)1`";
                let sh = crypto.createHash('sha256');
                sh.update(hintKey);
                let key = sh.digest();
                let hexStream = doEncrypt(`${chunk}`, key);
                r(curstr + `
const ${varname} = \`
${hexStream.toString('hex')}
\`

const hintKey = "${hintKey}"
`);
            }
        });
    });
}

function locateUnoGo(startd) {
    return new Promise((r, c) => {
        fs.readdir(startd, { encoding: 'utf8' }, (err, files) => {
            if (err) {
                c(err);
            } else {
                for (let f of files) {
                    // .gitignore may occur
                    if (f !== outName && f.match(/\.go$/)) {
                        r(f);
                        break;
                    }
                }
            }
        });
    });
}


locateUnoGo(nowdir).then((name) =>
    new Promise((r, c) => {
        execFile('grep', ['-i', "^package", path.join(nowdir, name)], { maxBuffer: 1024 * 1024 }, (err, stdout, stderr) => {
            if (err) {
                c(err);
            } else {
                let lines = stdout.split(/\r\n/);
                if (lines.length <= 0) {
                    c(new Error("unexpected lines"));
                    return;
                }
                let fields = lines[0].split(/\s/);
                if (fields.length <= 1) {
                    c(new Error("unexpected fields"));
                    return;
                }
                r(fields[1]); //~ package name
            }
        });
    })
)
    .then((packageName) => loadChunk(path.join(nowdir, 'allowedlist'), 'allowedConstantString', `
// This file is generated automatically.
// ${new Date()}
package ${packageName}
`))
    .then((data) =>
        new Promise((r, c) => {
            fs.writeFile(finOutName(), data, (err) => {
                if (err) {
                    c(err);
                } else {
                    r();
                }
            })
        })
    )
    .then(() => new Promise((c, r) => {
        execFile('gofmt', ['-w', '-s', finOutName()], (err, stdout, stderr) => {
            if (err) {
                console.error("error occurs during gofmt");
                c(err);
            } else {
                //console.log(stdout);
                console.log(`allow-list sync done.\n${finOutName()}`);
            }
        });
    }))
    .catch(e => {
        console.error(e);
    });
