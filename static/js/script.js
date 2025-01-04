const body = document.querySelector("body")
const sun = document.querySelector("#sun")
const moon = document.querySelector("#moon")

// logout handling program
function Logout(){
    document.cookie = "Authentication="
    window.location.href = "/"
}

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

// table data download functions =======================
function downloadDataToCsv(){
    const downloadButton = document.querySelector("#download-button")
    const data = downloadButton.dataset["csv"]
    downloadCsvFile(data)
}

function downloadCsvFile(csv_data){
    let CsvFile = new Blob([csv_data], {type : "text/csv"})

    let tempLink = document.createElement("a")

    tempLink.download = "data.csv"
    let url = window.URL.createObjectURL(CsvFile)
    tempLink.href = url

    tempLink.style.display ="none"
    document.body.appendChild(tempLink)

    tempLink.click()
    document.body.removeChild(tempLink)
}


// copy function ====================================
function copyText(event){
    const copyElement = event.querySelector(".copy")
    navigator.clipboard.writeText(copyElement.innerText)
}

// import funciton 
function importCsvFile(){
    const form = document.querySelector("#import-form")
    const input = document.querySelector("#csv-file")
    if(form && input){
        form.addEventListener("submit", (e)=>{
            e.preventDefault()
        })

        input.addEventListener("input", (e)=>{
            if (e.target.files.length >0){
                const ipBtn = document.querySelector(".import-btn")
                const ipLbl = document.querySelector(".input-label")
                ipLbl.classList.add("hidden")
                ipBtn.classList.remove("hidden")
            }
        })

    }
}

if (location.pathname==="/customers" || location.pathname === "/orders"){
    importCsvFile()
}
