//Navbar
function hideIconBar(){
    var iconBar = document.getElementById("iconBar");
    var navigation = document.getElementById("navigation")
    iconBar.setAttribute("style", "display:none;");
    navigation.classList.remove("hide");
}

function showIconBar() {
    var iconBar = document.getElementById("iconBar");
    var navigation = document.getElementById("navigation");
    iconBar.setAttribute("style", "display:block");
    navigation.classList.add("hide");
}

//Comment
//Pas fini (je crois)

    var commentArea = document.getElementById("comment-area");
    var comButtonCom = document.getElementById("com") ;
    var replyArea = document.getElementById("reply-area");
    var comButtonRep = document.getElementById("rep") ;
    

    var check = false;
    
    comButtonCom.addEventListener("click", () => {
        
        if (check == false) {
            check = true
            commentArea.setAttribute("style", "display:block;");
        } else if (check == true) {
            check = false
            commentArea.setAttribute("style", "display:hide;");
        }

      } );
    
    
    comButtonRep.addEventListener("click", () => {

        if (check == false) {
            check = true
            replyArea.setAttribute("style", "display:block;");
        } else if (check == true) {
            check = false
            replyArea.setAttribute("style", "display:hide;");
        }

    } );
    
    const searchBar = document.getElementById('search-bar');
    const tableRows = document.querySelectorAll('.table-row');

    searchBar.addEventListener('input', () => {
        const searchTerm = searchBar.value.toLowerCase();

        tableRows.forEach(row => {
            const subject = row.querySelector('.subjects').textContent.toLowerCase();

            if (subject.includes(searchTerm)) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        });
    });


