const body = document.querySelector("body")
const sun = document.querySelector("#sun")
const moon = document.querySelector("#moon")

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

