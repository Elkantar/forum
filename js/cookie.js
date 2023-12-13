//const myCookieValue = getCookieValue('session-id');

/**
 * Recupere la valeur d'un cookie
 */

 function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }
  

function checkCookie() {
    let user = getCookie("session-id");
    if (user != "") {
      console.log("Welcome again " + user);
    }
  }

// utiliser la fonction pour récupérer la valeur du cookie "monCookie"
//var session_id = getCookie("session-id");
checkCookie()
