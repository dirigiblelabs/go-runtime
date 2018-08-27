var request = require('http/v3/request');
var response = require('http/v3/response');
response.println('Your request: [' + request.getMethod() + '] ' + request.getPath());
