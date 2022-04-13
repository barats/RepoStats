$(document).ready(function() {

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
      //TODO: go ajax
      alert('go ajax');
    }
  });
});