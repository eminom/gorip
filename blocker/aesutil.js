const crypto = require('crypto');
const assert = require('assert');

//buff
function padKey(buffStr) {
    let buffKey = Buffer.from(buffStr);
    let targetLen = 16;
    if (buffKey.length <= 16) {
        targetLen = 16;
    } else if (buffKey.length <= 24) {
        targetLen = 24;
    } else {
        targetLen = 32;
    }

    let r = targetLen - buffKey.length;
    if (r > 0) {
        let padding = Buffer.alloc(r, r);
        buffKey = Buffer.concat([buffKey, padding]);
    }
    return buffKey.slice(0, targetLen);
}

function randInt8() {
    return Math.floor(Math.random() * 256);
}

// a buffer of length equals block-size
function randIv() {
    let rv = Buffer.alloc(16, 0);
    rv.forEach((c, i) => {
        rv[i] = randInt8();
    });
    return rv;
}

function makeCipherName(k) {
    let cipherName = 'aes-128-cbc';  //cipher block chain
    switch (k.length) {
        case 16:
            cipherName = 'aes-128-cbc';
            break;
        case 24:
            cipherName = 'aes-192-cbc';
            break;
        case 32:
            cipherName = 'aes-256-cbc';
            break;
        default:
            throw new Error(`unknown key length ${k.length}`);
    }
    return cipherName;
}

function makeCipher(k, iv) {
    let algoName = makeCipherName(k);
    return crypto.createCipheriv(algoName, k, iv);
}

function makeDecipher(k, iv) {
    let algoName = makeCipherName(k);
    return crypto.createDecipheriv(algoName, k, iv);
}


// using uint64 as the head
function encrypt(source, keyStr) {
    let srcBuff = Buffer.from(source);
    let headBuff = Buffer.alloc(8, 0)
    headBuff.writeUInt32BE(srcBuff.length, 4);
    srcBuff = Buffer.concat([headBuff, srcBuff]);
    let r = srcBuff.length % 16;
    if (r > 0) {
        srcPadding = Buffer.alloc(16 - r, 0);
        srcPadding.forEach((c, i) => { srcPadding[i] = randInt8() });
    }
    let startiv = randIv();
    let cipher = makeCipher(padKey(keyStr), startiv);

    let buffs = [startiv];
    buffs.push(cipher.update(srcBuff));
    buffs.push(cipher.final());
    return Buffer.concat(buffs);
}

function decrypt(chunk, keyStr) {
    if (Buffer.isBuffer(chunk)) {
        //nada
    } else {
        chunk = Buffer.from(chunk, 'hex');
    }
    let iv = chunk.slice(0, 16)
    chunk = chunk.slice(16);
    let cipher = makeDecipher(padKey(keyStr), iv);
    let decoded = [];
    decoded.push(cipher.update(chunk));
    decoded.push(cipher.final());
    decoded = Buffer.concat(decoded);
    let len = decoded.readUInt32BE(4);
    return decoded.slice(8, 8 + len);
}

module.exports = {
    encrypt: encrypt,
    decrypt: decrypt,
};

