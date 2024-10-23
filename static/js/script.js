const body = document.querySelector("body")
const sun = document.querySelector("#sun")
const moon = document.querySelector("#moon")

// theme handling programs -----------------------------
function LightTheme(){
    body.classList.remove("dark")
    body.classList.add("light")
    sun.style.display="none"    
    moon.style.display="block"
    localStorage.setItem("theme", "light")
}

function DarkTheme(){
    body.classList.remove("light")
    body.classList.add("dark")
    moon.style.display="none"
    sun.style.display="block"    
    localStorage.setItem("theme", "dark")
}

function getTheme(){
    const theme = localStorage.getItem("theme")
    if(theme === "light"){
        LightTheme()
    }else if(theme === "dark"){
        DarkTheme()
    }else{
        LightTheme()
    }
}

// active element in navbar programs -----------------------------
const id = window.location.pathname
if(id){
    document.querySelector(`#${id.replace("/", "")}`).classList.add("bg-green-600")
}

// element show and hide programs -----------------------------
function hideTheElement(elem){
    elem.style.display = "none"
}

function showTheElement(elem){
    elem.style.display = "block"
}

// form show and hide program -----------------------------
function showForm(elem){
    var formFilter = elem.nextSibling
    var formContainer = formFilter.nextSibling

    formFilter.classList.remove("hidden")
    formContainer.classList.remove("hidden")
}

function hideFormFromFilter(elem){
    var formFilter = elem
    var formContainer = formFilter.nextSibling

    formFilter.classList.add("hidden")
    formContainer.classList.add("hidden")
}

function hideFormFromContainer(elem){
    var formContainer = elem.parentNode
    var formFilter = formContainer.previousSibling

    formFilter.classList.add("hidden")
    formContainer.classList.add("hidden")
}
