$(document).ready(function() {
    loadCodeAndCommand($('#example-list').val())
    //load dropdown
    exampleList.map(function(example) {
        $("#example-list").append(
            `<option value="${example.name}">${example.name}</option>`
            )
    });
        
    loadParams()

    $('#run').on('click', compileCode)
    $('#example-list').on('change', function(){
        loadCodeAndCommand($(this).val())
    })

    function loadParams(){
        var queryString = window.location.search;
        var urlParams = new URLSearchParams(queryString);
        var example = urlParams.get('example');
        if(example) {
            $('#example-list').val(example)
            loadCodeAndCommand(example)
        }
        else {
            loadCodeAndCommand(exampleList[0].name)
        }
    }

    function loadCodeAndCommand(name) {
        var example = exampleList.find( x=> x.name == name)
        if(!example || example.length <1){
            return
        }
        $("#command").val(example.command)
        $.ajax({
            url : `/static/examples/${example.file}`,
            dataType: "text",
            success : function (data) {
                $("#code").val(data);
            }
        });
    }

    function compileCode(){
        //add loading
        $('#run').addClass('loading')
        //clear output
        $('#output').html('');
        //get code
        var codeSnippets = $('#code').val();
        var commandString = $('#command').val();
        var commandArgs = commandString.match(/\".+?\"|\S+/g)
        var postData = {
            code: codeSnippets,
            command: commandArgs.slice(1).map(function(value){
                return value.replace (/(^")|("$)/g, '')
            })
        }
        $.ajax({
            type: "POST",
            url: '/api/compile',
            data: JSON.stringify(postData),
            success: function(data, status, xhr) {
                result = JSON.parse(data)
                result.forEach(function(element) {
                    if(element.FileName.includes(".png")){
                        $('#output').append(`<img src="data:image/png;base64,${element.Content}" />`);
                        return
                    }
                    if(element.FileName.includes(".svg")){
                        $("#output").append(element.Content)
                        return
                    }
                    $('#output').append(`<textarea autocorrect="off" autocomplete="off" autocapitalize="off" spellcheck="false">${element.Content}</textarea>`);
                });
            },
            error: function(data) {
                $('#output').html(`<div class="error">${JSON.stringify(data)}</div>`);
            },
            complete: function() {
                $('#run').removeClass('loading')
            }
          });
    }
    function xmlToString(xmlData) { 

        var xmlString;
        //IE
        if (window.ActiveXObject){
            xmlString = xmlData.xml;
        }
        // code for Mozilla, Firefox, Opera, etc.
        else{
            xmlString = (new XMLSerializer()).serializeToString(xmlData);
        }
        return xmlString;
    } 
});

var exampleList = [
    {
        name: 'Sequence Diagram',
        command: `sysl sd -o call-login-sequence.png -s "MobileApp <- Login" tmp.sysl`,
        file: 'simple.sysl'
    },
    {
        name: 'Integration Diagram',
        command: `sysl integrations -o epa.png --project Project tmp.sysl`,
        file: 'GroceryStore.sysl'
    },
    {
        name: 'Datamodel Diagram',
        command: `sysl datamodel -d -o Payment.svg tmp.sysl`,
        file: 'Payment.sysl'
    },
    {
        name: 'Project Datamodel Diagram',
        command: `sysl datamodel -j Project -o "%(epname).svg" tmp.sysl`,
        file: 'PaymentService.sysl'
    },
    {
        name: 'Protobuf',
        command: `sysl protobuf --mode=textpb --output=simple.textpb tmp.sysl`,
        file: 'simple-pb.sysl'
    },
    {
        name: 'Code Gen',
        command: `sysl codegen --transform svc_types.sysl --grammar go.gen.g --start goFile --app-name Simple tmp.sysl`,
        file: 'simple-codegen.sysl'
    },
    {
        name: 'Import',
        command: `sysl import --input=simple-swagger.yaml --app-name=Simple --output=simple-swagger.sysl`,
        file: 'simple-swagger.yaml'
    },
    {
        name: 'Export',
        command: `sysl export --format=openapi3 --app-name=SimpleOpenAPI3 --output=simple-openapi3.yaml simple-openapi3.sysl`,
        file: 'simple-openapi3.sysl'
    }
]