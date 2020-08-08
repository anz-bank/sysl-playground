$(document).ready(function() {
    $('#run').on('click', compileCode)

    function compileCode(){
        //get code
        var codeSnippets = $('#code').val();
        var postData = {
            code: codeSnippets,
            command: 'diagram gen'
        }
        $.post( '/api/compile', postData)
        .done(function( data ) {
            $('#output').html(JSON.stringify(data));
        });
    }
});


// (function() {
//     const autocompleteUrlPrefix = '/api/autocomplete/';
//     document.getElementById('autocomplete-input').addEventListener('keyup', (ele) => {
//         //get input value
//         console.log('input')
//         let inputValue = ele.target.value;
//         //call map api to search
//         remoteCall()
//     })

//     function renderDropdown(suggestions) {
//         document.getElementById('dropdown-menu').innerHTML = ''
//         suggestions.forEach(element => {
//             document.getElementById('dropdown-menu')
//                 .insertAdjacentHTML('beforeend', 
//                 `<a href="#" class="dropdown-item">
//                 ${element.label}
//               </a>`
//                 )
//         });
//     }

//     async function remoteCall(input) {
        
//         const response = await fetch(autocompleteUrlPrefix + input)
//         const result =  await response.json();

//         renderDropdown(result.suggestions);
//     }


// })();
