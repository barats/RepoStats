// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

$(document).ready(function() {

  $('#login-form')
  .form({
    fields: {
      account: {
        identifier  : 'account',
        rules: [
          {
            type   : 'empty',
            prompt : '账户名不能为空'
          },
          {
            type   : 'length[5]',
            prompt : '账户名长度不得少于5位'
          }
        ]
      },
      password: {
        identifier  : 'password',
        rules: [
          {
            type   : 'empty',
            prompt : '密码不能为空'
          },
          {
            type   : 'length[8]',
            prompt : '密码长度不得少于8位'
          }
        ]
      },
      captcha: {
        identifier  : 'captcha-text',
        rules: [
          {
            type   : 'empty',
            prompt : '验证码不能为空'
          },
          {
            type   : 'length[6]',
            prompt : '验证码长度不得少于6位'
          }
        ]
      }
    }
  });

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

  $('#form-grafana-token').form({
    fields: {
      account: {
        rules: [{
          type:'empty',
          prompt:'帐号尚未填写'
        }]
      },
      password: {
        rules: [{
          type:'empty',
          prompt:'密码尚未填写'
        }]
      },
      host: {
        rules: [{
          type:'empty',
          prompt:'Grafana Host 尚未填写'
        }]
      },
      port: {
        rules: [{
          type:'empty',
          prompt:'Grafana Port 尚未填写'
        }]
      } 
    }
  });

  $('#add-repo-form').form({
    fields: {
      repo_url : {
        rules:[{
          type: 'empty',
          prompt: '仓库链接尚未填写'
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

  $('#btn-grafana-token').click(function(){
    tokenForm = $('#form-grafana-token');    
    tokenForm.form('validate form');
    if(tokenForm.form('is valid')) {
      CreateGrafanaToken(tokenForm)
    }
  });
});

function successToast(message) {
  $('body').toast({
    class: 'success',
    displayTime: 2500,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom',
    onHidden: function() {location.reload()}
  });
}

function errorToast(message) {
  $('body').toast({
    class: 'error',
    displayTime: 2500,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom'
  });
}

function CreateGrafanaToken(form) {
  authForm = $('#form-gitee-authorize'); 
  var data = JSON.stringify( {
    "account" : form.form('get value', 'account'),
    "password": form.form('get value', 'password'),
    "host" : form.form('get value', 'host'),
    "port" : form.form('get value', 'port')    
  })
  $.ajax({
    type: "POST",
    url: '/admin/grafana/token',    
    contentType: "application/json",
    dataType: "json",
    data: data,
    success: function(data) {            
      $('#grafana-token-message').html('<code>'+data+'</code>');
      $('#grafana-token-message').removeClass('error').addClass('positive').addClass('visible');
    },
    error: function(xhr) {
      $('#grafana-token-message').html('<code>'+xhr.responseText+'</code>');
      $('#grafana-token-message').removeClass('positive').addClass('error').addClass('visible');
    }
  });
}//end of function GetGiteeToken

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

function enableCrawl(repo_id,repo_name) {
  $('body').modal('confirm','开启爬取','确认开启爬取 <b>'+repo_name+'</b> 代码仓库吗？<br><br> 启用该仓库将在爬取周期内抓取相关数据并更新 Grfana 视图面板。', function(choice){
    if(choice) {
      $.ajax({
        type:"PUT",
        url: "/admin/repos/"+repo_id+"/change_state",
        data: {          
          "enable": true,
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function disableCrawl(repo_id,repo_name) {
  $('body').modal('confirm','禁止爬取','确认禁止爬取 <b>'+repo_name+'</b> 代码仓库吗？<br><br> 禁止该仓库将不再抓取相关数据，已有数据不受影响。', function(choice){
    if(choice) {
      $.ajax({
        type:"PUT",
        url: "/admin/repos/"+repo_id+"/change_state",
        data: {          
          "enable": false,
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function deleteRepo(repo_id,repo_name) {
  $('body').modal('confirm','删除仓库','确认删除 <b>'+repo_name+'</b> 代码仓库吗？<br><br> 删除该仓库将会同步删除相关的统计数据及 Grfana 视图面板。', function(choice){
    if(choice) {
      $.ajax({
        type:"POST",
        url: "/admin/repos/"+repo_id+"/delete",
        data: {
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function deleteCommit(sha) {
  $('body').modal('confirm','删除 Commit','确认删除 <b>'+sha+'</b> 的 Commit 记录吗？', function(choice){
    if(choice) {
      $.ajax({
        type:"POST",
        url: "/admin/commits/"+sha+"/delete",
        data: {
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function deletePR(prID) {
  $('body').modal('confirm','删除 PullRequest','确认删除 <b>'+prID+'</b> 的 PullRequest 记录吗？', function(choice){
    if(choice) {
      $.ajax({
        type:"POST",
        url: "/admin/prs/"+prID+"/delete",
        data: {
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function showAddRepoModal() {
  $('#add-repo-modal').modal({centered: false,onApprove : function() {    
    repoForm = $('#add-repo-form');    
    repoForm.form('validate form');
    if(repoForm.form('is valid')) {
      AddRepo(repoForm)
    }
    return false;
  }}).modal('show');
}

function AddRepo(form) {
  $('#add-form-loadder').removeClass('disabled').addClass('active');  
  $.ajax({
    type:"POST",
    url: "/admin/repos",
    data:{
      "repo_url" : form.form('get value', 'repo_url'),
      "type": "gitee" 
    },
    success: function() {            
      successToast('操作成功')
    },
    error: function(e) {
      errorToast($.parseJSON(e.responseText).message)
    },
    complete: function() {
      $('#add-form-loadder').removeClass('active').addClass('disabled');
    }
  });
}

function searchCommits() {
  searchForm = $('#search-commit-form');
  sha = searchForm.form('get value', 'sha');
  authorEmail = searchForm.form('get value', 'author_email');
  committerEmail = searchForm.form('get value', 'committer_email');

  if($.trim(sha).length === 0 && $.trim(authorEmail).length === 0 && $.trim(committerEmail).length === 0) {    
    return false;
  }

  searchForm.submit();
}

function searchRepos() {
  searchForm = $('#search-repo-form');
  repoName = searchForm.form('get value', 'repo_name');

  if($.trim(repoName).length === 0) {
    return false;
  }

  searchForm.submit();
}

function CreateOrUpdateGrafana(repoID) {
  $.ajax({
    type:"PUT",
    url: "/admin/grafana/"+repoID+"/update",
    data: {
      "type": "gitee"
    },
    success: function() {            
      successToast('操作成功')
    },
    error: function(e) {
      errorToast($.parseJSON(e.responseText).message)
    }
  });
}

function startToGrabRepos() {
  $('body').modal('confirm','开始爬取','确认现在就开始爬取数据吗？<br/> 如果需要爬取的仓库较多，过程可能比较久。', function(choice){
    if(choice) {
      $.ajax({
        type:"POST",
        url: "/admin/repos/grab",
        data: {
          "type": "gitee"
        },
        success: function() {
          successToast('操作成功')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }//end of if
  });
}

function reload_captcha() {
  $.ajax({
    type: "POST",
    url: '/captcha',
    dataType: 'json',
    success: function(r) {            
      $('#captcha-image').html('<img src="/captcha/'+r.result+'.png" />');
      $('<input>').attr({type: 'hidden', value:r.result ,name: 'captcha-id'}).appendTo('#login-form');
    },
    error: function(e) {
      errorToast($.parseJSON(e.responseText).message)
    }
  });
}

function sign_out() {
  $('body').modal('confirm','温馨提示','确认退出 RepoStats 管理后端吗？', function(choice){
    if (choice) {
      $.ajax({
        type:"POST",
        url: "/admin/logout",
        success: function() {
          successToast('操作成功，再见！')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }
  });
}