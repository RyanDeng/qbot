#!/usr/bin/env node

var request = require('request')
var QQ = require('../../wqq')
var clc = require('cli-color')
var async = require('async')
var express = require('express')
var bodyParser = require('body-parser');

var req = request.defaults({
  headers: {
    'Connection': 'keep-alive'
  }
})
var qq = new QQ()
var port = 8890

var red = clc.red

//var backendHost = 'http://127.0.0.1:8870'
var backendHost = 'http://192.168.200.105:8870'
var mainGroupNum = 472281325
//var account = '2892434023'
//var password = 'qboxtest'
var account = '286981980'
var password = ''


var run = function() {
  login(account, password, function(e, d) {
    if (e || !d) {
      console.log('登陆失败:', e)
    } else {
      qqReady(d, run)
    }
  })
}

//run()

var app = express()
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: false }))

// {"to": <to>, "msg": <msg>}
app.post("/msg", function(req, res){
  b = req.body
  console.log(b)
  qq.sendStrangerMsg(mainGroupNum, parseInt(b.to) || 0, b.msg || '', function(e, ok) {
    if (e) {
      console.log(e)
      res.status(500).end()
    } else {
      res.status(200).end()
    }
  })
})

// {"group": <group>, "msg": <msg>}
app.post("/grpMsg", function(req, res){
  b = req.body
  console.log(b)
  qq.sendGroupMsg2(parseInt(b.group) || 0, b.msg || '', function(e, ok) {
    if (e) {
      console.log(e)
      res.status(500).end()
    } else {
      res.status(200).end()
    }
  })
})

app.listen(port, function(e) {
  console.log("listen:", port, e)
  run()
})

function login(acc, psd, cb) {

  var vcode = ''
  qq.getVcode(acc, function(e, buffer) {
    if (e) {
      console.log("getVcode error:", e)
      return cb(e)
    }
    if (buffer) {
      // do not implement vcode now
      console.log("need vcode")
      return cb({error:'need vcode'})
    }
    qq.login(psd || '', vcode, function(e, d) {
      cb(e, d)
    })
  })
}

function qqReady(d, cb) {
  if (!d) {
    return
  }
  qq.getSelfInfo(function (e, d) {
    console.log('登陆成功: ',  d.nick)
    qq.cacheAllGroupMemberUin(function(e) {
      if (e) {
        console.log(red('缓存失败'))
      }
    })
    qq.on('disconnect', function() {
      console.log(red('连接断开'))
      //qq.stopPoll()
      //cb()
    })
    qq.on('message', function(d) {
      // do with message
      console.log(d)
      handleMsg(d)
    })
    qq.startPoll()
  })
}

function handleMsg(d) {
  if (!d) {
    return
  }

  if (d.poll_type === 'group_message') {
    if (d.content[0].indexOf("@QBot") != -1) {
      sendGrpMsgToBackend(qq.getAccount(d.group_code), d.send_account, d.content[0], function(e, code) {
        console.log(code)
      })
    }
    return
  }
  if (d.poll_type === 'message' || d.poll_type === '')  {
    sendMsgToBackend(d.send_account, d.content[0], function(e, code) {
      console.log(code)
    })
  }
  if (d.poll_type === 'sess_message')  {
    sendMsgToBackend(d.ruin, d.content[0], function(e, code) {
      console.log(code)
    })
  }
}


function sendMsgToBackend(from, msg, cb) {
  req.post({
    url: backendHost + '/msg',
    form: {
      'from': String(from),
      'msg': String(msg)
    }
  }, function(e, r, d) {
    cb(e, 0)
  })
}

function sendGrpMsgToBackend(grp, from, msg, cb) {
  req.post({
    url: backendHost + '/grpMsg',
    form: {
      'group': String(grp),
      'from': String(from),
      'msg': String(msg)
    }
  }, function(e, r, d) {
    cb(e, 0)
  })
}

