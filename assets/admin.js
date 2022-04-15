// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

$(document).ready(function() {

  $('.message .close')
  .on('click', function() {
    $(this)
      .closest('.message')
      .transition('fade')
    ;
  });

  $('#form-gitee-authorize').form({
    fields: {
      client_id: {
        rules: [{
          type:'empty',
          prompt:'Client ID 尚未填写'
        }]
      },
      client_secret: {
        rules: [{
          type:'empty',
          prompt:'Client Secret 尚未填写'
        }]
      },
      redirect_url: {
        rules: [{
          type:'empty',
          prompt:'应用回调地址尚未填写'
        }]
      }
    }
  });

  $('#btn-gitee-authorize').click(function(){    
    authForm = $('#form-gitee-authorize');    
    authForm.form('validate form');
    if(authForm.form('is valid')) {
      clientID = authForm.form('get value', 'client_id'),
      clientSecret = authForm.form('get value', 'client_secret'),
      redirectUrl = authForm.form('get value', 'redirect_url'),      
      scopes = authForm.form('get value', 'scopes'),      
      url = `https://gitee.com/oauth/authorize?client_id=`+clientID+`&redirect_uri=`+redirectUrl+`&response_type=code&scope=`+scopes
      window.open(url);      
    }//end of if
  });  

  $('#btn-gitee-token').click(function(){
    authForm = $('#form-gitee-authorize');
    authForm.form('add rule','code',{
      rules: [{
        type: 'empty',
        prompt: '授权 Code 尚未填写'
      }]
    });
    authForm.form('validate form');
    if(authForm.form('is valid')) {
      GetGiteeToken(authForm);
    }
  });
});

function GetGiteeToken(form) {
  authForm = $('#form-gitee-authorize'); 
  var data = JSON.stringify( {
    "client_id" : form.form('get value', 'client_id'),
    "client_secret": form.form('get value', 'client_secret'),
    "redirect_url" : form.form('get value', 'redirect_url'),  
    "code" : form.form('get value', 'code')
  })
  $.ajax({
    type: "POST",
    url: '/admin/gitee/token',    
    contentType: "application/json",
    dataType: "json",
    data: data,
    success: function(data) {            
      $('#gitee-token-message').html('<code>'+data+'</code>');
      $('#gitee-token-message').removeClass('error').addClass('positive').addClass('visible');
    },
    error: function(xhr) {
      $('#gitee-token-message').html('<code>'+xhr.responseText+'</code>');
      $('#gitee-token-message').removeClass('positive').addClass('error').addClass('visible');
    }
  });
}//end of function GetGiteeToken