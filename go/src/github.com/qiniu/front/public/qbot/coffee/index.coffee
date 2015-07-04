$(->
	$('#submit-btn').on 'click', ->
    
    groups = $('#qbot-submit').find('.form-group')
    name = groups.eq(0).find('input').val()
    phone = groups.eq(1).find('input').val()
    nickname = groups.eq(2).find('input').val()
    qq = groups.eq(3).find('input').val()
    email = groups.eq(4).find('input').val()
    department = groups.eq(5).find('input').val()
    $.ajax
      method: 'POST'
      url: '/qbot/post'
      dataType: 'json'
      data:
        name: name
        phone: phone
        qq: qq
        email: email
        nickname: nickname
        department: department
    .done (data) ->
      $.growl.notice { title: "",message: data.message }
    .fail ->
      $.growl { title: "",message: "添加失败，请稍后再试" }

    return false
)

$( ->
  $('#cancel-btn').on 'click', ->
    groups = $('#qbot-submit').find('.form-group')
    for item in groups
      $(item).find('input').val ''
)

$( ->
  uploader = Qiniu.uploader 
    runtimes: 'html5,flash,html4'   
    browse_button: 'uploads'       
    uptoken: '0XPyut__IqewGS275QmjU1ZiGQuSRxZUrcM4pkRj:9TYyDRRFHhD_nQbkBqGpa37azAY=:eyJzY29wZSI6Im1hLXRlc3QtMTAwIiwiZGVhZGxpbmUiOjE0MjU1MTc3MDd9'        
    domain: 'http://qiniu-plupload.qiniudn.com/'
    container: 'photo'          
    max_file_size: '100mb'          
    flash_swf_url: 'js/plupload/Moxie.swf' 
    max_retries: 3                   
    dragdrop: true                   
    drop_element: 'photo'        
    chunk_size: '4mb'                
    auto_start: true                 
    init:
      'UploadComplete': ->
      
)
